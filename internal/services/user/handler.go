package user

import (
	"fmt"
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
	router.HandleFunc("/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/refresh", h.HandleRefreshToken).Methods("POST")

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

	accessToken, refreshToken, err := h.service.Register(r.Context(), payload.Email, payload.Password)
	if err != nil {
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
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"access_token": accessToken, "refresh_token": refreshToken, "message": "User registered successfully"})

}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	token, refreshToken, err := h.service.Login(r.Context(), payload.Email, payload.Password)
	if err == types.ErrInvalidCredentials {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err == types.ErrInternalServer {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Internal Server error"))
		utils.Error("An error occured during login process" + err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]string{
		"token":         token,
		"refresh_token": refreshToken,
	})

}

func (h *Handler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var payload types.RefreshTokenPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	token, refreshToken, err := h.service.RefreshToken(r.Context(), payload.RefreshToken)
	if err == types.ErrInvalidCredentials {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err == types.ErrInternalServer {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Internal Server error"))
		utils.Error("An error occured during login process" + err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusAccepted, map[string]string{
		"token":         token,
		"refresh_token": refreshToken,
	})

}
