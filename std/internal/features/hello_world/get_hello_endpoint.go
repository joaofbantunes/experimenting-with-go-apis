package hello_world

import (
	"log/slog"
	"net/http"

	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
)

func GetHelloEndpoint(loggerProvider func(name string) *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	logger := loggerProvider("get_hello_endpoint")
	return func(w http.ResponseWriter, r *http.Request) {

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
			shared.InternalServerError(w, logger, err)
			return
		}
	}
}
