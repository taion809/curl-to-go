package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moonltg/curl-to-go/converter"
)

type ConvertRequest struct {
	CurlCommand string `json:"curl_command"`
}

type ConvertResponse struct {
	GoCode string `json:"go_code"`
	Error  string `json:"error,omitempty"`
}

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/", serveHome)
	r.Post("/api/convert", handleConvert)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	goCode, err := converter.CurlToGo(req.CurlCommand)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := ConvertResponse{
		GoCode: goCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ConvertResponse{
		Error: message,
	})
}
