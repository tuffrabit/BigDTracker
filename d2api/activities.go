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

func (getActivitiesRepsonse *GetActivitiesRepsonse) UnmarshalHttpResponseBody(responseBody []byte) error {
	err := json.Unmarshal(responseBody, getActivitiesRepsonse)
	if err != nil {
		return fmt.Errorf("GetActivitiesRepsonse.UnmarshalHttpResponseBody: could not json decode response: %w", err)
	}

	return nil
}

func (api *Api) GetActivities(page int, membershipId string, membershipType int, characterId string) (*GetActivitiesRepsonse, error) {
	log.Printf("Getting Activity data from Bungie for: MID:%v CID:%v P:%v\n", membershipId, characterId, page)

	entity := &GetActivitiesRepsonse{}
	err := api.DoGetRequest(
		fmt.Sprintf("%v/Account/%v/Character/%v/Stats/Activities/?mode=AllPvP&count=%v&page=%v", strconv.Itoa(membershipType), membershipId, characterId, api.ActivitiesPageSizeString, strconv.Itoa(page)),
		entity,
	)
	if err != nil {
		return nil, fmt.Errorf("Api.GetActivities: could not create bungie api request: %w", err)
	}

	return entity, nil
}
