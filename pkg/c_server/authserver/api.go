package authserver

import "net/http"

type Server interface {
}

type AuthMux interface {
	http.Handler
}

func NewMux(s Server, inner http.Handler, allowAnonymous bool) AuthMux {
	return &authenticationMux{
		server:         s.(*authenticationServer),
		inner:          inner,
		allowAnonymous: allowAnonymous,
	}
}

// authenticationMux implements http.Handler, and is used to provide session
// authentication for an arbitrary "inner" handler.
type authenticationMux struct {
	server         *authenticationServer
	inner          http.Handler
	allowAnonymous bool
}

func (am *authenticationMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//TODO implement me
	panic("implement me")
}
