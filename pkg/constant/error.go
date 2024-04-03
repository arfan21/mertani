package constant

import (
	"errors"
	"net/http"
)

const (
	ErrSQLUniqueViolation = "23505"
	ErrSQLInvalidUUID     = "22P02"
	ErrSQLFKViolation     = "23503"
)

var (
	ErrEmailAlreadyRegistered = &ErrWithCode{HTTPStatusCode: http.StatusConflict, Message: "email already registered"}
	ErrEmailOrPasswordInvalid = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "email or password invalid"}
	ErrInvalidUUID            = errors.New("invalid uuid length or format")
	ErrUnauthorizedAccess     = &ErrWithCode{HTTPStatusCode: http.StatusUnauthorized, Message: "unauthorized access"}
	ErrDeviceNotFound         = &ErrWithCode{HTTPStatusCode: http.StatusNotFound, Message: "device not found"}
	ErrSensorNotFound         = &ErrWithCode{HTTPStatusCode: http.StatusNotFound, Message: "sensor not found"}
)

type ErrWithCode struct {
	HTTPStatusCode int
	Message        string
}

func (e *ErrWithCode) Error() string {
	return e.Message
}

type ErrValidation struct {
	Message string
}

func (e *ErrValidation) Error() string {
	return e.Message
}
