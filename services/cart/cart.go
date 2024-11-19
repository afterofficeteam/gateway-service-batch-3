package cart

import (
	"encoding/json"
	"errors"
	"gateway-service/util/helper"
	"net/http"
)

func Delete(userID, productID string) (*int, error) {
	channel := make(chan helper.Response)
	clientRequest := helper.NetClientRequest{
		NetClient:  helper.DefaultClient,
		RequestUrl: "http://localhost:9992/cart/" + userID + "/" + productID,
	}

	clientRequest.Delete(nil, channel)

	response := <-channel
	if response.StatusCode != http.StatusOK {
		var responseError string
		if err := json.Unmarshal(response.Res, &responseError); err != nil {
			return nil, err
		}

		return nil, errors.New(string(responseError))
	}

	var rowsAffected int
	if err := json.Unmarshal(response.Res, &rowsAffected); err != nil {
		return nil, err
	}

	return &rowsAffected, nil
}
