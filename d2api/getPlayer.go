package d2api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type GetPlayerResponse struct {
	Response []GetPlayerRepsonseBody `json:"Response"`
}

type GetPlayerRepsonseBody struct {
	MembershipId                string `json:"membershipId"`
	MembershipType              int    `json:"membershipType"`
	BungieGlobalDisplayNameCode int    `json:"bungieGlobalDisplayNameCode"`
}

func GetPlayer(apiKey string, user string) (*GetPlayerResponse, error) {
	log.Printf("Getting Player data from Bungie for: %v\n", user)

	request, err := http.NewRequest(
		"GET",
		baseUrl+"SearchDestinyPlayer/-1/"+url.QueryEscape(user)+"/",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("GetPlayer: could not create bungie api request: %w", err)
	}

	request.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GetPlayer: request to bungie api failed: %w", err)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetPlayer: response from bungie api returned http %v", httpResponse.StatusCode)
	}

	jsonResponse := &GetPlayerResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(jsonResponse)

	if err != nil {
		return nil, fmt.Errorf("GetPlayer: could not json decode response: %w", err)
	}

	return jsonResponse, nil
}
