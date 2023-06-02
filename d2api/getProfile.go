package d2api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func GetProfile(apiKey string, membershipId string, membershipType int) (*GetProfileRepsonse, error) {
	log.Printf("Getting Profile data from Bungie for: %v\n", membershipId)

	request, err := http.NewRequest(
		"GET",
		baseUrl+strconv.Itoa(membershipType)+"/Profile/"+membershipId+"/?components=100",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("GetProfile: could not create bungie api request: %w", err)
	}

	request.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GetProfile: request to bungie api failed: %w", err)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetProfile: response from bungie api returned http %v", httpResponse.StatusCode)
	}

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("GetProfile: could not read response body: %w", err)
	}

	jsonResponse := &GetProfileRepsonse{}
	err = json.Unmarshal(body, jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("GetProfile: could not json decode response: %w", err)
	}

	jsonResponse.Json = string(body)

	return jsonResponse, nil
}
