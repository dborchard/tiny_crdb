package authserver

import "net/http"

type Server interface {
}
type AuthMux interface {
	http.Handler
}

// authenticationMux implements http.Handler, and is used to provide session
// authentication for an arbitrary "inner" handler.
type authenticationMux struct {
	server *authenticationServer

	inner http.Handler

	// allowAnonymous, if true, indicates that the authentication mux should
	// call its inner HTTP handler even if the request doesn't have a valid
	// session. If there is a valid session, the mux calls its inner handler
	// with a context containing the username and session ID.
	//
	// If allowAnonymous is false, the mux returns an error if there is no
	// valid session.
	allowAnonymous bool
}

func (am *authenticationMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//TODO implement me
	panic("implement me")
}

func NewMux(s Server, inner http.Handler, allowAnonymous bool) AuthMux {
	return &authenticationMux{
		server:         s.(*authenticationServer),
		inner:          inner,
		allowAnonymous: allowAnonymous,
	}
}
