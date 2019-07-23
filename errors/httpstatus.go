package errors

import "net/http"

// HTTPStatus sets an http status for the error
type HTTPStatus int

// Int returns the HTTP status
func (s HTTPStatus) Int() int {
	return int(s)
}

func (s HTTPStatus) String() string {
	return http.StatusText(int(s))
}

// GetHTTPStatus returns the HTTPStatus of an error.
func GetHTTPStatus(err error) HTTPStatus {
	for {
		e, ok := err.(Error)
		if !ok {
			break
		}
		if e.HTTPStatus != 0 {
			return e.HTTPStatus
		}
		err = e.Err
	}
	return 0
}
