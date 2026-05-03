package metrics

type ErrorType string

const (
	ErrBadRequest ErrorType = "bad_request_error"
	ErrForbidden  ErrorType = "forbidden_error"
	ErrInternal   ErrorType = "internal_error"
)

func (et ErrorType) String() string {
	return string(et)
}
