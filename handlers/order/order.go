package order

import (
	"encoding/json"
	model "gateway-service/models"
	"gateway-service/usecases/order"
	"gateway-service/util/helper"
	"gateway-service/util/middleware"
	"net/http"

	"github.com/go-playground/validator"
)

type Handler struct {
	validator *validator.Validate
}

func NewHandler(validate *validator.Validate) *Handler {
	return &Handler{
		validator: validate,
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		helper.HandleResponse(w, http.StatusBadRequest, "User ID cant not nil,", nil)
		return
	}

	var bReq model.RequestCreateOrder
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusConflict, err.Error(), nil)
		return
	}

	bReq.UserID = userID

	insertOK, err := order.CreateOrder(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusCreated, helper.SUCCESS_MESSSAGE, insertOK)
}
