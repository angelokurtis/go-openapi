package openapi

import (
	"fmt"
	"log"
	"net/http"

	kin "github.com/getkin/kin-openapi/openapi3"
	kinrl "github.com/getkin/kin-openapi/routers/legacy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router extends `chi.Mux` to enable it to deal with OpenAPI 3 contract.
type Router struct {
	*chi.Mux
	operations operations
	validation validation
}

// NewRouter returns a new Router object.
func NewRouter(openapi []byte) *Router {
	mux := chi.NewMux()
	loader := kin.NewLoader()
	contract, err := loader.LoadFromData(openapi)
	if err != nil {
		panic(fmt.Sprintf("openapi: failed to load OpenAPI document: %s", err))
	}
	optns := readOperations(contract)
	// we need to bypass the servers in the contract to be able to find the router
	ccopy := &kin.T{
		ExtensionProps: contract.ExtensionProps,
		OpenAPI:        contract.OpenAPI,
		Components:     contract.Components,
		Info:           contract.Info,
		Paths:          contract.Paths,
		Security:       contract.Security,
		Tags:           contract.Tags,
		ExternalDocs:   contract.ExternalDocs,
	}
	finder, err := kinrl.NewRouter(ccopy)
	if err != nil {
		panic(fmt.Sprintf("openapi: failed to create route finder from OpenAPI document: %s", err))
	}
	v := &validator{routeFinder: finder}
	router := &Router{Mux: mux, operations: optns, validation: newValidation(v)}

	return router
}

// HandleOperation adds the route `pattern` that matches the `method` HTTP method described by `operationID` of
// OpenAPI 3 contract to execute the `handlerFn` http.HandlerFunc.
func (r *Router) HandleOperation(operationID string, handlerFn http.HandlerFunc) {
	optn, ok := r.operations[operationID]
	if !ok {
		panic(fmt.Sprintf("openapi: operation %q must be defined on OpenAPI document", operationID))
	}

	r.Mux.
		With(r.validation).
		With(middleware.AllowContentType(optn.contentTypes...)).
		MethodFunc(optn.method, optn.path, handlerFn)

	log.Printf("%s %s\n", optn.method, optn.path)
}
