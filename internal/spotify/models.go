package spotify

// USAGE: models representing what we receive from Spotify API

type TrackResponse struct {
	Items		[]TrackItem		`json:"items"`
	Limit		int				`json:"limit"`
	Next		*string			`json:"next"`
	Offset		int				`json:"offset"`
	TotalItems	int				`json:"total"`
}

type TrackItem struct {
	Track	Track	`json:"track"`
}

type Track struct {
	Name  	string    	`json:"name"`
	ID    	string 		`json:"id"`
	URI   	string 		`json:"uri"`
	Album	Album		`json:"album"`
	Artists []Artist	`json:"artists"`
}

type Album struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	URI    string   `json:"uri"`
}

type Artist struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	URI    string   `json:"uri"`
}