package hello_world

import (
	"log/slog"
	"net/http"

	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
)

func NewHelloEndpoint(loggerProvider func(name string) *slog.Logger) http.Handler {
	logger := loggerProvider("get_hello_endpoint")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type response struct {
			Greeting string `json:"greeting"`
		}

		err := shared.Encode(
			w,
			http.StatusOK,
			response{
				Greeting: "Hello, World!",
			})

		if err != nil {
			shared.EncodeInternalServerError(r.Context(), w, logger, err)
			return
		}
	})
}
