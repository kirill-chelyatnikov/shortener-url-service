package app

import "fmt"

type DBErrors struct {
	function string
	err      error
}

func (db *DBErrors) Error() string {
	return fmt.Sprintf("[function: %s] %v", db.function, db.err)
}

func NewDBError(function string, err error) error {
	return &DBErrors{
		function: function,
		err:      err,
	}
}
