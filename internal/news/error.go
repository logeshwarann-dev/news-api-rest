package news

type CustomError struct {
	err        error
	httpStatus int
}

func NewCustomError(err error, httpStatus int) error {
	return &CustomError{
		err:        err,
		httpStatus: httpStatus,
	}
}

func (ce CustomError) Error() string {
	return ce.err.Error()
}

func (ce CustomError) Unwrap() error {
	return ce.err
}

func (ce CustomError) GetHttpStatus() int {
	return ce.httpStatus
}
