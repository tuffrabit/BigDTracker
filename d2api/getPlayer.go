package d2api

import (
	"encoding/json"
	"fmt"
	"log"
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

func (api *Api) GetPlayer(user string) (*GetPlayerResponse, error) {
	log.Printf("Getting Player data from Bungie for: %v\n", user)

	responseBody, err := api.DoGetRequest(fmt.Sprintf("SearchDestinyPlayer/-1/%v/", url.QueryEscape(user)))
	if err != nil {
		return nil, fmt.Errorf("Api.GetPlayer: could not create bungie api request: %w", err)
	}

	jsonResponse := &GetPlayerResponse{}
	err = json.Unmarshal(*responseBody, jsonResponse)

	if err != nil {
		return nil, fmt.Errorf("Api.GetPlayer: could not json decode response: %w", err)
	}

	return jsonResponse, nil
}
