package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/theskyinflames/quiz/csvreader/config"
	"github.com/theskyinflames/quiz/csvreader/pkg/domain"
)

type (
	Repository struct {
		cfg *config.Config
		db  *sql.DB
	}
)

func NewRepository(cfg *config.Config) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (r *Repository) Connect() (err error) {
	r.db, err = sql.Open("postgres", r.cfg.PostgreSQLConnStr)
	return r.db.Ping()
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) InsertBlock(block []domain.Record) error {
	ts := time.Now()

	txn, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("records", "id", "first_name", "last_name", "email", "phone"))
	if err != nil {
		return err
	}

	for _, record := range block {
		_, err = stmt.Exec(record.ID, record.FirstName, record.LastName, record.Email, record.Email)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	log.Printf("inserted %d block in %s ....\n", len(block), fmt.Sprint(time.Now().Sub(ts)))
	return nil
}

func (r *Repository) GetInsertedIDs() ([]string, error) {

	sqlStatement := `SELECT id FROM records`
	rows, err := r.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}
