package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tuffrabit/BigDTracker/d2api/entity"
)

func DoGet(uri string, handler entity.Handler) (entity.Entity, error) {
	request, err := http.NewRequest(
		"GET",
		handler.GetMeta().BaseUrl+uri,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Api.DoGetRequest: could not create bungie api request: %w", err)
	}

	request.Header.Add("X-API-KEY", handler.GetMeta().ApiKey)

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

	entity := handler.NewEntity()
	err = json.Unmarshal(body, entity)
	if err != nil {
		return nil, fmt.Errorf("d2api/entity/player.UnmarshalHttpResponseBody: could not json decode response: %w", err)
	}

	entity.SetRawJson(string(body))

	return entity, nil
}
