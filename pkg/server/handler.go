package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	debug "log"

	"github.com/kubeapps/common/response"

	"github.com/gorilla/mux"
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

func (s *APIServer) getLogs(w http.ResponseWriter, request *http.Request, params Params) {
	namespace := params["namespace"]
	podID := params["pod"]
	container := params["container"]
	client, err := log.CreateClient()
	if err != nil {
		message := fmt.Sprintf("Unable to create k8s client... => %v\n", err)
		debug.Println(message)
		response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
		return
	}
	result, err := log.GetLogs(client, namespace, podID, container)
	if err != nil {
		message := fmt.Sprintf("Unable to get log... => %v\n", err)
		debug.Println(message)
		response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
		return
	}
	renderer.JSON(w, http.StatusOK, result)
}

func (s *APIServer) streamLogs(w http.ResponseWriter, request *http.Request, params Params) {
	namespace := params["namespace"]
	podID := params["pod"]
	container := params["container"]
	client, err := log.CreateClient()
	if err != nil {
		message := fmt.Sprintf("Unable to create k8s client... => %v\n", err)
		debug.Println(message)
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

	content, err := log.StreamLogs(client, namespace, podID, container)
	if err != nil {
		message := fmt.Sprintf("Unable to get log... => %v\n", err)
		debug.Println(message)
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
			debug.Println("Client stopped listening")
			return
		default:
			/*
				result, err := ioutil.ReadAll(content)
				if err != nil {
					message := fmt.Sprintf("Unable to read all logs: %v", err)
					response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
					return
				}
				debug.Printf("log: %s\n", string(result))
			*/
			time.Sleep(2 * time.Second)
			// Send some data.
			var buf []byte
			n, err := content.Read(buf)
			if err != nil {
				message := fmt.Sprintf("Unable to stream logs: %v", err)
				response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
				return
			}

			debug.Printf("Sending some data: %d", n)

			w.Write(buf)

			flusher.Flush()
		}
	}
}
