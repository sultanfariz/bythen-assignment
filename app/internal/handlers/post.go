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

type PostHandler struct {
	usecases usecases.PostUsecase
}

func NewPostHandler(uc usecases.PostUsecase) *PostHandler {
	return &PostHandler{usecases: uc}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post entities.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		commons.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.usecases.CreatePost(r.Context(), &post)
	if err != nil {
		commons.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	commons.SuccessResponse(w, http.StatusCreated, res)
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	posts, err := h.usecases.GetAllPosts(r.Context(), limit, page)
	if err != nil {
		commons.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	commons.SuccessResponse(w, http.StatusOK, posts)
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	// retrieve id from URL and pass it to usecase
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert id to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		commons.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	post, err := h.usecases.GetPostByID(r.Context(), uint(id))
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

	commons.SuccessResponse(w, http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	// retrieve id from URL and pass it to usecase
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert id to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		commons.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	var post entities.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		commons.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	post.ID = uint(id)

	res, err := h.usecases.UpdatePost(r.Context(), &post)
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

	commons.SuccessResponse(w, http.StatusOK, res)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	// retrieve id from URL and pass it to usecase
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert id to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		commons.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.usecases.DeletePost(r.Context(), uint(id))
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

	commons.SuccessResponse(w, http.StatusOK, "Post deleted successfully")
}
