package shortener

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type sqlSlugToURL struct {
	db *sql.DB
}

func NewSQLSlugToURL(dbName string) (SlugToURL, error) {
	db, dbOpenErr := sql.Open("sqlite3", dbName)
	if dbOpenErr != nil {
		return nil, dbOpenErr
	}
	return &sqlSlugToURL{db}, nil
}

func (s *sqlSlugToURL) URLExists(slug string) (bool, error) {
	// Write SQL query
	query := fmt.Sprintf("SELECT url FROM Mappings WHERE slug = '%s'", slug)
	rows, sqlQueryErr := s.db.Query(query)
	if sqlQueryErr != nil {
		return false, sqlQueryErr
	}

	// Retrieve result from rows (there should be only one row,
	// so we use if rows.Next(), not for rows.Next()).
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (s *sqlSlugToURL) URL(slug string) (string, error) {
	// Write SQL query
	query := fmt.Sprintf("SELECT url FROM Mappings WHERE slug = '%s'", slug)
	rows, sqlQueryErr := s.db.Query(query)
	if sqlQueryErr != nil {
		return "", sqlQueryErr
	}

	// Retrieve result from rows (there should be only one row,
	// so we use if rows.Next(), not for rows.Next()).
	defer rows.Close()
	url := ""
	if rows.Next() {
		scanErr := rows.Scan(&url)
		if scanErr != nil {
			return url, scanErr
		}
	}
	return url, nil
}
