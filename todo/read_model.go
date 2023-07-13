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

type TodoList struct {
	Items []TodoItem `json:"items"`
	Count int        `json:"count"`
}

var emptyList TodoList

func findAllItems(ctx context.Context, tx pgx.Tx) (TodoList, error) {
	var itemCount int

	row := tx.QueryRow(ctx, "SELECT COUNT(id) as cnt FROM todolist;")
	err := row.Scan(&itemCount)

	if err != nil {
		log.Warn().Err(err).Msg("cannot find a count in todo list")
		return emptyList, err
	}

	if itemCount == 0 {
		return emptyList, nil
	}

	log.Debug().Int("count", itemCount).Msg("found todo items")

	items := make([]TodoItem, itemCount)

	rows, err := tx.Query(ctx, "SELECT id, title, created_at, done_at FROM todolist")

	if err != nil {
		return emptyList, err
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
			return emptyList, err
		}
		items[i] = TodoItem{
			Id: id, Title: title, CreatedAt: createdAt, DoneAt: doneAt,
		}
	}

	list := TodoList{
		Items: items,
		Count: itemCount,
	}

	return list, nil
}
