package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	spotify "iffy-blue-analytica/internal/spotify"
	core "iffy-blue-analytica/internal/core"
	db "iffy-blue-analytica/sql"
)

func main() {
	// load env vars
	err := godotenv.Load("../../.env")
	if err != nil {log.Fatal("Loading envs: ", err)}

	// create tables
	err = db.CreateTables(false)
	if err != nil {log.Fatal("Table create: ", err)}

	// setup router
	r, err := setupRouter()
	if err != nil {log.Fatal("Router setup: ", err)}

	// run server
	log.Println("Server listening at: localhost:" + os.Getenv("PORT") + "/")
	http.ListenAndServe(":" + os.Getenv("PORT"), r)
}

func setupRouter() (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Use(middleware.Logger) // default logger
	setupRoutes(r)
	return r, nil
}

func setupRoutes(r *chi.Mux) {
	// PRINTERS
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})

	r.Get("/print/successful-login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("successful Spotify login, nice!"))
	}) 

	// AUTH 
	r.Get("/login", spotify.Login)
	r.Get("/login/callback", spotify.Callback)

	// CORE
	r.Get("/tables/reset", core.ResetTables)
	r.Get("/tables/populate", core.PopulateTables)
	r.Get("/tables/populate/max={maxTracks}", core.PopulateTables)
}