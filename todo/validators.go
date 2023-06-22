package todo

import "errors"

var (
	ErrTitleTooLong  = errors.New("todo: title too long")
	ErrTitleTooShort = errors.New("todo: title too short")
	ErrTitleEmpty    = errors.New("todo: title empty")
)

const minTitle = 5
const maxTitle = 1000

func validateTitle(title string) error {
	l := len(title)

	switch {
	case l == 0:
		return ErrTitleEmpty
	case l < minTitle:
		return ErrTitleTooShort
	case l > maxTitle:
		return ErrTitleTooLong
	default:
		return nil
	}
}
