package internal

import (
	"net/http"

	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/hello_world"
)

func addRoutes(mux *http.ServeMux, root *CompositionRoot) {
	mux.HandleFunc("GET /hello", hello_world.NewHelloEndpoint(root.LoggerProvider))
}
