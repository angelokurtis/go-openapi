package openapi

import (
	"fmt"
	"strings"

	kin "github.com/getkin/kin-openapi/openapi3"
)

type operations map[string]*operation

func readOperations(doc *kin.T) operations {
	ops := make(operations)
	for path, item := range doc.Paths {
		for method, optn := range item.Operations() {
			id := optn.OperationID
			body := optn.RequestBody
			ops[id] = newOperation(id, method, path, body, append(item.Parameters, optn.Parameters...))
		}
	}
	fmt.Println(ops.Template())

	return ops
}

func (o operations) Template() string {
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "r := openapi.NewRouter(specification)\n")
	for _, op := range o {
		_, _ = fmt.Fprintf(&b, "%s\n", op.Template())
	}
	b.WriteString("http.ListenAndServe(\":3000\", r)\n")

	return b.String()
}
