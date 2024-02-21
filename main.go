package main

import (
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	// Handle "/login" route specifically
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/login.html")
	})

	http.ListenAndServe(":3000", nil)
}
