package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DbActivityHistoryData struct {
	Id             int
	MembershipId   string
	MembershipType int
	CharacterId    string
	ActivityCount  int
}

type DbActivityHistory struct {
	Db
	PreparedStatements *DbActivityHistoryPreparedStatements
}

type DbActivityHistoryPreparedStatements struct {
	GetByMembershipIdMembershipTypeCharacterId *sql.Stmt
	Create                                     *sql.Stmt
	UpdateById                                 *sql.Stmt
}

func (dbActivityHistory *DbActivityHistory) Init(db *sql.DB) error {
	preparedStatements := &DbActivityHistoryPreparedStatements{}
	statement, err := db.Prepare("SELECT * FROM activity_history WHERE membership_id = ? AND membership_type = ? AND character_id = ?")
	if err != nil {
		return fmt.Errorf("DbActivityHistory.Init: could not create query: %w", err)
	}

	preparedStatements.GetByMembershipIdMembershipTypeCharacterId = statement

	statement, err = db.Prepare("INSERT INTO activity_history(membership_id, membership_type, character_id, activity_count) values(?,?,?,?)")
	if err != nil {
		return fmt.Errorf("DbActivityHistory.Init: could not create query: %w", err)
	}

	preparedStatements.Create = statement

	statement, err = db.Prepare("UPDATE activity_history SET membership_id = ?, membership_type = ?, character_id = ?, activity_count = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("DbActivityHistory.Init: could not create query: %w", err)
	}

	preparedStatements.UpdateById = statement
	dbActivityHistory.PreparedStatements = preparedStatements
	dbActivityHistory.Db.Db = db

	return nil
}

func (dbActivityHistory *DbActivityHistory) GetActivityHistoryByMembershipIdMembershipTypeCharacterId(membershipId string, membershipType int, characterId string) (*DbActivityHistoryData, error) {
	log.Printf("Getting ActivityHistory data from DB for membership_id: %v and membership_type: %v and character_id: %v\n", membershipId, membershipType, characterId)

	rows, err := dbActivityHistory.PreparedStatements.GetByMembershipIdMembershipTypeCharacterId.Query(membershipId, membershipType, characterId)
	if err != nil {
		return nil, fmt.Errorf("DbActivityHistory.GetActivityHistoryByMembershipIdMembershipTypeCharacterId: could not execute query: %w", err)
	}

	defer rows.Close()

	activityHistory, err := dbActivityHistory.getActivityHistoryFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("DbActivityHistory.GetActivityHistoryByMembershipIdMembershipTypeCharacterId: could not handle query result: %w", err)
	}

	return activityHistory, nil
}

func (dbActivityHistory *DbActivityHistory) CreateActivityHistory(membershipId string, membershipType int, characterId string, activityCount int) error {
	log.Printf("Inserting ActivityHistory data from DB for membership_id: %v, membership_type: %v, character_id: %v\n", membershipId, membershipType, characterId)

	_, err := dbActivityHistory.PreparedStatements.Create.Exec(membershipId, membershipType, characterId, activityCount)
	if err != nil {
		return fmt.Errorf("DbActivityHistory.CreateActivityHistory: could not execute insert: %w", err)
	}

	return nil
}

func (dbActivityHistory *DbActivityHistory) UpdateActivityHistoryById(id int, membershipId string, membershipType int, characterId string, activityCount int) error {
	log.Printf("Updating ActivityHistory data in DB for membership_id: %v, membership_type: %v, character_id: %v\n", membershipId, membershipType, characterId)

	_, err := dbActivityHistory.PreparedStatements.UpdateById.Exec(membershipId, membershipType, characterId, activityCount, id)
	if err != nil {
		return fmt.Errorf("DbActivityHistory.UpdateActivityHistoryById: could not execute update: %w", err)
	}

	return nil
}

func (dbActivityHistory *DbActivityHistory) getActivityHistoryFromDbRows(rows *sql.Rows) (*DbActivityHistoryData, error) {
	activityHistory := &DbActivityHistoryData{}
	var id int
	var membershipId string
	var membershipType int
	var characterId string
	var activityCount int

	if rows.Next() {
		err := rows.Scan(&id, &membershipId, &membershipType, &characterId, &activityCount)
		if err != nil {
			return nil, fmt.Errorf("DbActivityHistory.getActivityHistoryFromDbRows: could not scan row: %w", err)
		}

		activityHistory.Id = id
		activityHistory.MembershipId = membershipId
		activityHistory.MembershipType = membershipType
		activityHistory.CharacterId = characterId
		activityHistory.ActivityCount = activityCount
	}

	return activityHistory, nil
}
