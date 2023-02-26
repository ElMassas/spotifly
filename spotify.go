package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
    token := os.Getenv("SPOTIFY_TOKEN") // Replace with your own Spotify access token

    req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/tracks?limit=50", nil)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }

    req.Header.Set("Authorization", "Bearer "+token)
    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return
    }

    defer resp.Body.Close()
    var tracks []struct {
        Track struct {
            Name string `json:"name"`
            ID   string `json:"id"`
        } `json:"track"`
    }

    err = json.NewDecoder(resp.Body).Decode(&tracks)
    if err != nil {
        fmt.Println("Error decoding response:", err)
        return
    }

    for _, track := range tracks {
        fmt.Printf("%s (%s)\n", track.Track.Name, track.Track.ID)
    }
}