package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

type DbPostGameCarnageReport struct {
	Id         int
	InstanceId string
	Json       string
}

func GetPostGameCarnageReportsByInstanceId(db *sql.DB, instanceId string) ([]*DbPostGameCarnageReport, error) {
	log.Printf("Getting PostGameCarnageReport data from DB for: %v\n", instanceId)

	statement, err := db.Prepare("SELECT * FROM post_game_carnage_report WHERE instance_id = ?")
	if err != nil {
		return nil, fmt.Errorf("GetPostGameCarnageReportByInstanceId: could not create query: %w", err)
	}

	rows, err := statement.Query(instanceId)
	if err != nil {
		return nil, fmt.Errorf("GetPostGameCarnageReportByInstanceId: could not create post_game_carnage_report query: %w", err)
	}

	defer rows.Close()

	postGameCarnageReports, err := getPostGameCarnageReportFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetPostGameCarnageReportByInstanceId: could not handle query result: %w", err)
	}

	return postGameCarnageReports, nil
}

func CreatePostGameCarnageReport(db *sql.DB, instanceId string, json string) error {
	log.Printf("Inserting PostGameCarnageReport data into DB for: %v\n", instanceId)

	statement, err := db.Prepare("INSERT INTO post_game_carnage_report(instance_id, json) values(?,?)")
	if err != nil {
		return fmt.Errorf("CreatePostGameCarnageReport: could not create query: %w", err)
	}

	_, err = statement.Exec(instanceId, json)
	if err != nil {
		return fmt.Errorf("CreatePostGameCarnageReport: could not execute insert: %w", err)
	}

	return nil
}

func UpdatePostGameCarnageReportById(db *sql.DB, id int, instanceId string, json string) error {
	statement, err := db.Prepare("UPDATE post_game_carnage_report SET instance_id = ?, json = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("UpdatePostGameCarnageReportById: could not create query: %w", err)
	}

	_, err = statement.Exec(instanceId, json, id)
	if err != nil {
		return fmt.Errorf("UpdatePostGameCarnageReportById: could not execute update: %w", err)
	}

	return nil
}

func getPostGameCarnageReportFromDbRows(rows *sql.Rows) ([]*DbPostGameCarnageReport, error) {
	var postGameCarnageReports []*DbPostGameCarnageReport
	var id int
	var instanceId string
	var json string

	for rows.Next() {
		err := rows.Scan(&id, &instanceId, &json)
		if err != nil {
			return nil, fmt.Errorf("getPostGameCarnageReportFromDbRows: could not scan row: %w", err)
		}

		postGameCarnageReport := &DbPostGameCarnageReport{}
		postGameCarnageReport.Id = id
		postGameCarnageReport.InstanceId = instanceId
		postGameCarnageReport.Json = json

		postGameCarnageReports = append(postGameCarnageReports, postGameCarnageReport)
	}

	return postGameCarnageReports, nil
}
