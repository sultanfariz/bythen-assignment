package handlers

import (
	"app/internal/commons"
	"app/internal/entities"
	usecases "app/internal/usecases"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	usecases usecases.UserUsecase
}

func NewUserHandler(u usecases.UserUsecase) *UserHandler {
	return &UserHandler{usecases: u}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user entities.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		commons.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	res, err := h.usecases.Register(r.Context(), &user)
	if err != nil {
		status := http.StatusInternalServerError
		if err == commons.ErrUserAlreadyExists {
			status = http.StatusConflict
		}
		commons.ErrorResponse(w, status, err)
		return
	}

	commons.SuccessResponse(w, http.StatusCreated, res)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req entities.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commons.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	token, err := h.usecases.Login(r.Context(), &req)
	if err != nil {
		commons.ErrorResponse(w, http.StatusUnauthorized, err)
		return
	}
	commons.SuccessResponse(w, http.StatusOK, token)
}
