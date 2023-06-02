package d2api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type GetActivitiesRepsonse struct {
	Response GetActivitiesRepsonseResponse `json:"Response"`
}

type GetActivitiesRepsonseResponse struct {
	Activities []GetActivitiesRepsonseActivity `json:"activities"`
}

type GetActivitiesRepsonseActivity struct {
	ActivityDetails GetActivitiesRepsonseActivityDetails `json:"activityDetails"`
}

type GetActivitiesRepsonseActivityDetails struct {
	ReferenceId int    `json:"referenceId"`
	InstanceId  string `json:"instanceId"`
}

func GetActivities(apiKey string, page int, membershipId string, membershipType int, characterId string) (*GetActivitiesRepsonse, error) {
	log.Printf("Getting Activity data from Bungie for: MID:%v CID:%v P:%v\n", membershipId, characterId, page)

	request, err := http.NewRequest(
		"GET",
		baseUrl+strconv.Itoa(membershipType)+"/Account/"+membershipId+"/Character/"+characterId+"/Stats/Activities/?mode=AllPvP&count=100&page="+strconv.Itoa(page),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("GetActivities: could not create bungie api request: %w", err)
	}

	request.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GetActivities: request to bungie api failed: %w", err)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetActivities: response from bungie api returned http %v", httpResponse.StatusCode)
	}

	jsonResponse := &GetActivitiesRepsonse{}
	err = json.NewDecoder(httpResponse.Body).Decode(jsonResponse)

	if err != nil {
		return nil, fmt.Errorf("GetActivities: could not json decode response: %w", err)
	}

	return jsonResponse, nil
}
