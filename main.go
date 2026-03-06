package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lib/pq"
)

type Cat struct {
	Id        int             `json:"id"`
	CatId     string          `json:"cat_id"`
	URL       string          `json:"url"`
	Width     int             `json:"width"`
	Height    int             `json:"height"`
	Breeds    json.RawMessage `json:"breeds"`
	ApiUsed   string          `json:"api_used"`
	CreatedAt *time.Time      `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
}

var db *sql.DB

func initDB() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("error pinging db: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Go API is running ok"})
}

func catsListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`SELECT id, cat_id, url, width, height, breeds, api_used, created_at, updated_at FROM cats`)
	if err != nil {
		log.Printf("error querying cats: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cats []Cat
	for rows.Next() {
		var c Cat
		var breedsRaw []byte
		var apiUsed sql.NullString
		if err := rows.Scan(&c.Id, &c.CatId, &c.URL, &c.Width, &c.Height, &breedsRaw, &apiUsed, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Printf("error scanning cat: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if apiUsed.Valid {
			c.ApiUsed = apiUsed.String
		}
		if len(breedsRaw) > 0 {
			c.Breeds = breedsRaw
		} else {
			c.Breeds = []byte("[]")
		}
		cats = append(cats, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cats)
}

func saveCatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		CatId  string          `json:"cat_id"`
		URL    string          `json:"url"`
		Width  int             `json:"width"`
		Height int             `json:"height"`
		Breeds json.RawMessage `json:"breeds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if input.Breeds == nil {
		input.Breeds = []byte("[]")
	}
	cat := Cat{
		CatId:   input.CatId,
		URL:     input.URL,
		Width:   input.Width,
		Height:  input.Height,
		Breeds:  input.Breeds,
		ApiUsed: "go",
	}
	err := db.QueryRow(
		`INSERT INTO cats (cat_id, url, width, height, breeds, api_used, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id, created_at, updated_at`,
		cat.CatId, cat.URL, cat.Width, cat.Height, cat.Breeds, cat.ApiUsed,
	).Scan(&cat.Id, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"detail": "Ya existe un gato con este ID en la base de datos."})
			return
		}
		log.Printf("error inserting cat: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cat)
}

func main() {
	initDB()

	http.HandleFunc("/health/", healthHandler)
	http.HandleFunc("/api/v1/cats/list/", catsListHandler)
	http.HandleFunc("/api/v1/cats/save/", saveCatHandler)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
