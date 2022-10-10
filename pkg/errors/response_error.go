package errors

import "net/http"

type Response interface {
	Error() string
	StatusCode() int
}

type ResponseImpl struct {
	statusCode int
	msg        string
}

func (h ResponseImpl) Error() string {
	return h.msg
}

func (h ResponseImpl) StatusCode() int {
	return h.statusCode
}

// generator function for the general http error
func NewResponseError(statusCode int, msg string) Response {
	return ResponseImpl{
		statusCode: statusCode,
		msg:        msg,
	}
}

// generator function of the common http errors
func NewResponseInternalServerError(msg string) Response {
	return NewResponseError(http.StatusInternalServerError, msg)
}

// generator function of the common http errors
func NewResponseBadRequestError(msg string) Response {
	return NewResponseError(http.StatusBadRequest, msg)
}

// generator function of the common http errors
func NewResponseForbiddenRequestError(msg string) Response {
	return NewResponseError(http.StatusForbidden, msg)
}

// generator function of the common http errors
func NewResponseNotFoundError(msg string) Response {
	return NewResponseError(http.StatusNotFound, msg)
}

// generator function of the common http errors
func NewResponseUnauthorizedError(msg string) Response {
	return NewResponseError(http.StatusUnauthorized, msg)
}

// generator function of the common http errors
func NewResponseExpectionFailedError(msg string) Response {
	return NewResponseError(http.StatusExpectationFailed, msg)
}

func NewResponseConflictError(msg string) Response {
	return NewResponseError(http.StatusConflict, msg)
}
