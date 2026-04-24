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

		query := r.URL.Query().Get("q")
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

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
			series, err = db.GetSeriesPaginated(database, limit, offset)
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

func SeriesRating(database *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/rating") {
			if r.Method == "GET" {
				GetRatings(database)(w, r)
				return
			}

			if r.Method == "POST" {
				AddOrUpdateRating(database)(w, r)
				return
			}
		}

		http.NotFound(w, r)
	}
}

func GetAllRatings(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := db.GetAllRatings(database)
		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		json.NewEncoder(w).Encode(data)
	}
}

func GetRatings(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := strings.TrimPrefix(r.URL.Path, "/series/")
		idStr = strings.TrimSuffix(idStr, "/rating")

		seriesID, _ := strconv.Atoi(idStr)

		ratings, err := db.GetRatingsBySeries(database, seriesID)
		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		json.NewEncoder(w).Encode(ratings)
	}
}

func AddOrUpdateRating(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/series/")
		idStr = strings.TrimSuffix(idStr, "/rating")

		seriesID, _ := strconv.Atoi(idStr)

		var body struct {
			Rating  int `json:"rating"`
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
