package Error

type CustomError struct {
	err       string
	userError string
}

func NewCustomError(err string, userError string) *CustomError {
	return &CustomError{err: err, userError: userError}
}

func (e *CustomError) Error() string {
	return e.err
}

func (ue *CustomError) UserError() string {
	return ue.userError
}
