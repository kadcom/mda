package main

import (
	"context"
	"mda/todo"
	"mda/todo/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {

	connStr := `host=localhost port=5432 dbname=kad_todo user=postgres`

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)

	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
	}

	todo.SetPool(pool)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Get("/todo", handlers.ListItems)
	r.Get("/todo/{itemId}", handlers.GetItem)
	r.Post("/todo", handlers.CreateItem)
	r.Post("/todo/done", handlers.MakeItemDone)

	log.Info().Msg("Starting up server...")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server")
		return
	}

	log.Info().Msg("Server Stopped")
}
