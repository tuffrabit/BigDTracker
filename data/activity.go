package data

import (
	"fmt"
	"log"
	"strings"

	"github.com/tuffrabit/BigDTracker/d2api"
	"github.com/tuffrabit/BigDTracker/database"
)

type Activity struct {
	InstanceId string
	DbActivity *database.DbActivityData
}

func GetActivities(dbHandlers *database.DbHandlers, api *d2api.Api, membershipId string, membershipType int, characterId string) ([]*Activity, error) {
	err := populateActivities(dbHandlers, api, membershipId, membershipType, characterId)
	if err != nil {
		return nil, fmt.Errorf("GetActivities: could not populate activities: %w", err)
	}

	dbActivities, err := dbHandlers.Activity.GetActivitiesByMembershipIdMembershipTypeCharacterId(
		membershipId,
		membershipType,
		characterId,
	)
	if err != nil {
		return nil, fmt.Errorf("GetActivities: could not get activities from db: %w", err)
	}

	log.Printf("Activities retrieved count: %v", len(dbActivities))

	var activities []*Activity

	for _, dbActivity := range dbActivities {
		activity := &Activity{}

		activity.InstanceId = dbActivity.InstanceId
		activity.DbActivity = dbActivity

		activities = append(activities, activity)
	}

	return activities, nil
}

func populateActivities(dbHandlers *database.DbHandlers, api *d2api.Api, membershipId string, membershipType int, characterId string) error {
	activityHistory, err := dbHandlers.ActivityHistory.GetActivityHistoryByMembershipIdMembershipTypeCharacterId(membershipId, membershipType, characterId)
	if err != nil {
		return fmt.Errorf("populateActivities: could not get activity history from db: %w", err)
	}

	log.Printf("ActivityHistory Count: %v for membership_id: %v, membership_type: %v, character_id: %v\n", activityHistory.ActivityCount, membershipId, membershipType, characterId)

	startPage := 0
	activitiesCount := 0

	if activityHistory.Id > 0 {
		startPage = getActivityD2ApiStartPage(activityHistory.ActivityCount)
	}

	log.Printf("Activity start page #: %v for membership_id: %v, membership_type: %v, character_id: %v\n", startPage, membershipId, membershipType, characterId)

	for {
		activityResponse, err := api.GetActivities(startPage, membershipId, membershipType, characterId)
		if err != nil {
			return fmt.Errorf("populateActivities: could not get activity history from db: %w", err)
		}

		activityResponseCount := len(activityResponse.Response.Activities)

		if activityResponseCount < 1 {
			break
		}

		activitiesCount = activitiesCount + activityResponseCount

		err = updateDbActivities(dbHandlers.Activity, activityResponse.Response.Activities, membershipId, membershipType, characterId)
		if err != nil {
			return fmt.Errorf("populateActivities: could not get activity history from db: %w", err)
		}

		if activityResponseCount < d2api.ActivitiesPageSize {
			break
		}

		startPage = startPage + 1
	}

	if activityHistory.Id == 0 {
		err := dbHandlers.ActivityHistory.CreateActivityHistory(membershipId, membershipType, characterId, activitiesCount)
		if err != nil {
			return fmt.Errorf("populateActivities: could not create activity history in db: %w", err)
		}
	} else {
		err := dbHandlers.ActivityHistory.UpdateActivityHistoryById(activityHistory.Id, membershipId, membershipType, characterId, activitiesCount)
		if err != nil {
			return fmt.Errorf("populateActivities: could not update activity history in db: %w", err)
		}
	}

	return nil
}

func updateDbActivities(db *database.DbActivity, apiActivities []d2api.GetActivitiesRepsonseActivity, membershipId string, membershipType int, characterId string) error {
	for _, activity := range apiActivities {
		dbActivities, err := db.GetActivitiesByInstanceId(activity.ActivityDetails.InstanceId)
		if err != nil {
			return fmt.Errorf("updateDbActivities: could not get activities from db: %w", err)
		}

		dbActivityUpdated := false

		for _, dbActivity := range dbActivities {
			if dbActivity.MembershipType == membershipType {
				doUpdate := false
				membershipIds := strings.Split(dbActivity.MembershipIds, ",")
				characterIds := strings.Split(dbActivity.CharacterIds, ",")

				if !contains(characterIds, characterId) {
					characterIds = append(characterIds, characterId)
					doUpdate = true
				}

				if !contains(membershipIds, membershipId) {
					membershipIds = append(membershipIds, membershipId)
					doUpdate = true
				}

				if doUpdate {
					err = db.UpdateActivityById(
						dbActivity.Id,
						dbActivity.InstanceId,
						strings.Join(membershipIds, ","),
						membershipType,
						strings.Join(characterIds, ","),
					)
					if err != nil {
						return fmt.Errorf("updateDbActivities: could not update activity in db: %w", err)
					}
				}

				dbActivityUpdated = true
			}
		}

		if !dbActivityUpdated {
			db.CreateActivity(activity.ActivityDetails.InstanceId, membershipId, membershipType, characterId)
			if err != nil {
				return fmt.Errorf("updateDbActivities: could not create activity in db: %w", err)
			}
		}
	}

	return nil
}

func getActivityD2ApiStartPage(currentActivityHistoryCount int) int {
	startPage := 0

	if currentActivityHistoryCount < d2api.ActivitiesPageSize {
		return startPage
	}

	startPage = currentActivityHistoryCount / d2api.ActivitiesPageSize

	return startPage
}

func contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if a == needle {
			return true
		}
	}

	return false
}
