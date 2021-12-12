package openapi

import (
	"fmt"
	"strings"

	kin "github.com/getkin/kin-openapi/openapi3"
)

type operation struct {
	id, method, path                      string
	queryParams, pathParams, contentTypes []string
}

func newOperation(id, method, path string, body *kin.RequestBodyRef, params kin.Parameters) *operation {
	o := &operation{id: id, method: method, path: path}
	o.addContentTypes(body)
	o.addParams(params)

	return o
}

func (o *operation) addParams(params kin.Parameters) {
	for _, parameter := range params {
		if parameter.Value == nil {
			continue
		}
		switch parameter.Value.In {
		case "path":
			o.pathParams = dedupe(o.pathParams, parameter.Value.Name)
		case "query":
			o.pathParams = dedupe(o.queryParams, parameter.Value.Name)
		}
	}
}

func (o *operation) addContentTypes(body *kin.RequestBodyRef) {
	contentTypes := make([]string, 0)
	if body != nil && body.Value != nil {
		for c := range body.Value.Content {
			contentTypes = append(contentTypes, c)
		}
	}
	o.contentTypes = contentTypes
}

func (o *operation) Template() string {
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "r.HandleOperation(\"%s\", func(w http.ResponseWriter, r *http.Request) {\n", o.id)
	for _, path := range o.pathParams {
		_, _ = fmt.Fprintf(&b, "\t_ = chi.URLParam(r, \"%s\")\n", path)
	}
	for _, query := range o.queryParams {
		_, _ = fmt.Fprintf(&b, "\t_ = r.URL.Query().Get(\"%s\")\n", query)
	}
	switch o.method {
	case "POST":
		_, _ = fmt.Fprintf(&b, "\trender.Status(r, http.StatusCreated)\n")
	case "PUT", "PATCH":
		_, _ = fmt.Fprintf(&b, "\trender.Status(r, http.StatusAccepted)\n")
	case "DELETE":
		_, _ = fmt.Fprintf(&b, "\trender.Status(r, http.StatusNoContent)\n")
	default:
		_, _ = fmt.Fprintf(&b, "\trender.Status(r, http.StatusOK)\n")
		_, _ = fmt.Fprintf(&b, "\trender.JSON(w, r, map[string]string{\"operation\": \"%s\"})\n", o.id)
	}
	b.WriteString("})")

	return b.String()
}

func dedupe(a []string, b ...string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}
	for letter := range check {
		res = append(res, letter)
	}

	return res
}
