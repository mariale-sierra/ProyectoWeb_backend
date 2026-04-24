package handlers

import (
	"ProyectoWeb_backend/db"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetSeries(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		query := r.URL.Query().Get("q")
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		sort := r.URL.Query().Get("sort")
		order := r.URL.Query().Get("order")

		allowedSort := map[string]bool{
			"id":              true,
			"name":            true,
			"current_episode": true,
			"total_episodes":  true,
		}

		if !allowedSort[sort] {
			sort = "id"
		}

		if order != "desc" {
			order = "asc"
		}

		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		var (
			series []db.Series
			err    error
		)

		if query != "" {
			series, err = db.SearchSeries(database, query)
		} else {
			series, err = db.GetSeriesPaginatedSorted(database, limit, offset, sort, order)
		}

		if err != nil {
			http.Error(w, "Error interno", 500)
			return
		}

		json.NewEncoder(w).Encode(series)
	}
}

func AddSeries(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		var body struct {
			ID int `json:"id"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		if body.ID < 1 {
			http.Error(w, "Invalid series id", 400)
			return
		}

		result, err := database.Exec(`
			UPDATE series
			SET added = 1
			WHERE id = ?
		`, body.ID)

		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		log.Printf("AddSeries id=%d rows_updated=%d", body.ID, rowsAffected)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":      "Series added",
			"rows_updated": rowsAffected,
		})
	}
}

func GetAllRatings(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		data, err := db.GetAllRatings(database)
		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		json.NewEncoder(w).Encode(data)
	}
}


func AddOrUpdateRating(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/series/")
		idStr = strings.TrimSuffix(idStr, "/rating")

		seriesID, _ := strconv.Atoi(idStr)

		var body struct {
			Rating int `json:"rating"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		err = db.EditRating(database, seriesID, body.Rating)
		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "rating saved",
		})
	}
}

func UpdateEpisode(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		id := r.URL.Query().Get("id")

		err := db.UpdateEpisode(database, id)
		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Episode updated",
		})
	}
}

func CreateSeries(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		var body struct {
			Name  string `json:"name"`
			Total int    `json:"total_episodes"`
			Image string `json:"image"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		_, err = database.Exec(`
			INSERT INTO series (name, current_episode, total_episodes, added, image)
			VALUES (?, 1, ?, 1, ?)
		`, body.Name, body.Total, body.Image)

		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Series created",
		})
	}
}

func DeleteSeries(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "DELETE" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/series/")

		_, err := database.Exec("DELETE FROM series WHERE id = ?", id)
		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "deleted",
		})
	}
}

func SeriesHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		path := r.URL.Path

		if strings.HasSuffix(path, "/rating") {
			if r.Method == "GET" {
				GetAllRatings(database)(w, r)
				return
			}
			if r.Method == "POST" {
				AddOrUpdateRating(database)(w, r)
				return
			}
		}

		if r.Method == "DELETE" {
			DeleteSeries(database)(w, r)
			return
		}

		http.NotFound(w, r)
	}
}
