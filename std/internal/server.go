package internal

import (
	"fmt"

	"log/slog"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewServer(root *CompositionRoot) http.Handler {
	mux := http.NewServeMux()

	// replacement for mux.Handle to add OpenTelemetry instrumentation
	handle := func(pattern string, handler http.Handler) {
		instrumented := otelhttp.NewHandler(
			handler,
			pattern,
			otelhttp.WithTracerProvider(root.O11y.TracerProvider),
			otelhttp.WithMeterProvider(root.O11y.MeterProvider))
		instrumented = otelhttp.WithRouteTag(pattern, instrumented)
		mux.Handle(pattern, instrumented)
	}
	addRoutes(handle, root)
	addScalar(handle, root.CreateLogger("scalar"))
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
