package model

type Book struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Authors          []string `json:"authors"`
	FirstPublishYear int      `json:"first_publish_year,omitempty"`
	CoverURL         string   `json:"cover_url,omitempty"`
	ISBN             string   `json:"isbn"`
}
