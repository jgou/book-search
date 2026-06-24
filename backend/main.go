package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	openLibraryBase = "https://openlibrary.org"
	coverBase       = "https://covers.openlibrary.org/b/id"
	defaultLimit    = 20
	maxLimit        = 100
)

type openLibDoc struct {
	Key              string   `json:"key"`
	Title            string   `json:"title"`
	AuthorName       []string `json:"author_name"`
	FirstPublishYear *int     `json:"first_publish_year"`
	CoverI           *int     `json:"cover_i"`
	ISBN             []string `json:"isbn"`
}

type openLibResponse struct {
	NumFound int          `json:"numFound"`
	Docs     []openLibDoc `json:"docs"`
}

type Book struct {
	Key              string   `json:"key"`
	Title            string   `json:"title"`
	AuthorNames      []string `json:"author_names"`
	FirstPublishYear *int     `json:"first_publish_year,omitempty"`
	CoverURL         string   `json:"cover_url,omitempty"`
	ISBN             []string `json:"isbn,omitempty"`
}

type SearchResponse struct {
	Books []Book `json:"books"`
	Total int    `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func searchBooksHandler(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	if author == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{"author query param is required"})
		return
	}

	limit := defaultLimit
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= maxLimit {
		limit = l
	}

	apiURL := fmt.Sprintf(
		"%s/search.json?author=%s&limit=%d&fields=key,title,author_name,first_publish_year,cover_i,isbn",
		openLibraryBase, url.QueryEscape(author), limit,
	)

	resp, err := http.Get(apiURL) //nolint:gosec
	if err != nil {
		log.Printf("open library request failed: %v", err)
		writeJSON(w, http.StatusBadGateway, ErrorResponse{"failed to reach Open Library"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{"failed to read upstream response"})
		return
	}

	var olResp openLibResponse
	if err := json.Unmarshal(body, &olResp); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{"failed to parse upstream response"})
		return
	}

	books := make([]Book, 0, len(olResp.Docs))
	for _, doc := range olResp.Docs {
		book := Book{
			Key:              doc.Key,
			Title:            doc.Title,
			AuthorNames:      doc.AuthorName,
			FirstPublishYear: doc.FirstPublishYear,
			ISBN:             doc.ISBN,
		}
		if doc.CoverI != nil {
			book.CoverURL = fmt.Sprintf("%s/%d-M.jpg", coverBase, *doc.CoverI)
		}
		books = append(books, book)
	}

	writeJSON(w, http.StatusOK, SearchResponse{Books: books, Total: olResp.NumFound})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/books", cors(searchBooksHandler))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
