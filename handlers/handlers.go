package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"ProyectoWeb_backend/db"
)

func GetSeries(conn *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        series, err := db.GetAllSeries(conn)

        if err != nil {
            http.Error(w, "Error interno", 500)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(series)
    }
}