package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// API CALLS

func GetUserSavedTracks(offset int, limit int, client *http.Client) (TrackResponse, error) {
	if limit < 0 {limit = 10 } // default to 10
	url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?limit=%d&offset=%d", limit, offset)

	return GetUserSavedTracksByURL(url, client)
}

func GetUserSavedTracksByURL(url string, client *http.Client) (TrackResponse, error) {
	var trackResp TrackResponse = TrackResponse{}
	
	resp, err := client.Get(url)
	if err != nil {
		log.Println("GetUserSavedTracks API error: ", err)
		return trackResp, err
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&trackResp)
	return trackResp, nil
}

// HELPERS
func CreateClient(r *http.Request) (*http.Client, error) {
	token, err := createToken(r)
	if err != nil {return http.DefaultClient, err}

	conf := getSpotifyConfig()
	return conf.Client(context.Background(), token), nil
}

func createToken(r *http.Request) (*oauth2.Token, error)  {
	accToken, err := r.Cookie("iffy_blue.AccessToken")
	if err != nil {return new(oauth2.Token), err}

	refToken, err := r.Cookie("iffy_blue.RefreshToken")
	if err != nil {return new(oauth2.Token), err}

	tokenType, err := r.Cookie("iffy_blue.TokenType")
	if err != nil {return new(oauth2.Token), err}

	exp, err := r.Cookie("iffy_blue.TokenExpiry")
	if err != nil {return new(oauth2.Token), err}
	expTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", exp.Value)
	if err != nil {return new(oauth2.Token), err}

	return &oauth2.Token{
		AccessToken: accToken.Value,
		RefreshToken: refToken.Value,
		TokenType: tokenType.Value,
		Expiry: expTime,
	}, nil
}

