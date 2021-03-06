# Go OpenAPI

This project extends the [go-chi](https://github.com/go-chi/chi) router to support OpenAPI 3, bringing to you a simple
interface to build a router conforming your API contract.

## Examples

See [_examples/](https://github.com/angelokurtis/go-openapi/blob/main/_examples/) for more examples.

**As easy as:**

```go
package main

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/angelokurtis/go-openapi"
)

//go:embed openapi.yaml
var contract []byte

func main() {
	r := openapi.NewRouter(contract)
	r.Use(middleware.Logger)
	r.HandleOperation("getInvites", func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"operation": "getInvites"})
	})
	r.HandleOperation("createInvite", func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusCreated)
	})
	r.HandleOperation("deleteInvite", func(w http.ResponseWriter, r *http.Request) {
		_ = chi.URLParam(r, "inviteId")
		render.Status(r, http.StatusNoContent)
	})
	http.ListenAndServe(":3000", r)
}
```