package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/willzfrank/redteam-hackathon/target-api/internal/config"
	"github.com/willzfrank/redteam-hackathon/target-api/internal/handlers"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.Get("/accounts/{id}", handlers.HandleGetAccount)
	router.Get("/users", handlers.HandleGetUser)
	router.Post("/login", handlers.HandleLogin)

	log.Fatal(http.ListenAndServe(config.Load().Port, router))

	srv := &http.Server{
		Addr:    config.Load().Port,
		Handler: router,
	}

	// Start the server in a separate goroutine so it doesn't block main()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// main() continues immediately here, without waiting for the server to finish
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit // main() blocks HERE instead, waiting for Ctrl+C

	// once Ctrl+C happens, we get here and gracefully shut down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("Server shut down cleanly")

}
