//go:build !fake

package todo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

func FindAllItems(ctx context.Context, tx pgx.Tx) ([]TodoItem, error) {
	var itemCount int

	row := tx.QueryRow(ctx, "SELECT COUNT(id) as cnt FROM todolist;")
	err := row.Scan(&itemCount)

	if err != nil {
		log.Warn().Err(err).Msg("cannot find a count in todo list")
		return nil, err
	}

	if itemCount == 0 {
		return nil, nil
	}

	log.Debug().Int("count", itemCount).Msg("found todo items")

	items := make([]TodoItem, itemCount)

	rows, err := tx.Query(ctx, "SELECT id, title, created_at, done_at FROM todolist")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var i int

	for i = range items {
		var id ulid.ULID
		var title string
		var createdAt time.Time
		var doneAt null.Time

		if !rows.Next() {
			break
		}

		if err := rows.Scan(&id, &title, &createdAt, &doneAt); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return nil, err
		}
		items[i] = TodoItem{
			Id: id, Title: title, CreatedAt: createdAt, DoneAt: doneAt,
		}
	}
	return items, nil
}

func FindItemById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (TodoItem, error) {
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

func SaveItem(ctx context.Context, tx pgx.Tx, item TodoItem) error {
	q := `INSERT INTO todolist(id, title, created_at, done_at) VALUES ( $1, $2, $3, $4 )
        ON CONFLICT(id)
				DO UPDATE SET title=$2, done_at=$4`

	_, err := tx.Exec(ctx, q, item.Id, item.Title, item.CreatedAt, item.DoneAt)

	if err != nil {
		return err
	}

	return nil
}
