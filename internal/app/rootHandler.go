package app

import (
	"fmt"
	"net/http"
)

type RootHandler struct{}

func NewRootHandler(router *http.ServeMux) {
	handler := &RootHandler{}
	router.HandleFunc("GET /", handler.apiReady)
}

func (h *RootHandler) apiReady(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API is ready. v1")
}
