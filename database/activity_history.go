package database

import (
	"database/sql"
	"fmt"
	"log"

	//_ "github.com/glebarez/go-sqlite"
	_ "github.com/mattn/go-sqlite3"
)

type DbActivityHistory struct {
	Id             int
	MembershipId   string
	MembershipType int
	CharacterId    string
	ActivityCount  int
}

func GetActivityHistoryByMembershipIdMembershipTypeCharacterId(db *sql.DB, membershipId string, membershipType int, characterId string) (*DbActivityHistory, error) {
	log.Printf("Getting ActivityHistory data from DB for membership_id: %v and membership_type: %v and character_id: %v\n", membershipId, membershipType, characterId)

	statement, err := db.Prepare("SELECT * FROM activity_history WHERE membership_id = ? AND membership_type = ? AND character_id = ?")
	if err != nil {
		return nil, fmt.Errorf("GetActivityHistoryByMembershipIdMembershipTypeCharacterId: could not create query: %w", err)
	}

	rows, err := statement.Query(membershipId, membershipType, characterId)
	if err != nil {
		return nil, fmt.Errorf("GetActivityHistoryByMembershipIdMembershipTypeCharacterId: could not execute query: %w", err)
	}

	defer rows.Close()

	activityHistory, err := getActivityHistoryFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetActivityHistoryByMembershipIdMembershipTypeCharacterId: could not handle query result: %w", err)
	}

	return activityHistory, nil
}

func CreateActivityHistory(db *sql.DB, membershipId string, membershipType int, characterId string, activityCount int) error {
	log.Printf("Inserting ActivityHistory data from DB for membership_id: %v, membership_type: %v, character_id: %v\n", membershipId, membershipType, characterId)

	statement, err := db.Prepare("INSERT INTO activity_history(membership_id, membership_type, character_id, activity_count) values(?,?,?,?)")
	if err != nil {
		return fmt.Errorf("CreateActivityHistory: could not create query: %w", err)
	}

	_, err = statement.Exec(membershipId, membershipType, characterId, activityCount)
	if err != nil {
		return fmt.Errorf("CreateActivityHistory: could not execute insert: %w", err)
	}

	return nil
}

func UpdateActivityHistoryById(db *sql.DB, id int, membershipId string, membershipType int, characterId string, activityCount int) error {
	log.Printf("Updating ActivityHistory data in DB for membership_id: %v, membership_type: %v, character_id: %v\n", membershipId, membershipType, characterId)

	statement, err := db.Prepare("UPDATE activity_history SET membership_id = ?, membership_type = ?, character_id = ?, activity_count = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("UpdateActivityHistoryById: could not create query: %w", err)
	}

	_, err = statement.Exec(membershipId, membershipType, characterId, activityCount, id)
	if err != nil {
		return fmt.Errorf("UpdateActivityHistoryById: could not execute update: %w", err)
	}

	return nil
}

func getActivityHistoryFromDbRows(rows *sql.Rows) (*DbActivityHistory, error) {
	activityHistory := &DbActivityHistory{}
	var id int
	var membershipId string
	var membershipType int
	var characterId string
	var activityCount int

	if rows.Next() {
		err := rows.Scan(&id, &membershipId, &membershipType, &characterId, &activityCount)
		if err != nil {
			return nil, fmt.Errorf("getActivityHistoryFromDbRows: could not scan row: %w", err)
		}

		activityHistory.Id = id
		activityHistory.MembershipId = membershipId
		activityHistory.MembershipType = membershipType
		activityHistory.CharacterId = characterId
		activityHistory.ActivityCount = activityCount
	}

	return activityHistory, nil
}
