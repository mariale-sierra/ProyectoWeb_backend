package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"ProyectoWeb_backend/db"
)

func GetSeries(database *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        series, err := db.GetAllSeries(database)
        if err != nil {
            http.Error(w, "Error interno", 500)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(series)
    }
}