package database

import (
	"database/sql"
	"fmt"
	"log"

	//_ "github.com/glebarez/go-sqlite"
	_ "github.com/mattn/go-sqlite3"
)

type DbActivity struct {
	Id             int
	InstanceId     string
	MembershipIds  string
	MembershipType int
	CharacterIds   string
}

func GetActivitiesByInstanceId(db *sql.DB, instanceId string) ([]*DbActivity, error) {
	log.Printf("Getting Activity data from DB for instance: %v\n", instanceId)

	statement, err := db.Prepare("SELECT * FROM activity WHERE instance_id = ?")
	if err != nil {
		return nil, fmt.Errorf("GetActivitiesByInstanceId: could not create query: %w", err)
	}

	rows, err := statement.Query(instanceId)
	if err != nil {
		return nil, fmt.Errorf("GetActivitiesByInstanceId: could not execute query: %w", err)
	}

	defer rows.Close()

	activities, err := getActivitiesFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetActivitiesByInstanceId: could not handle query result: %w", err)
	}

	return activities, nil
}

func GetActivitiesByMembershipIdMembershipTypeCharacterId(db *sql.DB, membershipId string, membershipType int, characterId string) ([]*DbActivity, error) {
	log.Printf("Getting Activity data from DB for membership_id: %v and membership_type: %v and character_id: %v\n", membershipId, membershipType, characterId)

	statement, err := db.Prepare("SELECT * FROM activity WHERE membership_ids LIKE '%' || ? || '%' AND membership_type= ? AND character_ids LIKE '%' || ? || '%'")
	if err != nil {
		return nil, fmt.Errorf("GetActivitiesByMembershipIdMembershipTypeCharacterId: could not create query: %w", err)
	}

	rows, err := statement.Query(membershipId, membershipType, characterId)
	if err != nil {
		return nil, fmt.Errorf("GetActivitiesByMembershipIdMembershipTypeCharacterId: could not execute query: %w", err)
	}

	defer rows.Close()

	activities, err := getActivitiesFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetActivitiesByMembershipIdMembershipTypeCharacterId: could not handle query result: %w", err)
	}

	rowsErr := rows.Err()
	if rowsErr != nil {
		return nil, fmt.Errorf("GetActivitiesByMembershipIdMembershipTypeCharacterId: result set error: %w", err)
	}

	return activities, nil
}

func CreateActivity(db *sql.DB, instanceId string, membershipIds string, membershipType int, characterIds string) error {
	log.Printf("Inserting Activity data to DB for instance: %v\n", instanceId)
	statement, err := db.Prepare("INSERT INTO activity(instance_id, membership_ids, membership_type, character_ids) values(?,?,?,?)")
	if err != nil {
		return fmt.Errorf("CreateActivity: could not create query: %w", err)
	}

	_, err = statement.Exec(instanceId, membershipIds, membershipType, characterIds)
	if err != nil {
		return fmt.Errorf("CreateActivity: could not execute insert: %w", err)
	}

	return nil
}

func UpdateActivityById(db *sql.DB, id int, instanceId string, membershipIds string, membershipType int, characterIds string) error {
	log.Printf("Updating Activity data in DB for instance: %v\n", instanceId)
	statement, err := db.Prepare("UPDATE activity SET instance_id = ?, membership_ids = ?, membership_type = ?, character_ids = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("UpdateActivityById: could not create query: %w", err)
	}

	_, err = statement.Exec(instanceId, membershipIds, membershipType, characterIds, id)
	if err != nil {
		return fmt.Errorf("UpdateActivityById: could not execute update: %w", err)
	}

	return nil
}

func getActivitiesFromDbRows(rows *sql.Rows) ([]*DbActivity, error) {
	var activities []*DbActivity
	var id int
	var instanceId string
	var membershipIds string
	var membershipType int
	var characterIds string

	for rows.Next() {
		err := rows.Scan(&id, &instanceId, &membershipIds, &membershipType, &characterIds)
		if err != nil {
			return nil, fmt.Errorf("getActivitiesFromDbRows: could not scan row: %w", err)
		}

		activity := &DbActivity{}
		activity.Id = id
		activity.InstanceId = instanceId
		activity.MembershipIds = membershipIds
		activity.MembershipType = membershipType
		activity.CharacterIds = characterIds

		activities = append(activities, activity)
	}

	return activities, nil
}
