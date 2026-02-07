package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	r.Group(func(pr chi.Router) {
		pr.Use(AuthMiddleware)
		pr.Get("/videos/{videoID}", h.GetVideo)
	})

	return r
}
