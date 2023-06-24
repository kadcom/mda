package todo

import (
	"context"

	"github.com/oklog/ulid/v2"
)

func ListItems(ctx context.Context) ([]TodoItem, error) {
	tx, err := pool.Begin(ctx)

	if err != nil {
		return nil, err
	}

	items, err := FindAllItems(ctx, tx)

	if err != nil {
		return nil, err
	}

	tx.Commit(ctx)

	return items, nil
}

func CreateItem(ctx context.Context, title string) (id ulid.ULID, err error) {
	todoItem, err := NewTodoItem(title)

	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)

	if err != nil {
		return
	}

	err = SaveItem(ctx, tx, todoItem)

	if err != nil {
		tx.Rollback(ctx)
		return
	}

	err = tx.Commit(ctx)

	if err != nil {
		return
	}

	return todoItem.Id, nil
}

func FindItem(ctx context.Context, id ulid.ULID) (item TodoItem, err error) {
	tx, err := pool.Begin(ctx)

	if err != nil {
		return
	}

	item, err = FindItemById(ctx, tx, id)

	if err != nil {
		return TodoItem{}, err
	}

	err = tx.Commit(ctx)
	return
}

func MakeItemDone(ctx context.Context, id ulid.ULID) error {
	tx, err := pool.Begin(ctx)

	if err != nil {
		return err
	}

	item, err := FindItemById(ctx, tx, id)

	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err = item.MakeDone(); err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err = SaveItem(ctx, tx, item); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
