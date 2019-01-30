package main

import (
	"encoding/json"
	"net/http"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
	skills := fetchNeoGraph()
	js, err := json.Marshal(skills)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
