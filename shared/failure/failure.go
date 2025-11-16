package failure

import (
	"errors"
	"fmt"
	"net/http"
)

type Failure struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Failure) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func GetFailure(err error) Failure {
	if err == nil {
		return Failure{}
	}

	failureErr := &Failure{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	errors.As(err, &failureErr)

	return *failureErr
}

func SetFailure(err error, code int, msg string) error {
	if err == nil {
		return nil
	}

	return WrapE(err, &Failure{Code: code, Message: msg})
}

func GetCode(err error) int {
	return GetFailure(err).Code
}

func SetCode(err error, code int) error {
	failureErr := GetFailure(err)
	return SetFailure(err, code, failureErr.Message)
}

func GetMsg(err error) string {
	return GetFailure(err).Message
}

func SetMsg(err error, msg string) error {
	failureErr := GetFailure(err)
	return SetFailure(err, failureErr.Code, msg)
}

// ------------------------------------ wrapping error ------------------------------------ //

func BadRequest(err error) error {
	if err != nil {
		return SetFailure(err, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	return nil
}

func InternalError(err error) error {
	if err != nil {
		return SetFailure(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func Unauthorized(err error) error {
	if err != nil {
		return SetFailure(err, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}
	return nil
}

// ------------------------------------ new error ------------------------------------ //

func BadRequestFromString(msg string) error {
	return &Failure{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

func NotFound(domainName string) error {
	return &Failure{
		Code:    http.StatusNotFound,
		Message: domainName,
	}
}

func Conflict(operationName string, domainName string, message string) error {
	return &Failure{
		Code:    http.StatusConflict,
		Message: fmt.Sprintf("%s on %s: %s", operationName, domainName, message),
	}
}
