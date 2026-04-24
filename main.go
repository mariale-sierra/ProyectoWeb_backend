package main

import (
	"ProyectoWeb_backend/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/series.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	absPath, _ := filepath.Abs("./db/series.db")
	fmt.Println("Using DB at:", absPath)

	mux := http.NewServeMux()
	mux.HandleFunc("/series", handlers.GetSeries(db))
	mux.HandleFunc("/series/", handlers.SeriesRating(db))
	mux.HandleFunc("/add", handlers.AddSeries(db))
	mux.HandleFunc("/update", handlers.UpdateEpisode(db))
	mux.HandleFunc("/ratings", handlers.GetAllRatings(db))

	loggedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    loggedMux.ServeHTTP(w, r)
})))
}
