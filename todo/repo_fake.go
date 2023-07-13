//go:build fake

package todo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func findItemById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (TodoItem, error) {

	log.Debug().Msg("Fake find item")

	var found bool
	var item TodoItem

	for _, v := range fake_items {
		if id == v.Id {
			item = v
			found = true
			break
		}
	}

	if !found {
		return TodoItem{}, ErrTodoNotFound
	}
	return item, nil
}

func saveItem(ctx context.Context, tx pgx.Tx, item TodoItem) error {

	log.Debug().Msg("Fake save item")

	var found bool

	for i, v := range fake_items {
		if item.Id == v.Id {
			fake_items[i] = item
			return nil
		}
	}

	if !found {
		fake_items = append(fake_items, item)
	}
	return nil

}
