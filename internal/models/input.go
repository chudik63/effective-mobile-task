package models

type AddSongInput struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type UpdateSongInput struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
