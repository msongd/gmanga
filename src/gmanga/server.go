package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ServerCfg struct {
	Address string
	Route   *mux.Router
	Context *AppContext
	Server  *http.Server
}

func NewServerCfg(address string, c *AppContext) *ServerCfg {
	cfg := ServerCfg{Address: address, Route: nil, Context: c}
	cfg.SetupRoute()
	return &cfg
}

func (cfg *ServerCfg) SetupRoute() {
	cfg.Route = DefineRoutes(cfg.Context)
}

func (cfg *ServerCfg) Serve() {

	loggedRouter := handlers.CombinedLoggingHandler(os.Stderr, cfg.Route)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: loggedRouter,
	}
	cfg.Server = server
	server.ListenAndServe()
}

func (cfg *ServerCfg) ServeTLS(serverCert string, serverKey string) {
	loggedRouter := handlers.CombinedLoggingHandler(os.Stderr, cfg.Route)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: loggedRouter,
	}
	cfg.Server = server
	server.ListenAndServeTLS(serverCert, serverKey)
}
