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

func (getPlayerResponse *GetPlayerResponse) UnmarshalHttpResponseBody(responseBody []byte) error {
	err := json.Unmarshal(responseBody, getPlayerResponse)
	if err != nil {
		return fmt.Errorf("GetPlayerResponse.UnmarshalHttpResponseBody: could not json decode response: %w", err)
	}

	return nil
}

func (api *Api) GetPlayer(user string) (*GetPlayerResponse, error) {
	log.Printf("Getting Player data from Bungie for: %v\n", user)

	entity := &GetPlayerResponse{}
	err := api.DoGetRequest(
		fmt.Sprintf("SearchDestinyPlayer/-1/%v/", url.QueryEscape(user)),
		entity,
	)
	if err != nil {
		return nil, fmt.Errorf("Api.GetPlayer: could not create bungie api request: %w", err)
	}

	return entity, nil
}
