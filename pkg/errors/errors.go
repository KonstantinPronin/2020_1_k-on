package errors

type InvalidArgument struct {
	error string
}

func NewInvalidArgument(error string) *InvalidArgument {
	return &InvalidArgument{error: error}
}

func (ia *InvalidArgument) Error() string {
	return ia.error
}
