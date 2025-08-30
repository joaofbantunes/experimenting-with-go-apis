package internal

import (
	"fmt"

	"log/slog"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewServer(compositionRoot *CompositionRoot) http.Handler {
	mux := http.NewServeMux()

	// replacement for mux.Handle to enrich traces with route pattern
	handle := func(pattern string, handler http.Handler) {
		instrumented := otelhttp.NewHandler(handler, pattern)
		instrumented = otelhttp.WithRouteTag(pattern, instrumented)
		mux.Handle(pattern, instrumented)
	}
	addRoutes(handle, compositionRoot)
	addScalar(handle, compositionRoot.CreateLogger("scalar"))
	return mux
}

func addScalar(handle func(pattern string, handler http.Handler), logger *slog.Logger) {
	handle("GET /scalar", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
}
