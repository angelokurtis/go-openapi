package openapi

import (
	"fmt"
	"net/http"

	kinf "github.com/getkin/kin-openapi/openapi3filter"
	kinr "github.com/getkin/kin-openapi/routers"
	"github.com/go-chi/render"
)

type validation func(http.Handler) http.Handler

func newValidation(v *validator) validation {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			vlt, err := v.validate(req)
			if err != nil {
				fmt.Println(err.Error())
				res.WriteHeader(http.StatusInternalServerError)

				return
			}
			if vlt != nil {
				render.Status(req, http.StatusBadRequest)
				render.JSON(res, req, vlt)

				return
			}
			next.ServeHTTP(res, req)
		})
	}
}

type validator struct{ routeFinder kinr.Router }

func (v *validator) validate(req *http.Request) (*violation, error) {
	route, params, err := v.routeFinder.FindRoute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to find route: %w", err)
	}
	input := &kinf.RequestValidationInput{
		Request:    req,
		PathParams: params,
		Route:      route,
		Options:    &kinf.Options{ExcludeResponseBody: true, MultiError: true},
	}
	if err = kinf.ValidateRequest(req.Context(), input); err != nil {
		return violationsFromError(err), nil
	}

	return new(violation), nil
}
