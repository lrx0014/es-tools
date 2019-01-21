package server

import (
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"github.com/urfave/negroni"

	"github.com/lrx0014/log-tools/cmd/apiserver/app/options"
	"github.com/lrx0014/log-tools/pkg/log"
)

// APIServer is a http.Handler which exposes catalog aggregator functionality over HTTP.
type APIServer struct {
	cfg  *options.ServerRunOptions
	logs *log.LogService
}

func NewAPIServer(cfg *options.ServerRunOptions) (*APIServer, error) {
	server := &APIServer{
		cfg: cfg,
	}

	return server, nil
}

func (s *APIServer) Run() error {

	routes := s.setupRoutes()

	return http.ListenAndServe(net.JoinHostPort(s.cfg.Address, strconv.FormatUint(uint64(s.cfg.Port), 10)), routes)
}

// setupRoutes registers a set of supported HTTP request patterns
func (s *APIServer) setupRoutes() http.Handler {
	r := mux.NewRouter()

	// Healthcheck
	health := healthcheck.NewHandler()
	r.Handle("/live", health)
	r.Handle("/ready", health)

	// Routes
	// apiv1 := r.PathPrefix("/v1").Subrouter()

	/*
		apiv1.Methods("GET").Path("/catalogs").HandlerFunc(s.listCatalog)
		apiv1.Methods("GET").Path("/catalogs/{catalog}").Handler(WithParams(s.getCatalog))
	*/

	n := negroni.Classic()
	n.UseHandler(r)
	return n
}
