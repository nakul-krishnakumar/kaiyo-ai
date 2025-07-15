package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nakul-krishnakumar/kaiyo-ai/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// http mux constructor
	mux := http.NewServeMux()

	// default endpoint
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Welcome to Kaiyo AI!"))

		if err != nil {
			slog.Error("Could not establish connection", slog.String("error", err.Error()))
			os.Exit(1)
		}
	})

	addr := strings.Join([]string{cfg.Host, cfg.Port}, ":")

	server := http.Server{
		Addr: addr,
		Handler: mux,
	}

	slog.Info("Server listening on http://" + addr)

	//* graceful shutdown
	done := make(chan os.Signal, 1)

	// to read interrupts
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) 

	go func() {
		err := server.ListenAndServe(); 
		if err != nil && err !=  http.ErrServerClosed {
			slog.Error("Failed to start server, " + err.Error())
			os.Exit(1)
		}
	} ()

	<-done

	slog.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down the server,", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")

}