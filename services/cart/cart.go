package cart

import (
	"encoding/json"
	"errors"
	"gateway-service/util/helper"
	"gateway-service/util/middleware"
	"net/http"
)

func Delete(userID, productID string) (*int, error) {
	cb := middleware.CircuitBreaker

	url_valid := "http://localhost:9992/cart/" + userID + "/" + productID
	url_non_valid := "http://localhost:1001/cart/" + userID + "/" + productID

	var response int
	for i := 0; i < 20; i++ {
		url := url_non_valid
		if i > 15 {
			url = url_valid
		}

		res, err := middleware.CircuitBreakerExecute(cb, func() (interface{}, error) {
			channel := make(chan helper.Response)
			client_request := helper.NetClientRequest{
				NetClient:  helper.DefaultClient,
				RequestUrl: url,
			}

			client_request.Delete(nil, channel)

			response := <-channel
			if response.StatusCode != http.StatusOK {
				var responseError string
				if err := json.Unmarshal(response.Res, &responseError); err != nil {
					return nil, err
				}
				return nil, errors.New(responseError)
			}

			var rows_affected int
			if err := json.Unmarshal(response.Res, &rows_affected); err != nil {
				return nil, err
			}

			return rows_affected, nil
		})

		if err != nil {
			return nil, err
		}

		response = res.(int)
	}

	return &response, nil
}
