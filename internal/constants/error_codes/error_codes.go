package errorcodes

import "net/http"

type ErrorCodes int

const (
	InvalidCredentials ErrorCodes = 1000 + iota
	InvalidInput
	InternalServerError
)

var errToStr map[ErrorCodes]string = map[ErrorCodes]string{
	InvalidCredentials:  "Invalid credentials",
	InvalidInput:        "Invalid input, please check and try again",
	InternalServerError: "Internal server error",
}

var errToStatus map[ErrorCodes]int = map[ErrorCodes]int{
	InvalidCredentials:  http.StatusUnauthorized,
	InvalidInput:        http.StatusBadRequest,
	InternalServerError: http.StatusInternalServerError,
}

func (code ErrorCodes) String() string {
	return errToStr[code]
}

func (code ErrorCodes) Status() int {
	return errToStatus[code]
}
