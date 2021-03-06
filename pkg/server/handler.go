package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang/glog"

	"github.com/kubeapps/common/response"

	"github.com/gorilla/mux"
	"github.com/lrx0014/log-tools/pkg/api"
	"github.com/lrx0014/log-tools/pkg/es"
	"github.com/lrx0014/log-tools/pkg/log"
	"github.com/unrolled/render"
)

var renderer *render.Render

func init() {
	renderer = render.New(render.Options{})
}

// Params a key-value map of path params
type Params map[string]string

// WithParams can be used to wrap handlers to take an extra arg for path params
type WithParams func(http.ResponseWriter, *http.Request, Params)

func (h WithParams) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	h(w, request, vars)
}

func bind(request *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(request.Body)
	return decoder.Decode(obj)
}

func (s *APIServer) handleLogs(w http.ResponseWriter, request *http.Request, params Params) {
	s.streamLogs(w, request, params)
}

func (s *APIServer) handlerListPods(w http.ResponseWriter, request *http.Request, params Params) {
	namespace := params["namespace"]
	container := params["container"]
	k8sClient, err := log.CreateClient()
	if err != nil {
		message := fmt.Sprintf("Unable to create k8s client... => %v\n", err)
		glog.Errorln(message)
		response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
		return
	}
	podlist, err := es.GetPodList(k8sClient, namespace, container)
	if err != nil {
		message := fmt.Sprintf("Unable to get pod list from es... => %v\n", err)
		glog.Errorln(message)
		response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
		return
	}
	renderer.JSON(w, http.StatusOK, podlist)
}

func (s *APIServer) streamLogs(w http.ResponseWriter, request *http.Request, params Params) {
	instance := &api.ContainerInfo{
		Namespace: params["namespace"],
		PodID:     params["pod"],
		Container: params["container"],
	}
	client, err := log.CreateClient()
	if err != nil {
		message := fmt.Sprintf("Unable to create k8s client... => %v\n", err)
		glog.Errorln(message)
		response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
		return
	}

	cn, ok := w.(http.CloseNotifier)
	if !ok {
		http.NotFound(w, request)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.NotFound(w, request)
		return
	}

	content, err := log.StreamLogs(client, instance)
	if err != nil {
		message := fmt.Sprintf("Unable to get log... => %v\n", err)
		glog.Errorln(message)
		response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
		return
	}

	defer content.Close()

	// Send the initial headers saying we're gonna stream the response.
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	for {
		select {
		case <-cn.CloseNotify():
			glog.Infoln("Client stopped listening")
			return
		default:
			// Send some data.
			buf := make([]byte, 2048)
			n, err := content.Read(buf)
			if err != nil && err != io.EOF {
				message := fmt.Sprintf("Unable to stream logs: %v", err)
				response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
				return
			}

			if err == io.EOF {
				break
			}

			glog.Infof("Sending some data: %s", string(buf[:n]))

			w.Write(buf[:n])

			flusher.Flush()
		}
	}
}
