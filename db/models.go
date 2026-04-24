
package db
type Series struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Current int    `json:"current_episode"`
	Total   int    `json:"total_episodes"`
	Image   *string `json:"image"`
}

type SeriesRating struct {
    ID     int  `json:"id"`
    Name   string `json:"name"`
    Rating *int   `json:"rating"`
}