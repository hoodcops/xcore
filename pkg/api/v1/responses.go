package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Error represents an API error code
// and it's associated human-readable message
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewMissingParamError ...
func NewMissingParamError(paramName string) Error {
	return Error{
		Code:    402,
		Message: fmt.Sprintf("Required param %s has empty value", paramName),
	}
}

// ErrorResponse is the response payload sent
// to clients when an error occurs during
// request handling.
type ErrorResponse struct {
	Summary string  `json:"summary"`
	Errors  []Error `json:"errors"`
}

// HasErrors ...
func (res *ErrorResponse) HasErrors() bool {
	return len(res.Errors) > 0
}

// NewErrorResponse ...
func NewErrorResponse(summary string) *ErrorResponse {
	return &ErrorResponse{
		Summary: summary,
		Errors:  []Error{},
	}
}

// AddError appends a new Error to an ErrorResponse
func (res *ErrorResponse) AddError(e Error) {
	res.Errors = append(res.Errors, e)
}

// NewInvalidPayloadResponse ...
func NewInvalidPayloadResponse(err error) ErrorResponse {
	return ErrorResponse{
		Summary: "Failed parsing request payload",
		Errors: []Error{
			{
				Code:    401,
				Message: err.Error(),
			},
		},
	}
}

// NewInternalServerErrorResponse ...
func NewInternalServerErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Summary: "Oops! something bad happened on the server. Please try again.",
		Errors: []Error{
			{
				Code:    500,
				Message: err.Error(),
			},
		},
	}
}

// OkResponse represent a response sent to
// clients when request is successful
type OkResponse struct {
	Data interface{} `json:"data"`
	Info string      `json:"info"`
}

func renderData(w http.ResponseWriter, payload interface{}) {
	renderJSON(w, http.StatusOK, payload)
}

func renderBadRequest(w http.ResponseWriter, payload interface{}) {
	renderJSON(w, http.StatusBadRequest, payload)
}

func renderInternalServerError(w http.ResponseWriter, payload interface{}) {
	renderJSON(w, http.StatusInternalServerError, payload)
}

func renderJSON(w http.ResponseWriter, status int, payload interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}
