package exception

import (
	"encoding/json"
	"net/http"
)

// ResponseError adalah struktur standar untuk error response
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// WriteErrorResponse menulis error response dalam format JSON
func WriteErrorResponse(writer http.ResponseWriter, statusCode int, message string) {
	response := ResponseError{
		Code:    statusCode,
		Message: message,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	json.NewEncoder(writer).Encode(response)
}

// NotFoundHandler menangani error 404 (route tidak ditemukan)
func NotFoundHandler(writer http.ResponseWriter, request *http.Request) {
	WriteErrorResponse(writer, http.StatusNotFound, "Route not found")
}

// MethodNotAllowedHandler menangani error 405 (method tidak diizinkan)
func MethodNotAllowedHandler(writer http.ResponseWriter, request *http.Request) {
	WriteErrorResponse(writer, http.StatusMethodNotAllowed, "Method not allowed")
}

// InternalServerErrorHandler menangani error 500 (kesalahan server)
func InternalServerErrorHandler(writer http.ResponseWriter, err error) {
	WriteErrorResponse(writer, http.StatusInternalServerError, "Internal server error: "+err.Error())
}
