package database

import (
	"database/sql"
	"fmt"
	"os"

	//_ "github.com/glebarez/go-sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func GetDb() (*sql.DB, error) {
	initDb := false

	if _, err := os.Stat("./data.db"); err != nil {
		initDb = true
	}

	dsn := "./data.db"
	//db, err := sql.Open("sqlite", dsn)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("GetDb: could not open data.db: %w", err)
	}

	if initDb {
		if err = initializeDb(db); err != nil {
			return nil, fmt.Errorf("GetDb: could not initialize data.db: %w", err)
		}
	}

	return db, nil
}

func initializeDb(db *sql.DB) error {
	createStatements, err := getInitializeDbTableCreateStatements(db)
	if err != nil {
		return fmt.Errorf("InitializeDb: could not get table create statements: %w", err)
	}

	for _, statement := range createStatements {
		_, err := statement.Exec()
		if err != nil {
			return fmt.Errorf("InitializeDb: could not execute table create statement: %w", err)
		}
	}

	return nil
}

func getInitializeDbTableCreateStatements(db *sql.DB) ([]*sql.Stmt, error) {
	var statements []*sql.Stmt

	statement, err := db.Prepare(`
		CREATE TABLE profile (
			id INTEGER PRIMARY KEY,
			username TEXT NOT NULL,
			membership_type INTEGER NOT NULL,
			membership_id TEXT NOT NULL UNIQUE,
			json TEXT NOT NULL
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("getInitializeDbTableCreationStatements: could not create profile table query: %w", err)
	}

	statements = append(statements, statement)

	statement, err = db.Prepare(`
		CREATE TABLE activity (
			id INTEGER PRIMARY KEY,
			instance_id TEXT NOT NULL,
			membership_ids TEXT NOT NULL,
			membership_type INTEGER NOT NULL,
			character_ids TEXT NOT NULL
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("getInitializeDbTableCreationStatements: could not create activity table: %w", err)
	}

	statements = append(statements, statement)

	statement, err = db.Prepare(`
		CREATE TABLE activity_history (
			id INTEGER PRIMARY KEY,
			membership_id TEXT NOT NULL,
			membership_type INTEGER NOT NULL,
			character_id TEXT NOT NULL,
			activity_count INTEGER NOT NULL
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("getInitializeDbTableCreationStatements: could not create activity_history table: %w", err)
	}

	statements = append(statements, statement)

	statement, err = db.Prepare(`
			CREATE UNIQUE INDEX activity_history_unique_index
			ON activity_history (membership_id, membership_type, character_id
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("getInitializeDbTableCreationStatements: could not create activity_history_unique_index index: %w", err)
	}

	statements = append(statements, statement)

	statement, err = db.Prepare(`
		CREATE TABLE post_game_carnage_report (
			id INTEGER PRIMARY KEY,
			instance_id TEXT NOT NULL UNIQUE,
			json TEXT NOT NULL
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("getInitializeDbTableCreationStatements: could not create post_game_carnage_report table: %w", err)
	}

	statements = append(statements, statement)

	return statements, nil
}
