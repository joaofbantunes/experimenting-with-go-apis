package shared

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func InternalServerError(w http.ResponseWriter, logger *slog.Logger, err error) {
	logger.Error("internal server error", slog.Any("error", err))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
