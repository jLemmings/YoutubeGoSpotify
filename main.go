package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistModifyPublic)
	ch    = make(chan *spotify.Client)
	state = "Abc123"
)

func main() {


	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)


	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

		list, err := client.GetPlaylistsForUser("in0cent")
	for _, playlist := range list.Playlists {
		fmt.Println("  ", playlist.Name, playlist.ID)
	}

	id, err := client.AddTracksToPlaylist("1pRYdwSAMA0CJjfeuujL89", "11dFghVXANMlKmJXsNCbNl")
	if err != nil {
		log.Fatal("Could not add Track to playlist", err)
	}
	fmt.Println(id)

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}