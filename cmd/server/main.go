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
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/http/auth"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/http/chat"
	mw "github.com/nakul-krishnakumar/kaiyo-ai/internal/middlewares"
)

func main() {
	// load configs
	cfg := config.MustLoad()
	authConfig := auth.MustLoad()

	// middlewares
	Authenticate := mw.Auth(authConfig)

	// http mux constructor
	mainMux := http.NewServeMux()

	apiMux := http.NewServeMux()
	apiMux.Handle("/chats/", http.StripPrefix("/chats", chat.New())) // /api/v1/chats
	
	mainMux.Handle("/api/v1/", Authenticate(http.StripPrefix("/api/v1", apiMux))) // /api/v1/
	mainMux.Handle("/auth/", http.StripPrefix("/auth", auth.New(authConfig)))   // /auth


	// default endpoint - {$} makes it very specific
	mainMux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Welcome to Kaiyo AI!"))

		if err != nil {
			slog.Error("Could not write response", slog.String("error", err.Error()))
		}
	})

	addr := strings.Join([]string{cfg.Host, cfg.Port}, ":")

	server := http.Server{
		Addr: addr,
		Handler: mw.CORS(mainMux),
	}

	slog.Info("Server listening on http://" + addr)

	//* graceful shutdown
	done := make(chan os.Signal, 1)

	// to read interrupts
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) 
	// when any of the above mention signals is observed, it sent it into the 'done' channel which stop the goroutine

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