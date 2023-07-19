package storage

import "fmt"

type DBErrors struct {
	function string
	msg      string
	err      error
}

func (db *DBErrors) Error() string {
	return fmt.Sprintf("function: %s, msg: %s, err: %v", db.function, db.msg, db.err)
}

func NewDBError(function, msg string, err error) error {
	return &DBErrors{
		function: function,
		msg:      msg,
		err:      err,
	}
}
