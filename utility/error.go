package utility

type InternalServerError struct {
	err error
}

type ValidationError struct {
	err     error
	errType string
	field   string
}

type ResourceNotFoundError struct {
	err error
	ID  interface{}
}

func (e *InternalServerError) Error() string {
	return e.err.Error()
}

func (e *ValidationError) Error() string {
	return e.err.Error()
}

func (e *ResourceNotFoundError) Error() string {
	return e.err.Error()
}

func NewInternalServerError(err error) error {
	return &InternalServerError{err: err}
}

func NewValidationError(err error, errType string, field string) error {
	return &ValidationError{err: err, errType: errType, field: field}
}

func NewResourceNotFoundError(err error, id interface{}) *ResourceNotFoundError {
	return &ResourceNotFoundError{
		err: err,
		ID:  id,
	}
}
