package core

import (
	"log"
	"net/http"

	spotify "iffy-blue-analytica/internal/spotify"
	db "iffy-blue-analytica/sql"
)

// USAGE: create data in DB from track response
func uploadTracksToDB(tracks []db.Track) error {

	// set maximum rows to insert with a single command
	maxGroupSize := 10
	
	totalTracks := len(tracks)
	for i := 0; i < totalTracks; i+=maxGroupSize {

		var trackGroup []db.Track
		if (i + maxGroupSize < totalTracks) {
			trackGroup = tracks[i:i+maxGroupSize]
		} else {
			trackGroup = tracks[i:]
		}

		_, err := db.InsertTracks(trackGroup)
		if err != nil { return err }
	}

	return nil
}

// USAGE: retrieves all of user's saved tracks
// PARAMS: 
//	- maxTracks = (x <= 0, for all tracks), (x>0, for track limit)
func getSavedTracks(client *http.Client, maxTracks int) ([]db.Track, error) {
	var tracks = make([]db.Track, 0)
	limit := 25

	if maxTracks < limit { limit = maxTracks }

	totalTracks := -1 // arbitrary default value
	var err error
	var trackResp spotify.TrackResponse
	var nextUrl *string = nil
	hasNextPage := true
	for hasNextPage && len(tracks) < maxTracks {
		if nextUrl != nil {
			trackResp, err = spotify.GetUserSavedTracksByURL(*nextUrl, client)
		} else {
			trackResp, err = spotify.GetUserSavedTracks(0, limit, client)
		}

		if totalTracks == -1 {
			totalTracks = trackResp.TotalItems
			if maxTracks == -1 { maxTracks = totalTracks }
		}

		if err != nil {
			log.Println("PopulateTable error: ", err)
			return tracks, err
		}

		for _, trackItem := range trackResp.Items {

			currTrack := db.Track{
				TrackID: trackItem.Track.ID,
				Name: trackItem.Track.Name,
				URI: trackItem.Track.URI,
				AlbumID: trackItem.Track.Album.ID,
			}

			tracks = append(tracks, currTrack)
		}

		nextUrl = trackResp.Next
		hasNextPage = nextUrl != nil
	}

	return tracks, nil
}