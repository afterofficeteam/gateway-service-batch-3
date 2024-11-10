package cart

import (
	"encoding/json"
	"gateway-service/proto/cart"
	"gateway-service/util/helper"
	"gateway-service/util/middleware"
	"net/http"

	cartSvc "gateway-service/usecases/cart"
)

type Handler struct {
	cart cartSvc.CartSvc
}

func NewHandler(cart cartSvc.CartSvc) *Handler {
	return &Handler{
		cart: cart,
	}
}

func (h *Handler) InsertCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	var bReq cart.CartInsertRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bReq.UserId = userID

	bRes, err := h.cart.InsertCart(ctx, &bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, bRes)
}
