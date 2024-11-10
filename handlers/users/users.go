package users

import (
	"encoding/json"
	model "gateway-service/models"
	"gateway-service/usecases/users"
	"gateway-service/util/helper"
	"net/http"

	"github.com/go-playground/validator"
)

type Handler struct {
	userSvc   users.UserSvc
	validator *validator.Validate
}

func NewHandler(userSvc users.UserSvc, validator *validator.Validate) *Handler {
	return &Handler{
		userSvc:   userSvc,
		validator: validator,
	}
}

func (h *Handler) SignUpByEmail(w http.ResponseWriter, r *http.Request) {
	var bReq model.Users
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if bReq.CategoryPreferences == nil {
		bReq.CategoryPreferences = []string{}
	}

	if err := h.validator.Struct(bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, err := h.userSvc.UserRegister(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusCreated, helper.SUCCESS_MESSSAGE, userID)
}

func (h *Handler) SignInByEmail(w http.ResponseWriter, r *http.Request) {
	var bReq model.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.validator.Struct(bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	bRes, err := h.userSvc.UserLogin(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, bRes)
}
