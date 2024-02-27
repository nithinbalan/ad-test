package utility

type InternalServerError struct {
	err error
}

func (e *InternalServerError) Error() string {
	return e.err.Error()
}

func NewInternalServerError(err error) error {
	return &InternalServerError{err: err}
}
