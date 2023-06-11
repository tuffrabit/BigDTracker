package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

type DbActivityData struct {
	Id             int
	InstanceId     string
	MembershipIds  string
	MembershipType int
	CharacterIds   string
}

type DbActivity struct {
	Db
	PreparedStatements *DbActivityPreparedStatements
}

type DbActivityPreparedStatements struct {
	GetByInstanceId                            *sql.Stmt
	GetByMembershipIdMembershipTypeCharacterId *sql.Stmt
	Create                                     *sql.Stmt
	UpdateById                                 *sql.Stmt
}

func (dbActivity *DbActivity) Init(db *sql.DB) error {
	preparedStatements := &DbActivityPreparedStatements{}
	statement, err := db.Prepare("SELECT * FROM activity WHERE instance_id = ?")
	if err != nil {
		return fmt.Errorf("DbActivity.Init: could not create query: %w", err)
	}

	preparedStatements.GetByInstanceId = statement

	statement, err = db.Prepare("SELECT * FROM activity WHERE membership_ids LIKE '%' || ? || '%' AND membership_type= ? AND character_ids LIKE '%' || ? || '%'")
	if err != nil {
		return fmt.Errorf("DbActivity.Init: could not create query: %w", err)
	}

	preparedStatements.GetByMembershipIdMembershipTypeCharacterId = statement

	statement, err = db.Prepare("INSERT INTO activity(instance_id, membership_ids, membership_type, character_ids) values(?,?,?,?)")
	if err != nil {
		return fmt.Errorf("DbActivity.Init: could not create query: %w", err)
	}

	preparedStatements.Create = statement

	statement, err = db.Prepare("UPDATE activity SET instance_id = ?, membership_ids = ?, membership_type = ?, character_ids = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("DbActivity.Init: could not create query: %w", err)
	}

	preparedStatements.UpdateById = statement
	dbActivity.PreparedStatements = preparedStatements
	dbActivity.Db.Db = db

	return nil
}

func (dbActivity *DbActivity) GetActivitiesByInstanceId(instanceId string) ([]*DbActivityData, error) {
	log.Printf("Getting Activity data from DB for instance: %v\n", instanceId)

	rows, err := dbActivity.PreparedStatements.GetByInstanceId.Query(instanceId)
	if err != nil {
		return nil, fmt.Errorf("DbActivity.GetActivitiesByInstanceId: could not execute query: %w", err)
	}

	defer rows.Close()

	activities, err := dbActivity.getActivitiesFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("DbActivity.GetActivitiesByInstanceId: could not handle query result: %w", err)
	}

	return activities, nil
}

func (dbActivity *DbActivity) GetActivitiesByMembershipIdMembershipTypeCharacterId(membershipId string, membershipType int, characterId string) ([]*DbActivityData, error) {
	log.Printf("Getting Activity data from DB for membership_id: %v & membership_type: %v & character_id: %v\n", membershipId, membershipType, characterId)

	rows, err := dbActivity.PreparedStatements.GetByMembershipIdMembershipTypeCharacterId.Query(membershipId, membershipType, characterId)
	if err != nil {
		return nil, fmt.Errorf("DbActivity.GetActivitiesByMembershipIdMembershipTypeCharacterId: could not execute query: %w", err)
	}

	defer rows.Close()

	activities, err := dbActivity.getActivitiesFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("DbActivity.GetActivitiesByMembershipIdMembershipTypeCharacterId: could not handle query result: %w", err)
	}

	rowsErr := rows.Err()
	if rowsErr != nil {
		return nil, fmt.Errorf("DbActivity.GetActivitiesByMembershipIdMembershipTypeCharacterId: result set error: %w", err)
	}

	return activities, nil
}

func (dbActivity *DbActivity) CreateActivity(instanceId string, membershipIds string, membershipType int, characterIds string) error {
	log.Printf("Inserting Activity data to DB for instance: %v\n", instanceId)

	_, err := dbActivity.PreparedStatements.Create.Exec(instanceId, membershipIds, membershipType, characterIds)
	if err != nil {
		return fmt.Errorf("DbActivity.CreateActivity: could not execute insert: %w", err)
	}

	return nil
}

func (dbActivity *DbActivity) UpdateActivityById(id int, instanceId string, membershipIds string, membershipType int, characterIds string) error {
	log.Printf("Updating Activity data in DB for instance: %v\n", instanceId)

	_, err := dbActivity.PreparedStatements.UpdateById.Exec(instanceId, membershipIds, membershipType, characterIds, id)
	if err != nil {
		return fmt.Errorf("DbActivity.UpdateActivityById: could not execute update: %w", err)
	}

	return nil
}

func (dbActivity *DbActivity) getActivitiesFromDbRows(rows *sql.Rows) ([]*DbActivityData, error) {
	var activities []*DbActivityData
	var id int
	var instanceId string
	var membershipIds string
	var membershipType int
	var characterIds string

	for rows.Next() {
		err := rows.Scan(&id, &instanceId, &membershipIds, &membershipType, &characterIds)
		if err != nil {
			return nil, fmt.Errorf("DbActivity.getActivitiesFromDbRows: could not scan row: %w", err)
		}

		activity := &DbActivityData{}
		activity.Id = id
		activity.InstanceId = instanceId
		activity.MembershipIds = membershipIds
		activity.MembershipType = membershipType
		activity.CharacterIds = characterIds

		activities = append(activities, activity)
	}

	return activities, nil
}
