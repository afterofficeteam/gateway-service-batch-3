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

	if limiter := middleware.GetLimiter(userID); !limiter.Allow() {
		helper.HandleResponse(w, http.StatusTooManyRequests, "To many request, please try again later", nil)
		return
	}

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
	userID, productID := middleware.GetUserID(ctx), r.PathValue("id")

	if limiter := middleware.GetLimiter(userID); !limiter.Allow() {
		helper.HandleResponse(w, http.StatusTooManyRequests, "To many request, please try again later", nil)
		return
	}

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

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, productID := middleware.GetUserID(ctx), r.PathValue("product_id")

	if limiter := middleware.GetLimiter(userID); !limiter.Allow() {
		helper.HandleResponse(w, http.StatusTooManyRequests, "To many request, please try again later", nil)
		return
	}

	if userID == "" || productID == "" {
		helper.HandleResponse(w, http.StatusBadRequest, "user id and product id cant not nil", nil)
		return
	}

	deleteOK, err := h.cart.Delete(userID, productID)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, deleteOK)
}
