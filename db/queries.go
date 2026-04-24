package db

import "database/sql"

func SearchSeries(database *sql.DB, query string) ([]Series, error) {
	rows, err := database.Query(`
		SELECT id, name, current_episode, total_episodes
		FROM series
		WHERE name LIKE ?
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Series

	for rows.Next() {
		var s Series
		err := rows.Scan(&s.ID, &s.Name, &s.Current, &s.Total)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}

	return list, nil
}

func GetSeriesPaginated(database *sql.DB, limit int, offset int) ([]Series, error) {
	rows, err := database.Query(`
		SELECT id, name, current_episode, total_episodes
		FROM series
		WHERE added = 1
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Series

	for rows.Next() {
		var s Series
		err := rows.Scan(&s.ID, &s.Name, &s.Current, &s.Total)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}

	return list, nil
}

func EditRating(database *sql.DB, seriesID int, rating int) error {
	_, err := database.Exec(`
		INSERT INTO ratings (series_id, rating)
		VALUES (?, ?)
		ON CONFLICT(series_id)
		DO UPDATE SET rating = excluded.rating
	`, seriesID, rating)

	return err
}

func GetRatingsBySeries(database *sql.DB, seriesID int) ([]map[string]int, error) {
	rows, err := database.Query(`
		SELECT episode, rating
		FROM ratings
		WHERE series_id = ?
	`, seriesID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []map[string]int{}

	for rows.Next() {
		var ep, rating int
		rows.Scan(&ep, &rating)

		list = append(list, map[string]int{
			"episode": ep,
			"rating":  rating,
		})
	}

	return list, nil
}

func GetAllRatings(database *sql.DB) ([]map[string]interface{}, error) {
	rows, err := database.Query(`
		SELECT 
    		s.id,
    		s.name,
   			r.rating
		FROM series s
		LEFT JOIN ratings r ON s.id = r.series_id
		WHERE s.added = 1
		`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []map[string]interface{}{}

	for rows.Next() {
		var id int
		var name string
		var rating sql.NullInt64

		rows.Scan(&id, &name, &rating)

		list = append(list, map[string]interface{}{
			"id": id,
			"name": name,
			"rating": func() interface{} {
				if rating.Valid {
					return rating.Int64
				}
				return nil
			}(),
		})
	}

	return list, nil
}

func UpdateEpisode(database *sql.DB, id string) error {
	_, err := database.Exec(`
		UPDATE series
		SET current_episode = current_episode + 1
		WHERE id = ? AND current_episode < total_episodes
	`, id)

	return err
}
