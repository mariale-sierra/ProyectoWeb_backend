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

        query := r.URL.Query().Get("q")

        var (
            series []db.Series
            err    error
        )

        if query != "" {
            series, err = db.SearchSeries(database, query)
        } else {
            series, err = db.GetAllSeries(database)
        }

        if err != nil {
            http.Error(w, "Error interno", 500)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(series)
    }
}


func AddSeries(database *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {   
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        if r.Method != "POST" {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var newSeries struct {
            Name    string `json:"name"`
            Current int    `json:"current_episode"`
            Total   int    `json:"total_episodes"`
        }

        err := json.NewDecoder(r.Body).Decode(&newSeries)
        if err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }

        _, err = database.Exec(
            "INSERT INTO series (name, current_episode, total_episodes) VALUES (?, ?, ?)",
            newSeries.Name, newSeries.Current, newSeries.Total,
        )

        if err != nil {
            http.Error(w, "DB error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{
            "message": "Series created",
        })
    }
}

func GetRatings(database *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        rows, err := database.Query(`
            SELECT s.id, s.name, r.rating
            FROM ratings r
            JOIN series s ON s.id = r.series_id
        `)

        if err != nil {
            http.Error(w, "DB error", 500)
            return
        }
        defer rows.Close()

        var results []map[string]interface{}

        for rows.Next() {
            var id int
            var name string
            var rating int

            err := rows.Scan(&id, &name, &rating)
            if err != nil {
                http.Error(w, "Scan error", 500)
                return
            }

            results = append(results, map[string]interface{}{
                "id":     id,
                "name":   name,
                "rating": rating,
            })
        }

        json.NewEncoder(w).Encode(results)
    }
}

func UpdateEpisode(database *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        // CORS
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        if r.Method != "POST" {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        id := r.URL.Query().Get("id")

        if id == "" {
            http.Error(w, "Missing id", http.StatusBadRequest)
            return
        }

        _, err := database.Exec(`
            UPDATE series
            SET current_episode = current_episode + 1
            WHERE id = ? AND current_episode < total_episodes
        `, id)

        if err != nil {
            http.Error(w, "DB error", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{
            "message": "Episode updated",
        })
    }
}