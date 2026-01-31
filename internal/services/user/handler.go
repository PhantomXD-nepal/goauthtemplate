package user

import (
	"net/http"

	"github.com/PhantomXD-nepal/goauthtemplate/internal/types"
	"github.com/PhantomXD-nepal/goauthtemplate/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	service types.UserService
}

func NewHandler(service types.UserService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.HandleRegister).Methods("POST")

}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.service.Register(r.Context(), payload.Email, payload.Password)
	utils.Error("Failed to register user: " + err.Error())
	if err == types.ErrInternalServer {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err == types.ErrEmailAlreadyExists {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

}
