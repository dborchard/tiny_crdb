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

func (s *httpServer) setupRoutes(
	ctx context.Context,
	authnServer authserver.Server,
) error {

	// Define the http.Handler for UI assets.
	assetHandler := ui.Handler()

	// The authentication mux used here is created in "allow anonymous" mode so that the UI
	// assets are served up whether or not there is a session. If there is a session, the mux
	// adds it to the context, and it is templated into index.html so that the UI can show
	// the username of the currently-logged-in user.
	authenticatedUIHandler := authserver.NewMux(
		authnServer, assetHandler, true /* allowAnonymous */)
	s.mux.Handle("/", authenticatedUIHandler)

	return nil
}
