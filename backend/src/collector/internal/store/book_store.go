package store

import (
	"book-search/backend/collector/internal/model"
	"database/sql"
	"strconv"
)

type Store interface {
	GetAll() ([]*model.Book, error)
	GetByID(id string) (*model.Book, error)
	Create(book *model.Book) (*model.Book, error)
	Update(id string, book *model.Book) (*model.Book, error)
	Delete(id string) error
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return &store{db: db}
}

func (s *store) GetAll() ([]*model.Book, error) {
	query := "SELECT id, title, authors, first_publish_year, cover_url, isbn FROM books"

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*model.Book
	for rows.Next() {
		book := &model.Book{}
		err := rows.Scan(&book.ID, &book.Title, &book.Authors, &book.FirstPublishYear, &book.CoverURL, &book.ISBN)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (s *store) GetByID(id string) (*model.Book, error) {
	query := "SELECT id, title, authors, first_publish_year, cover_url, isbn FROM books WHERE id = ?"
	row := s.db.QueryRow(query, id)
	book := &model.Book{}
	err := row.Scan(&book.ID, &book.Title, &book.Authors, &book.FirstPublishYear, &book.CoverURL, &book.ISBN)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // This is not an error, just means no book found with the given ID
		}
		return nil, err
	}
	return book, nil
}

func (s *store) Create(book *model.Book) (*model.Book, error) {
	query := "INSERT INTO books (title, authors, first_publish_year, cover_url, isbn) VALUES (?, ?, ?, ?, ?)"
	row, err := s.db.Exec(query, book.Title, book.Authors, book.FirstPublishYear, book.CoverURL, book.ISBN)
	if err != nil {
		return nil, err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}
	book.ID = strconv.FormatInt(id, 10)
	return book, nil
}

func (s *store) Update(id string, book *model.Book) (*model.Book, error) {
	query := "UPDATE books SET title=?, authors=?, first_publish_year=?, cover_url=?, isbn=? where id = ?"
	_, err := s.db.Exec(query, book.Title, book.Authors, book.FirstPublishYear, book.CoverURL, book.ISBN, id)
	if err != nil {
		return nil, err
	}
	book.ID = id
	return book, nil
}

func (s *store) Delete(id string) error {
	query := "DELETE FROM books WHERE id = ?"
	_, err := s.db.Exec(query, id)
	return err
}
