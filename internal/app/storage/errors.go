package storage

import "fmt"

type DBError struct {
	function string
	msg      string
	err      error
}

func (db *DBError) Error() string {
	return fmt.Sprintf("function: %s, msg: %s, err: %v", db.function, db.msg, db.err)
}

func NewDBError(function, msg string, err error) error {
	return &DBError{
		function: function,
		msg:      msg,
		err:      err,
	}
}
