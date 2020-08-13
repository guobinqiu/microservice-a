package myerror

type Error struct {
	err        string
	statusCode int
}

func New(err string, statusCode int) *Error {
	return &Error{err, statusCode}
}

func (e *Error) Error() string {
	return e.err
}

func (e *Error) StatusCode() int {
	return e.statusCode
}
