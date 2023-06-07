package d2api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const baseUrl string = "https://www.bungie.net/Platform/Destiny2/"
const ActivitiesPageSize int = 100
const BigDApiHash int = 3824106094

type Api struct {
	ApiKey                   string
	ActivitiesPageSizeString string
}

func (api *Api) Init(apiKey string) {
	api.ApiKey = apiKey
	api.ActivitiesPageSizeString = strconv.Itoa(ActivitiesPageSize)
}

func (api *Api) DoGetRequest(uri string) (*[]byte, error) {
	request, err := http.NewRequest(
		"GET",
		baseUrl+uri,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Api.DoGetRequest: could not create bungie api request: %w", err)
	}

	request.Header.Add("X-API-KEY", api.ApiKey)

	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Api.DoGetRequest: request to bungie api failed: %w", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Api.DoGetRequest: response from bungie api returned http %v", httpResponse.StatusCode)
	}

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("Api.DoGetRequest: could not read response body: %w", err)
	}

	return &body, nil
}
