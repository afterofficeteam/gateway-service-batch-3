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

	if userID == "" {
		helper.HandleResponse(w, http.StatusBadRequest, "user id is required", nil)
		return
	}

	var bReq cart.CartInsertRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	bReq.UserId = userID

	if bReq.ProductId == "" {
		helper.HandleResponse(w, http.StatusBadRequest, "product id is required", nil)
		return
	}

	if bReq.Qty == 0 {
		helper.HandleResponse(w, http.StatusBadRequest, "quantity is required", nil)
		return
	}

	bRes, err := h.cart.InsertCart(ctx, &bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, bRes)
}

func (h *Handler) GetDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)
	productID := r.PathValue("id")

	bRes, err := h.cart.GetDetail(ctx, &cart.CartDetailRequest{
		Id:        userID,
		ProductId: productID,
	})
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, bRes)
}
