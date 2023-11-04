package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"url-shortener/internal/model"

	_ "github.com/lib/pq"
)

type PgDatabase struct {
	db *sql.DB
}

func NewPgDatabase(dbPtr *sql.DB) *PgDatabase {
	return &PgDatabase{db: dbPtr}
}

func (pg *PgDatabase) AddEntry(ctx context.Context, originUrl, shortUrl string) error {
	const query = `INSERT INTO url_shorturl(origin_url, short_url)
				   VALUES ($1, $2);`
	if _, err := pg.db.ExecContext(ctx, query, originUrl, shortUrl); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") { // TODO sql.err
			switch {
			case strings.Contains(err.Error(), "_origin_url_key"):
				return model.ErrOriginConflict
			case strings.Contains(err.Error(), "_short_url_key"):
				return model.ErrShortConflict
			}
		}
		return fmt.Errorf("pg: %w", err)
	}
	return nil
}

func (pg PgDatabase) FindEntry(ctx context.Context, shortUrl string) (string, error) {
	const query = `SELECT origin_url FROM url_shorturl WHERE short_url = $1;`
	row := pg.db.QueryRowContext(ctx, query, shortUrl)
	var originUrl string
	if err := row.Scan(&originUrl); err != nil {
		return "", model.ErrShortBadParam
	}
	return originUrl, nil
}
