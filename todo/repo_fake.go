//go:build fake

package todo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var items []TodoItem

func FindAllItems(ctx context.Context, tx pgx.Tx) ([]TodoItem, error) {

	log.Debug().Msg("Fake find all item")
	return items, nil
}

func FindItemById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (TodoItem, error) {

	log.Debug().Msg("Fake find item")

	var found bool
	var item TodoItem

	for _, v := range items {
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

func SaveItem(ctx context.Context, tx pgx.Tx, item TodoItem) error {

	log.Debug().Msg("Fake save item")

	var found bool

	for i, v := range items {
		if item.Id == v.Id {
			items[i] = item
			return nil
		}
	}

	if !found {
		items = append(items, item)
	}
	return nil

}
