// Package main implements the TaskFlow API server.
//
// @title TaskFlow API
// @version 1.0
// @description TaskFlow MVP backend (Go implementation).
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"taskflow-backend/internal/app"
	"taskflow-backend/internal/config"
	"taskflow-backend/internal/db"
	"taskflow-backend/internal/server"
)

func main() {
	settings := config.Load()

	conn, err := db.Connect(settings.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close()

	a := &app.App{DB: conn, Settings: settings}

	router := chi.NewRouter()
	router.Get("/openapi.json", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "docs/swagger.json")
	})
	router.Handle("/docs/swaggerui/*", http.StripPrefix("/docs/", swaggerUIFileHandler))
	router.Get("/docs", swaggerUIHandler)
	router.Get("/docs/*", swaggerUIHandler)
	router.Mount("/", server.NewRouter(a))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("TaskFlow Go backend listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
