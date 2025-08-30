package internal

import (
	"fmt"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"log/slog"
	"net/http"
)

func NewServer(compositionRoot *CompositionRoot) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, compositionRoot)
	addScalar(mux, compositionRoot.CreateLogger("scalar"))
	var handler http.Handler = mux
	// TODO: add middleware here if needed
	// handler = someMiddleware(handler)
	handler = otelhttp.NewHandler(handler, "/")
	return handler
}

func addScalar(mux *http.ServeMux, logger *slog.Logger) {
	mux.HandleFunc("GET /scalar", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./api/v1.yaml",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Sample API",
			},
			DarkMode: true,
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		_, err = fmt.Fprintln(w, htmlContent)
		if err != nil {
			shared.EncodeInternalServerError(r.Context(), w, logger, err)
			return
		}
	})
}
