package spotify

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
)

func Login(w http.ResponseWriter, r *http.Request) {
	spotifyConfig := getSpotifyConfig()
	oauthState := generateStateCookie(w)
	u := spotifyConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	spotifyConfig := getSpotifyConfig()

	// get bearer token
	tokenData, err := spotifyConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Print(fmt.Errorf("code exchange error: %s", err.Error()))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	
	expDate := tokenData.Expiry.Format("2006-01-02 15:04:05.999999999 -0700 MST")

	// save token data to cookies
	exp := time.Now().Add(24 * time.Hour) // all tokens expires in 24 hours

	http.SetCookie(w, &http.Cookie{Name: "iffy_blue.AccessToken", Value: tokenData.AccessToken, 
		Expires: exp, SameSite: http.SameSiteLaxMode, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "iffy_blue.RefreshToken", Value: tokenData.RefreshToken, 
		Expires: exp, SameSite: http.SameSiteLaxMode, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "iffy_blue.TokenType", Value: tokenData.TokenType, 
		Expires: exp, SameSite: http.SameSiteLaxMode, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "iffy_blue.TokenExpiry", Value: expDate, 
		Expires: exp, SameSite: http.SameSiteLaxMode, Path: "/"})

	http.Redirect(w, r, "/print/successful-login", http.StatusPermanentRedirect)
}

func getSpotifyConfig() oauth2.Config {
	return oauth2.Config{
		RedirectURL:	"http://localhost:" + os.Getenv("PORT") + "/login/callback",
		ClientID:     	os.Getenv("CLIENT_ID"),
		ClientSecret: 	os.Getenv("CLIENT_SECRET"),
		Scopes:       	[]string{"playlist-read-private", "playlist-modify-private", 
							"playlist-modify-public", "user-library-read"},
		Endpoint:     	oauth2.Endpoint{
							AuthURL:  "https://accounts.spotify.com/authorize",
							TokenURL: "https://accounts.spotify.com/api/token"},
	}
}

func generateStateCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(24 * time.Hour) // 24 hours

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}
