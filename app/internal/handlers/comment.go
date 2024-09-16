package handlers

import (
	"app/internal/commons"
	"app/internal/entities"
	usecases "app/internal/usecases"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CommentHandler struct {
	usecases usecases.CommentUsecase
}

func NewCommentHandler(uc usecases.CommentUsecase) *CommentHandler {
	return &CommentHandler{usecases: uc}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	// retrieve postId from URL
	vars := mux.Vars(r)
	postIDStr := vars["id"]

	// Convert postId to uint
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		commons.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	var comment entities.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		commons.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	comment.PostID = uint(postID)

	res, err := h.usecases.CreateComment(r.Context(), &comment)
	if err != nil {
		status := http.StatusInternalServerError
		if err == commons.ErrForbidden {
			status = http.StatusForbidden
		} else if err == commons.ErrBadRequest {
			status = http.StatusBadRequest
		} else if err == commons.ErrNotFound {
			status = http.StatusNotFound
		}

		commons.ErrorResponse(w, status, err)
		return
	}

	commons.SuccessResponse(w, http.StatusCreated, res)
}

func (h *CommentHandler) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	// retrieve postId from URL
	vars := mux.Vars(r)
	postIDStr := vars["id"]

	// Convert postId to uint
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		commons.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	// Handle pagination (limit and page)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	comments, err := h.usecases.GetCommentsByPostID(r.Context(), uint(postID), limit, page)
	if err != nil {
		status := http.StatusInternalServerError
		if err == commons.ErrForbidden {
			status = http.StatusForbidden
		} else if err == commons.ErrBadRequest {
			status = http.StatusBadRequest
		} else if err == commons.ErrNotFound {
			status = http.StatusNotFound
		}

		commons.ErrorResponse(w, status, err)
		return
	}

	commons.SuccessResponse(w, http.StatusOK, comments)
}
