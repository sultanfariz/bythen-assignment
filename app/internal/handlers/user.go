package handlers

import (
	"app/internal/commons"
	"app/internal/entities"
	usecases "app/internal/usecases"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	usecase usecases.UserUsecase
}

func NewUserHandler(u usecases.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user entities.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	res, err := h.usecase.Register(r.Context(), &user)
	if err != nil {
		status := http.StatusInternalServerError
		if err == commons.ErrUserAlreadyExists {
			status = http.StatusConflict
		}
		ErrorResponse(w, status, err)
		return
	}

	SuccessResponse(w, http.StatusCreated, res)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req entities.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	token, err := h.usecase.Login(r.Context(), &req)
	if err != nil {
		ErrorResponse(w, http.StatusUnauthorized, err)
		return
	}
	SuccessResponse(w, http.StatusOK, token)
}
