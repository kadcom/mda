package todo

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrIsDone = errors.New("todo: the item is done")
)

type TodoItem struct {
	Id        ulid.ULID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	DoneAt    null.Time `json:"done_at"`
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

func (t TodoItem) MarshalJSON() ([]byte, error) {
	var j struct {
		Id        ulid.ULID  `json:"id"`
		Title     string     `json:"title"`
		CreatedAt time.Time  `json:"created_at"`
		DoneAt    *time.Time `json:"done_at,omitempty"`
		IsDone    bool       `json:"is_done"`
	}

	j.Id = t.Id
	j.Title = t.Title
	j.CreatedAt = t.CreatedAt
	j.DoneAt = t.DoneAt.Ptr()
	j.IsDone = t.IsDone()

	return json.Marshal(j)
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
