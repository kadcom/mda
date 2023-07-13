//go:build fake

package todo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func findAllItems(ctx context.Context, tx pgx.Tx) (TodoList, error) {

	log.Debug().Msg("Fake find all item")

	list := TodoList{
		Items: fake_items,
		Count: len(fake_items),
	}

	return list, nil
}
