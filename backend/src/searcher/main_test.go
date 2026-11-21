package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	mux := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status to be 'ok', got %q", body["status"])
	}
}

func TestSearchBooksHandler_Success(t *testing.T) {
	mockYear := 1937
	mockCoverI := 12345
	mockResponse := openLibResponse{
		NumFound: 1,
		Docs: []openLibDoc{
			{
				Key:              "/works/OL27479W",
				Title:            "The Hobbit",
				AuthorName:       []string{"J.R.R. Tolkien"},
				FirstPublishYear: &mockYear,
				CoverI:           &mockCoverI,
				ISBN:             []string{"9780007440832"},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		author := r.URL.Query().Get("author")
		if author != "tolkien" {
			t.Errorf("expected author query param to be 'tolkien', got %q", author)
		}

		limit := r.URL.Query().Get("limit")
		if limit != "20" {
			t.Errorf("expected limit query param to be '20', got %q", limit)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	originalBase := openLibraryBase
	openLibraryBase = server.URL
	defer func() { openLibraryBase = originalBase }()

	mux := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/books?author=tolkien", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res SearchResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Total != 1 {
		t.Errorf("expected total 1, got %d", res.Total)
	}

	if len(res.Books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(res.Books))
	}

	book := res.Books[0]
	if book.Title != "The Hobbit" {
		t.Errorf("expected title 'The Hobbit', got %q", book.Title)
	}

	if len(book.AuthorNames) != 1 || book.AuthorNames[0] != "J.R.R. Tolkien" {
		t.Errorf("expected author 'J.R.R. Tolkien', got %v", book.AuthorNames)
	}

	expectedCoverURL := fmt.Sprintf("%s/%d-M.jpg", coverBase, mockCoverI)
	if book.CoverURL != expectedCoverURL {
		t.Errorf("expected cover URL %q, got %q", expectedCoverURL, book.CoverURL)
	}
}

func TestSearchBooksHandler_MissingAuthor(t *testing.T) {
	mux := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var errRes ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errRes); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errRes.Error != "author query param is required" {
		t.Errorf("expected error message to be 'author query param is required', got %q", errRes.Error)
	}
}
