package server

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/c_server/authserver"
	ui "github.com/dborchard/tiny_crdb/pkg/e_ui"
	"net/http"
)

type httpServer struct {
	mux   http.ServeMux
	gzMux http.Handler
}

func newHTTPServer() *httpServer {
	server := &httpServer{}
	return server
}

func (s *httpServer) setupRoutes(ctx context.Context, authnServer authserver.Server) error {
	assetHandler := ui.Handler()
	authenticatedUIHandler := authserver.NewMux(authnServer, assetHandler, true)
	s.mux.Handle("/", authenticatedUIHandler)
	return nil
}
