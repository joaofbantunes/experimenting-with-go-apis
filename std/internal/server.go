package internal

import (
	"net/http"
)

func NewServer(compositionRoot *CompositionRoot) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, compositionRoot)
	var handler http.Handler = mux
	// TODO: add middleware here if needed
	// handler = someMiddleware(handler)
	return handler
}
