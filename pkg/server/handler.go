package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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

/*
func (s *APIServer) listCatalog(w http.ResponseWriter, request *http.Request) {
	catalogList, err := s.CatalogManager.ListCatalog()
	if err != nil {
		message := fmt.Sprintf("can not fetch all catalogs: %v", err)
		glog.Error(message)
		response.NewErrorResponse(http.StatusInternalServerError, message).Write(w)
		return
	}
	renderer.JSON(w, http.StatusOK, catalogList)
}
*/
