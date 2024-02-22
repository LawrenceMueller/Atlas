package main

import (
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	// Handle "/lobby" route specifically
	http.HandleFunc("/lobby", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/lobby.html")
	})

	// Handle "/instructions" route specifically
	http.HandleFunc("/instructions", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/instructions.html")
	})

	http.ListenAndServe(":3000", nil)
}
