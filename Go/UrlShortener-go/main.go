package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const domain = "http://localhost:3001/"

type URLRecord struct {
	LongURL  string
	ShortURL string
	Expiry   time.Time
}

type urlStore struct {
	mu           sync.RWMutex
	db           map[string]URLRecord
	counterStore map[string]int
}

var store = &urlStore{
	db:           make(map[string]URLRecord),
	counterStore: map[string]int{"UNIQUE_COUNTER": 0},
}

var shortCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9]{6}$`)

func main() {
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/", handleRedirect)
	fmt.Println("Server started at port 3001")
	http.ListenAndServe(":3001", nil)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		URL    string `json:"url"`
		Expiry *int   `json:"expiry"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"success": false, "errors": []string{"invalid body"}})
		return
	}

	body.URL = strings.TrimSpace(body.URL)
	if body.URL == "" {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"success": false, "errors": []string{"URL must not be empty!"}})
		return
	}
	parsed, err := url.Parse(body.URL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"success": false, "errors": []string{"URL must contain valid protocol!"}})
		return
	}
	if len(body.URL) > 2000 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"success": false, "errors": []string{"URL is too long!"}})
		return
	}

	expiryDays := 30
	if body.Expiry != nil {
		if *body.Expiry < 1 || *body.Expiry > 365 {
			writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"success": false, "errors": []string{"expiry must be between 1 and 365 days"}})
			return
		}
		expiryDays = *body.Expiry
	}

	store.mu.Lock()
	shortCode := generateShortCode(store.counterStore)
	shortURL := domain + shortCode
	store.db[shortCode] = URLRecord{
		LongURL:  body.URL,
		ShortURL: shortURL,
		Expiry:   time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour),
	}
	store.mu.Unlock()
	writeJSON(w, http.StatusCreated, map[string]any{"success": true, "shortUrl": shortURL})
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	if shortCode == "" || !shortCodeRegex.MatchString(shortCode) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, "<html><head><title>URL Shortener</title></head><body><h1>Invalid URL ;)</h1></body></html>")
		return
	}
	store.mu.RLock()
	record, ok := store.db[shortCode]
	store.mu.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<html><head><title>URL Shortener</title></head><body><h1>URL not found ;)</h1></body></html>")
		return
	}
	if record.Expiry.Before(time.Now()) {
		store.mu.Lock()
		delete(store.db, shortCode)
		store.mu.Unlock()
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<html><head><title>URL Shortener</title></head><body><h1>URL Expired ;)</h1></body></html>")
		return
	}
	http.Redirect(w, r, record.LongURL, http.StatusFound)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
