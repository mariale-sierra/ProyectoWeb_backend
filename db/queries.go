package db

import "database/sql"


func GetAllSeries(database *sql.DB) ([]Series, error) {

	rows, err := database.Query(`
		SELECT s.id, s.name, s.current_episode, s.total_episodes, r.rating
		FROM series s
		LEFT JOIN ratings r ON s.id = r.series_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Series

	for rows.Next() {
		var s Series

		err := rows.Scan(&s.ID, &s.Name, &s.Current, &s.Total, &s.Rating)
		if err != nil {
			return nil, err   
		}

		list = append(list, s)
	}
	return list, nil
}