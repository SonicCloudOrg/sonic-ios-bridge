package util

import "fmt"

const (
	ErrConnect     = "failed connecting to "
	ErrReadingMsg  = "failed reading msg "
	ErrSendCommand = "failed send the command "
	ErrMissingArgs = "missing arg(s)"
	ErrUnknown = "unknown error"
)

func NewErrorPrint(t string, msg string, err error) error {
	if len(msg) == 0 && err == nil {
		return fmt.Errorf("%s", t)
	}
	if len(msg) == 0 {
		return fmt.Errorf("%s : %w", t, err)
	}
	if err == nil {
		return fmt.Errorf("%s %s", t, msg)
	}
	return fmt.Errorf("%s %s : %w", t, msg, err)
}
