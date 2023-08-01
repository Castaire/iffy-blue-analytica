package core

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	
	spotify "iffy-blue-analytica/internal/spotify"
	db "iffy-blue-analytica/sql"
)

func PopulateTables(w http.ResponseWriter, r *http.Request) {
	client, err := spotify.CreateClient(r)
	if err != nil {
		log.Println("Could not create client: ", err)
		return
	}

	maxTracks, err := strconv.Atoi(chi.URLParam(r, "maxTracks"))
	if err != nil {
		// arbitrary negative value; ALL tracks will be retrieved
		maxTracks = -1
	}
	tracks, err := getSavedTracks(client, maxTracks)

	if err != nil {
		log.Println("Error in getting saved tracks: ", err)
		return
	}

	err = uploadTracksToDB(tracks)
	if err != nil {
		log.Println("Error in uploading tracks: ", err)
		return
	}
}

func ResetTables(w http.ResponseWriter, r *http.Request) {
	err := db.CreateTables(true)

	if err != nil {
		log.Println("Error in resetting tables: ", err)
		return
	}
}