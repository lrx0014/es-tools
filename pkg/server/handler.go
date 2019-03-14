package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kubeapps/common/response"
	"k8s.io/klog/glog"

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
		glog.Error(message)
		response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
		return
	}
	result, err := log.GetLogs(client, namespace, podID, container)
	r, _ := ioutil.ReadAll(result)
	if err != nil {
		message := fmt.Sprintf("Unable to get log... => %v\n", err)
		glog.Error(message)
		response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
		return
	}
	renderer.JSON(w, http.StatusOK, r)
}
