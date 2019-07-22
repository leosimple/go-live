package restfulapi

import (
	"encoding/json"
	"go-live/models"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AppsResponse struct {
	Code    int          `json:"code"`
	Data    []models.App `json:"data"`
	Message string       `json:"message"`
}

type AppResponse struct {
	Code    int         `json:"code"`
	Data    *models.App `json:"data"`
	Message string      `json:"message"`
}

type LiveTokenResponse struct {
	Code    int    `json:"code"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type LivesResponse struct {
	Code    int           `json:"code"`
	Data    []models.Live `json:"data"`
	Message string        `json:"message"`
}

type LiveResponse struct {
	Code    int         `json:"code"`
	Data    models.Live `json:"data"`
	Message string      `json:"message"`
}

type StatusResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	PlayerCount int    `json:"playercount"`
	IsPublisher bool   `json:"ispublisher"`
}

func SendErrorResponse(w http.ResponseWriter, code int, message string) {
	SendResponse(w, code, &ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func SendResponse(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data, err := json.Marshal(v); err == nil {
		w.Write(data)
	} else {
		w.Write([]byte("Error."))
	}
}
