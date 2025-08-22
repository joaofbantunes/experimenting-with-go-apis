package main

import (
	"context"
	"errors"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal"
)

func run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	config, err := internal.LoadConfig()
	if err != nil {
		return err
	}

	compositionRoot := internal.NewCompositionRoot(config)

	compositionRoot.InitApp()

	logger := compositionRoot.LoggerProvider("main")

	srv := internal.NewServer(compositionRoot)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(compositionRoot.Config.Server.Host, compositionRoot.Config.Server.Port),
		Handler: srv,
	}
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorContext(ctx, "error listening and serving", slog.Any("error", err))
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.ErrorContext(ctx, "error shutting down server", slog.Any("error", err))
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(
		ctx,
		os.Args,
		func(s string) string {
			return os.Getenv(s)
		},
		os.Stdin,
		os.Stdout,
		os.Stderr,
	); err != nil {
		log.Fatalf("error: %v", err)
	}
}
