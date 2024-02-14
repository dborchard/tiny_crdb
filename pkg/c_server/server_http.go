package server

import (
	"context"
	ui "github.com/dborchard/tiny_crdb/pkg/e_ui"
	"net/http"
)

func (s *httpServer) setupRoutes(
	ctx context.Context,
	authnServer authserver.Server,
	adminAuthzCheck privchecker.CheckerForRPCHandlers,
	metricSource metricMarshaler,
	runtimeStatSampler *status.RuntimeStatSampler,
	handleRequestsUnauthenticated http.Handler,
	handleDebugUnauthenticated http.Handler,
	handleInspectzUnauthenticated http.Handler,
	apiServer http.Handler,
	flags serverpb.FeatureFlags,
) error {

	// Define the http.Handler for UI assets.
	assetHandler := ui.Handler(ui.Config{
		Insecure: s.cfg.InsecureWebAccess(),
		NodeID:   s.cfg.IDContainer,
		OIDC:     oidc,
		GetUser: func(ctx context.Context) *string {
			if user, ok := authserver.MaybeUserFromHTTPAuthInfoContext(ctx); ok {
				ustring := user.Normalized()
				return &ustring
			}
			return nil
		},
		Flags: flags,
	})

	// The authentication mux used here is created in "allow anonymous" mode so that the UI
	// assets are served up whether or not there is a session. If there is a session, the mux
	// adds it to the context, and it is templated into index.html so that the UI can show
	// the username of the currently-logged-in user.
	authenticatedUIHandler := authserver.NewMux(
		authnServer, assetHandler, true /* allowAnonymous */)
	s.mux.Handle("/", authenticatedUIHandler)

}
