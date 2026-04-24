package main

import (
	"ProyectoWeb_backend/handlers"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "./db/series.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    http.HandleFunc("/create", handlers.CreateSeries(db))
    http.HandleFunc("/series", handlers.GetSeries(db))
    http.HandleFunc("/series/", handlers.SeriesHandler(db))
    http.HandleFunc("/add", handlers.AddSeries(db))
    http.HandleFunc("/update", handlers.UpdateEpisode(db))
    http.HandleFunc("/ratings", handlers.GetAllRatings(db))

    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}