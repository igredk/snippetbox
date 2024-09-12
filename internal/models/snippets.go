package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires)
    VALUES($1, $2, NOW(), NOW() + INTERVAL '1 day' * $3) RETURNING id`

	var id int
	err := m.DB.QueryRow(context.Background(), query, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	if id < 1 {
		return nil, ErrNoRecord
	}

	query := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > NOW() AND id = $1`

	rows, _ := m.DB.Query(context.Background(), query, id)
	defer rows.Close()

	snippet, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Snippet])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > NOW() ORDER BY id DESC LIMIT 10`

	rows, _ := m.DB.Query(context.Background(), query)
	defer rows.Close()

	snippets, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Snippet])
	if err != nil {
		return nil, err
	}

	return snippets, nil
}
