package db

// USAGE: database models

type Track struct {
	TrackID		string
	Name		string
	URI			string
	AlbumID		string
}

// NOTE: link table
type TrackArtist struct {
	TrackID		string
	ArtistID	string
}

type Artist struct {
	ArtistID		string
	Name			string
	URI				string
	Genres			[]string
}

type Album struct {
	AlbumID			string
	Name			string
	URI				string
	Genres			[]string
}