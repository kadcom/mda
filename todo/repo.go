//go:build !fake

package todo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func findItemById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (TodoItem, error) {
	q := `SELECT id, title, created_at, done_at FROM todolist WHERE id = $1`

	row := tx.QueryRow(ctx, q, id)

	var item TodoItem
	if err := row.Scan(&item.Id, &item.Title, &item.CreatedAt, &item.DoneAt); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return TodoItem{}, ErrTodoNotFound
		}
		return TodoItem{}, err
	}

	return item, nil
}

func saveItem(ctx context.Context, tx pgx.Tx, item TodoItem) error {
	q := `INSERT INTO todolist(id, title, created_at, done_at) VALUES ( $1, $2, $3, $4 )
        ON CONFLICT(id)
				DO UPDATE SET title=$2, done_at=$4`

	_, err := tx.Exec(ctx, q, item.Id, item.Title, item.CreatedAt, item.DoneAt)

	if err != nil {
		return err
	}

	return nil
}
