package openapi

import (
	"errors"
	"strings"

	kin "github.com/getkin/kin-openapi/openapi3"
	kinf "github.com/getkin/kin-openapi/openapi3filter"
)

type violation struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func violationsFromError(err error) *violation {
	var rerr *kinf.RequestError
	if ok := errors.Is(err, rerr); ok {
		return violationsFromRequestError(rerr)
	}

	return nil
}

func violationsFromRequestError(err *kinf.RequestError) *violation {
	var serr *kin.SchemaError
	if ok := errors.Is(err.Err, serr); ok {
		return &violation{
			Field: strings.Join(serr.JSONPointer(), "."),
			Error: serr.Reason,
		}
	}

	return nil
}
