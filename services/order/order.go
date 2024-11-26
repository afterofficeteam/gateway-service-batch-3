package order

import (
	"encoding/json"
	"errors"
	model "gateway-service/models"
	"gateway-service/util/helper"
	"net/http"
)

func CreateOrder(req model.RequestCreateOrder) (*string, error) {
	channel := make(chan helper.Response)
	clientRequest := helper.NetClientRequest{
		NetClient:  helper.DefaultClient,
		RequestUrl: "http://localhost:9992/order",
	}

	clientRequest.Post(req, channel)

	response := <-channel
	if response.StatusCode != http.StatusCreated {
		var responseError string
		if err := json.Unmarshal(response.Res, &responseError); err != nil {
			return nil, err
		}

		return nil, errors.New(string(responseError))
	}

	var userID string
	if err := json.Unmarshal(response.Res, &userID); err != nil {
		return nil, err
	}

	return &userID, nil
}
