package gopg

type DBConnError struct {
	errMsg string
}

func NewDBConnError(errMsg string) DBConnError {
	return DBConnError{errMsg}
}

func (db DBConnError) Error() string {
	return db.errMsg
}
