package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
    // Replace with your own Tidal access token and user ID
    token := os.Getenv("TIDAL_TOKEN")
    userID := os.Getenv("TIDAL_USER_ID")

    // Read the song names from the JSON file
    file, err := os.Open("songs.json")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    var songNames []string
    err = json.NewDecoder(file).Decode(&songNames)
    if err != nil {
        fmt.Println("Error decoding JSON:", err)
        return
    }

    // Like each song in the list
    for _, songName := range songNames {
        songID, err := getSongID(songName, token)
        if err != nil {
            fmt.Printf("Error getting ID for song '%s': %v\n", songName, err)
            continue
        }
        err = likeSong(songID, userID, token)
        if err != nil {
            fmt.Printf("Error liking song '%s': %v\n", songName, err)
            continue
        }
        fmt.Printf("Liked song '%s'\n", songName)
    }
}

func getSongID(songName string, token string) (string, error) {
    // Search for the song by name
    url := fmt.Sprintf("https://api.tidal.com/v1/search?query=%s&limit=1&offset=0&types=TRACKS", songName)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", fmt.Errorf("error creating request: %v", err)
    }
    req.Header.Set("Authorization", "Bearer "+token)

    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    var searchResult struct {
        Tracks struct {
            Items []struct {
                ID string `json:"id"`
            } `json:"items"`
        } `json:"tracks"`
    }

    err = json.NewDecoder(resp.Body).Decode(&searchResult)
    if err != nil {
        return "", fmt.Errorf("error decoding response: %v", err)
    }

    if len(searchResult.Tracks.Items) == 0 {
        return "", fmt.Errorf("song not found")
    }

    return searchResult.Tracks.Items[0].ID, nil
}

func likeSong(songID string, userID string, token string) error {
    // Like the song by adding it to the user's favorites
    url := fmt.Sprintf("https://api.tidal.com/v1/users/%s/favorites/tracks/%s", userID, songID)
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return fmt.Errorf("error creating request: %v", err)
    }
    req.Header.Set("Authorization", "Bearer "+token)

    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return nil
}