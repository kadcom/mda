package todo

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool

	ErrTodoNotFound = errors.New("todo: not found")
)

func SetPool(newPool *pgxpool.Pool) error {

	if newPool == nil {
		return errors.New("cannot assign nil pool")
	}

	pool = newPool

	return nil
}
