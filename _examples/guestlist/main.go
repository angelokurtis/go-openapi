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
