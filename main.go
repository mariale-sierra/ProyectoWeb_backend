package main

import (
	
	"log"
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
	"ProyectoWeb_backend/handlers"
    "fmt"
    "path/filepath"
    "net/http"
)

func main() {

    db, err := sql.Open("sqlite3", "./db/series.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    absPath, _ := filepath.Abs("./db/series.db")
    fmt.Println("Using DB at:", absPath)

    http.HandleFunc("/series", handlers.GetSeries(db))
    http.HandleFunc("/add", handlers.AddSeries(db)) 
    http.HandleFunc("/update", handlers.UpdateEpisode(db))
    http.HandleFunc("/ratings", handlers.GetRatings(db))

    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}