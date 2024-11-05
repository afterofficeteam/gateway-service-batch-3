package helper

import (
	"encoding/json"
	model "gateway-service/models"
	"net/http"
)

const (
	SUCCESS_MESSSAGE string = "Success"
)

func HandleResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{
		Message: message,
		Data:    data,
	})
}
