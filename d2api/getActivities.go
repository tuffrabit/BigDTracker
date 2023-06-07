package d2api

import (
	"encoding/json"
	"fmt"
	"log"
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

func (api *Api) GetActivities(page int, membershipId string, membershipType int, characterId string) (*GetActivitiesRepsonse, error) {
	log.Printf("Getting Activity data from Bungie for: MID:%v CID:%v P:%v\n", membershipId, characterId, page)

	responseBody, err := api.DoGetRequest(fmt.Sprintf("%v/Account/%v/Character/%v/Stats/Activities/?mode=AllPvP&count=%v&page=%v", strconv.Itoa(membershipType), membershipId, characterId, api.ActivitiesPageSizeString, strconv.Itoa(page)))
	if err != nil {
		return nil, fmt.Errorf("Api.GetActivities: could not create bungie api request: %w", err)
	}

	jsonResponse := &GetActivitiesRepsonse{}
	err = json.Unmarshal(*responseBody, jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("Api.GetActivities: could not json decode response: %w", err)
	}

	return jsonResponse, nil
}
