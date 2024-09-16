package commons

import (
	"encoding/json"
	"net/http"
)

type BaseResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	res := BaseResponse{
		Status:  status,
		Success: true,
		Message: "Success",
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func ErrorResponse(w http.ResponseWriter, status int, err error) {
	res := BaseResponse{
		Status:  status,
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}
