package repository

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	ErrNotFound       error = errors.New("resource not found")
	ErrInternalServer error = errors.New("internal server error")
	ErrBadRequest     error = errors.New("bad request error")
	ErrIdk            error = errors.New("idk what's happened")
)

var mappingError map[int]error = map[int]error{
	http.StatusBadRequest:          ErrBadRequest,
	http.StatusInternalServerError: ErrInternalServer,
	http.StatusNotFound:            ErrNotFound,
}

func readAndUnmarshal[T any](body io.Reader, model *T) error {
	buffBody, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buffBody, model); err != nil {
		return err
	}
	return nil
}

func errorCodeMapping(errorCode int) error {
	if err, ok := mappingError[errorCode]; !ok {
		return ErrIdk
	} else {
		return err
	}
}

func treatResult(response *http.Response, expectedReturnCode int) error {
	if returnCode := response.StatusCode; expectedReturnCode != returnCode {
		return errorCodeMapping(returnCode)
	}
	return nil
}
