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

func (api *Api) GetProfile(membershipId string, membershipType int) (*GetProfileRepsonse, error) {
	log.Printf("Getting Profile data from Bungie for: %v\n", membershipId)

	responseBody, err := api.DoGetRequest(fmt.Sprintf("%v/Profile/%v/?components=100", strconv.Itoa(membershipType), membershipId))
	if err != nil {
		return nil, fmt.Errorf("GetProfile: could not create bungie api request: %w", err)
	}

	jsonResponse := &GetProfileRepsonse{}
	err = json.Unmarshal(*responseBody, jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("Api.GetProfile: could not json decode response: %w", err)
	}

	jsonResponse.Json = string(*responseBody)

	return jsonResponse, nil
}
