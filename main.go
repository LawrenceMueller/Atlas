package main

import (
	//"context"
	"database/sql"
	"fmt"
	"github.com/a-h/templ"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {

	// Connect to DB
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Fetch the token from environment variable
	authToken := os.Getenv("TURSO_DB_TOKEN")

	// URL string with placeholder
	url := "libsql://atlas.turso.io?authToken=[TOKEN]"

	// Replace the placeholder with the actual token
	url = strings.Replace(url, "[TOKEN]", authToken, 1)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	defer db.Close()

	// Web Server stuff

	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	// Handle "/lobby" route specifically
	http.HandleFunc("/lobby", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/lobby.html")
	})

	// Handle "/createMatch" route specifically
	http.HandleFunc("/createMatch", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/createMatch.html")
	})

	// Handle "/instructions" route specifically
	http.HandleFunc("/instructions", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/instructions.html")
	})

	// Handle "/waitingRoom" route specifically
	http.HandleFunc("/waitingRoom", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/waitingRoom.html")
	})

	// test data
	testData := []string{
		"Lobby 1: Casual Games",
		"Lobby 2: Competitive Matches",
		"Lobby 3: Strategy Games",
		"Lobby 4: Racing Games",
		"Lobby 5: Adventure Games",
		"Lobby 6: Puzzle Games",
		"Lobby 7: Shooter Games",
		"Lobby 8: Sports Games",
		"Lobby 9: Multiplayer Online Battle Arena",
		"Lobby 10: Role-Playing Games",
		"Lobby 11: Simulation Games",
		"Lobby 12: Party Games",
		"Lobby 13: Virtual Reality Games",
		"Lobby 14: Augmented Reality Games",
		"Lobby 15: Indie Games"}

	// Route to get lobby data
	lobbyDataComponent := lobbyData(testData)
	http.Handle("/lobbyData", templ.Handler(lobbyDataComponent))
	//lobbyDataComponent.Render(context.Background(), os.Stdout)

	// Route to confirm lobby selection
	joinLobbyConfirmationComponent := joinLobbyConfirmation("test")
	http.Handle("/joinLobbyConfirmation", templ.Handler(joinLobbyConfirmationComponent))

	http.ListenAndServe(":3000", nil)
}
