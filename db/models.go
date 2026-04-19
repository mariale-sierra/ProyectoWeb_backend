package db

import "database/sql"

type Series struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Current int            `json:"current_episode"`
	Total   int            `json:"total_episodes"`
	Rating  sql.NullInt64  `json:"rating"`
}