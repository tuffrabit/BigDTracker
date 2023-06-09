package d2api

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type GetProfileRepsonse struct {
	Response GetProfileRepsonseResponse `json:"Response"`
	Json     string
}

type GetProfileRepsonseResponse struct {
	Profile GetProfileRepsonseProfile `json:"profile"`
}

type GetProfileRepsonseProfile struct {
	Data GetProfileRepsonseData `json:"data"`
}

type GetProfileRepsonseData struct {
	CharacterIds []string `json:"characterIds"`
}

func (getProfileResponse *GetProfileRepsonse) UnmarshalHttpResponseBody(responseBody []byte) error {
	err := json.Unmarshal(responseBody, getProfileResponse)
	if err != nil {
		return fmt.Errorf("GetProfileRepsonse.UnmarshalHttpResponseBody: could not json decode response: %w", err)
	}

	getProfileResponse.Json = string(responseBody)

	return nil
}

func (api *Api) GetProfile(membershipId string, membershipType int) (*GetProfileRepsonse, error) {
	log.Printf("Getting Profile data from Bungie for: %v\n", membershipId)

	entity := &GetProfileRepsonse{}
	err := api.DoGetRequest(
		fmt.Sprintf("%v/Profile/%v/?components=100", strconv.Itoa(membershipType), membershipId),
		entity,
	)
	if err != nil {
		return nil, fmt.Errorf("GetProfile: could not create bungie api request: %w", err)
	}

	return entity, nil
}
