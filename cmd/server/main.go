package main

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/nakul-krishnakumar/kaiyo-ai/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad("local", "./config")


	// http mux
	mux := http.NewServeMux()

	// default endpoint
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Welcome to Kaiyo AI!"))

		if err != nil {
			slog.Error("Could not establish connection", slog.String("error", err.Error()))
		}
	})

	addr := strings.Join([]string{cfg.Host, cfg.Port}, ":")

	server := http.Server{
		Addr: addr,
		Handler: mux,
	}

	slog.Info("Server listening on", slog.String("address", "http://" + addr))

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server,", slog.String("error", err.Error()))
	}


}