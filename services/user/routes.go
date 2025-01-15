package user

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.Handlelogin).Methods("POST")
	router.HandleFunc("/register", h.HandleRegister).Methods("POST")
}

func (h *Handler) Handlelogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {

}
