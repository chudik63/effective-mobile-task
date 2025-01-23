package models

import "time"

// Song model godoc
// @Description Information about a song in the library
type Song struct {
	Id uint64 `json:"song_id"`
	// Musician group
	// Required: true
	Group string `json:"group"`
	// Song name
	// Required: true
	Song string `json:"song"`
	// Date when song was released
	// Required: true
	ReleaseDate time.Time `json:"release_date"`
	// Song lyrics
	Text string `json:"text"`
	// Link to the song
	Link string `json:"link"`
}

type SongInfo struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
