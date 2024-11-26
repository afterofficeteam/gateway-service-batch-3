package order

import (
	model "gateway-service/models"
	"gateway-service/services/order"
)

func CreateOrder(req model.RequestCreateOrder) (*string, error) {
	return order.CreateOrder(req)
}
