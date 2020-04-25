package errors

type InvalidArgumentError struct {
	error string
}

func NewInvalidArgument(error string) *InvalidArgumentError {
	return &InvalidArgumentError{error: error}
}

func (ia *InvalidArgumentError) Error() string {
	return ia.error
}

type DbInternalError struct {
	error string
}

func NewDbInternalError(error string) *DbInternalError {
	return &DbInternalError{error: error}
}

func (di *DbInternalError) Error() string {
	return di.error
}

type NotFoundError struct {
	error string
}

func NewNotFoundError(error string) *NotFoundError {
	return &NotFoundError{error: error}
}

func (nf *NotFoundError) Error() string {
	return nf.error
}

type ForbiddenError struct {
	error string
}

func NewForbiddenError(error string) *ForbiddenError {
	return &ForbiddenError{error: error}
}

func (f *ForbiddenError) Error() string {
	return f.error
}
