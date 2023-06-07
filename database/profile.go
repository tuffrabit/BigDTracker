package database

import (
	"database/sql"
	"fmt"
	"log"

	//_ "github.com/glebarez/go-sqlite"
	_ "github.com/mattn/go-sqlite3"
)

type DbProfileData struct {
	Id             int
	Username       string
	MembershipType int
	MembershipId   string
	Json           string
}

type DbProfile struct {
	Db
	PreparedStatements *DbProfilePreparedStatements
}

type DbProfilePreparedStatements struct {
	GetByUsername     *sql.Stmt
	GetByMembershipId *sql.Stmt
	Create            *sql.Stmt
	UpdateById        *sql.Stmt
}

func (dbProfile *DbProfile) Init(db *sql.DB) error {
	preparedStatements := &DbProfilePreparedStatements{}
	statement, err := db.Prepare("SELECT * FROM profile WHERE username = ?")
	if err != nil {
		return fmt.Errorf("DbProfile.Init: could not create query: %w", err)
	}

	preparedStatements.GetByUsername = statement

	statement, err = db.Prepare("SELECT * FROM profile WHERE membership_id = ?")
	if err != nil {
		return fmt.Errorf("DbProfile.Init: could not create query: %w", err)
	}

	preparedStatements.GetByMembershipId = statement

	statement, err = db.Prepare("INSERT INTO profile(username, membership_type, membership_id, json) values(?,?,?,?)")
	if err != nil {
		return fmt.Errorf("DbProfile.Init: could not create query: %w", err)
	}

	preparedStatements.Create = statement

	statement, err = db.Prepare("UPDATE profile SET username = ?, membership_type = ?, membership_id = ?, json = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("DbProfile.Init: could not create query: %w", err)
	}

	preparedStatements.UpdateById = statement
	dbProfile.PreparedStatements = preparedStatements
	dbProfile.Db.Db = db

	return nil
}

func (dbProfile *DbProfile) GetProfilesByUsername(username string) ([]*DbProfileData, error) {
	log.Printf("Getting Profile data from DB for: %v\n", username)

	rows, err := dbProfile.PreparedStatements.GetByUsername.Query(username)
	if err != nil {
		return nil, fmt.Errorf("DbProfile.GetProfilesByUsername: could not execute query: %w", err)
	}

	defer rows.Close()

	profiles, err := dbProfile.getProfilesFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("DbProfile.GetProfilesByUsername: could not handle query result: %w", err)
	}

	return profiles, nil
}

func (dbProfile *DbProfile) GetProfileByMembershipId(membershipId string) ([]*DbProfileData, error) {
	log.Printf("Getting Profile data from DB for: %v\n", membershipId)

	rows, err := dbProfile.PreparedStatements.GetByMembershipId.Query(membershipId)
	if err != nil {
		return nil, fmt.Errorf("DbProfile.GetProfileByMembershipId: could not create profile table: %w", err)
	}

	defer rows.Close()

	profiles, err := dbProfile.getProfilesFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("DbProfile.GetProfileByMembershipId: could not handle query result: %w", err)
	}

	return profiles, nil
}

func (dbProfile *DbProfile) CreateProfile(username string, membershipType int, membershipId string, json string) error {
	log.Printf("Inserting Profiile data into DB for username: %v & membership_type: %v & membership_id: %v\n", username, membershipType, membershipId)

	_, err := dbProfile.PreparedStatements.Create.Exec(username, membershipType, membershipId, json)
	if err != nil {
		return fmt.Errorf("DbProfile.CreateProfile: could not execute insert: %w", err)
	}

	return nil
}

func (dbProfile *DbProfile) UpdateProfileById(id int, username string, membershipType int, membershipId string, json string) error {
	log.Printf("Updating Profiile data in DB for username: %v & membership_type: %v & membership_id: %v\n", username, membershipType, membershipId)

	_, err := dbProfile.PreparedStatements.UpdateById.Exec(username, membershipType, membershipId, json, id)
	if err != nil {
		return fmt.Errorf("DbProfile.UpdateProfileById: could not execute update: %w", err)
	}

	return nil
}

func (dbProfile *DbProfile) getProfilesFromDbRows(rows *sql.Rows) ([]*DbProfileData, error) {
	var profiles []*DbProfileData
	var profileId int
	var profileUsername string
	var profileMembershipType int
	var profileMembershipId string
	var profileJson string

	for rows.Next() {
		err := rows.Scan(&profileId, &profileUsername, &profileMembershipType, &profileMembershipId, &profileJson)
		if err != nil {
			return nil, fmt.Errorf("DbProfile.getProfilesFromDbRows: could not scan row: %w", err)
		}

		profile := &DbProfileData{}
		profile.Id = profileId
		profile.Username = profileUsername
		profile.MembershipType = profileMembershipType
		profile.MembershipId = profileMembershipId
		profile.Json = profileJson

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
