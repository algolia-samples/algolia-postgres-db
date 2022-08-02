package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = 5432
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

type dbStruct struct {
	db *sql.DB
	tx *sql.Tx
}

type DB interface {
	OpenDB() error
	QueryRow(statement string, args ...interface{}) *sql.Row
	AddAuditLogEntry(actionID, userID, contentItemID int) (int, error)
}

func (d *dbStruct) OpenDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	d.db = db
	return err
}

func (d *dbStruct) QueryRow(statement string, args ...interface{}) *sql.Row {
	if d.db != nil {
		return d.db.QueryRow(statement, args...)
	}
	return &sql.Row{}
}

func (d *dbStruct) AddAuditLogEntry(actionID, userID, contentItemID int) (int, error) {
	// log.Printf("AddAuditLogEntry (userID: %v, contentItemID: %v, actionID: %v)",
	// 	actionID, userID, contentItemID,
	// )

	var id int
	err := d.db.QueryRow(sqlInsertLogRecord, actionID, userID, contentItemID).Scan(&id)
	return id, err
}

func NewDB() DB {
	return &dbStruct{}
}
