package cart

import (
	"encoding/json"
	"errors"
	"gateway-service/util/helper"
	"net/http"
)

func Delete(userID, productID string) (*int, error) {
	println(userID)
	channel := make(chan helper.Response)
	clientRequest := helper.NetClientRequest{
		NetClient:  helper.DefaultClient,
		RequestUrl: "http://localhost:9992/cart/" + userID + "/" + productID,
	}

	clientRequest.Delete(nil, channel)

	response := <-channel
	if response.StatusCode != http.StatusOK || response.Err != nil {
		var responseError string
		if err := json.Unmarshal(response.Res, &responseError); err != nil {
			return nil, errors.New(responseError)
		}
	}

	var rowsAffected int
	if err := json.Unmarshal(response.Res, &rowsAffected); err != nil {
		return nil, err
	}

	return &rowsAffected, nil
}
