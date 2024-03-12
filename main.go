package main

import (
	//"context"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"

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
	url := "libsql://atlas-lawrencemueller.turso.io?authToken=[TOKEN]"

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
		queryMatches(db)
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
		// If this path is invoked the user wants to create a match

		// Get the name and create and ID for the match
		nameOfMatch := r.URL.Query().Get("playerName")
		idCreatedForMatch := generateRandomString(10)

		createMatch(db, nameOfMatch, idCreatedForMatch)

		// Create and add to data base a match with this name
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

	// Route to confirm lobby selection
	joinLobbyConfirmationComponent := joinLobbyConfirmation("test")
	http.Handle("/joinLobbyConfirmation", templ.Handler(joinLobbyConfirmationComponent))

	http.ListenAndServe(":3000", nil)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Function to insert match into data base
func createMatch(db *sql.DB, matchName string, matchID string) {
	_, err := db.Exec("INSERT INTO matches (id, name) VALUES (?, ?)", matchID, matchName)
	if err != nil {
		fmt.Println("Error inserting new match into database")
		os.Exit(1)
	}
}

// Struct to create a data structure for matches
type Match struct {
	ID   string
	Name string
}

// Function to get matches for lobby
func queryMatches(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM matches")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query get matches: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var matches []Match

	for rows.Next() {
		var match Match

		if err := rows.Scan(&match.ID, &match.Name); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		matches = append(matches, match)
		fmt.Print(matches)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}

	fmt.Print(matches)
}
