package db

import "database/sql"

func SearchSeries(database *sql.DB, query string) ([]Series, error) {
	rows, err := database.Query(`
		SELECT id, name, current_episode, total_episodes, image
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
		err := rows.Scan(&s.ID, &s.Name, &s.Current, &s.Total, &s.Image)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}

	return list, nil
}


func GetSeriesPaginatedSorted(database *sql.DB, limit int, offset int, sort string, order string) ([]Series, error) {
	query := `
		SELECT id, name, current_episode, total_episodes, image
		FROM series
		WHERE added = 1
		ORDER BY ` + sort + ` ` + order + `
		LIMIT ? OFFSET ?
	`

	rows, err := database.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Series

	for rows.Next() {
		var s Series
		err := rows.Scan(&s.ID, &s.Name, &s.Current, &s.Total, &s.Image)
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

func GetAllRatings(database *sql.DB) ([]SeriesRating, error) {
    rows, err := database.Query(`
        SELECT s.id, s.name, r.rating
        FROM series s
        LEFT JOIN ratings r ON s.id = r.series_id
        WHERE s.added = 1
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    list := []SeriesRating{}

    for rows.Next() {
        var s SeriesRating
        rows.Scan(&s.ID, &s.Name, &s.Rating)
        list = append(list, s)
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
