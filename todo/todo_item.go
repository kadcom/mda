package todo

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrIsDone = errors.New("todo: the item is done")
)

type TodoItem struct {
	Id        ulid.ULID
	Title     string
	CreatedAt time.Time
	DoneAt    null.Time
}

func (t TodoItem) IsDone() bool {
	return t.DoneAt.Valid && t.DoneAt.Time.After(t.CreatedAt)
}

func (t *TodoItem) MakeDone() error {
	if t.IsDone() {
		return ErrIsDone
	}

	t.DoneAt = null.TimeFrom(time.Now())
	return nil
}

func NewTodoItem(title string) (TodoItem, error) {
	if err := validateTitle(title); err != nil {
		return TodoItem{}, err
	}

	item := TodoItem{
		Id:        ulid.Make(),
		Title:     title,
		CreatedAt: time.Now(),
	}

	return item, nil
}
