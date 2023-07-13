package todo

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

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

func parseNullStringToNullTime(s null.String) (t null.Time) {
	if !s.Valid {
		return
	}

	ts, err := time.Parse(time.RFC3339, s.String)

	if err != nil {
		return
	}

	return null.TimeFrom(ts)
}

func (t *TodoItem) UnmarshalJSON(data []byte) error {
	var j struct {
		Id        ulid.ULID   `json:"id"`
		Title     string      `json:"title"`
		CreatedAt string      `json:"created_at"`
		DoneAt    null.String `json:"done_at"`
	}

	err := json.Unmarshal(data, &j)

	if err != nil {
		return err
	}

	createdAt, err := time.Parse(time.RFC3339, j.CreatedAt)

	if err != nil {
		return err
	}

	doneAt := parseNullStringToNullTime(j.DoneAt)

	t = &TodoItem{
		Id:        j.Id,
		Title:     j.Title,
		CreatedAt: createdAt,
		DoneAt:    doneAt,
	}

	return nil
}
