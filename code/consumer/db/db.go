package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	index "algolia/postgres/consumer/index"

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
	BeginTx() error
	Rollback() error
	Commit() error
	Query(statement string, args ...interface{}) (*sql.Rows, error)
	QueryRow(statement string, args ...interface{}) *sql.Row
	Exec(statement string, args ...interface{}) (sql.Result, error)
	GetLastNRecordsFromPostgresQueue(int) ([]index.AuditLogRecord, int, error)
	RemoveUploadedRows() error
}

func (d *dbStruct) OpenDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	d.db = db
	return err
}

func (d *dbStruct) BeginTx() error {
	if d.db == nil {
		return fmt.Errorf("Database connection is not opened yet")
	}
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	d.tx = tx
	return nil
}

func (d *dbStruct) Rollback() error {
	return d.tx.Rollback()
}

func (d *dbStruct) Commit() error {
	return d.tx.Commit()
}

func (d *dbStruct) Query(statement string, args ...interface{}) (*sql.Rows, error) {
	if d.tx != nil {
		return d.tx.Query(statement, args...)
	}
	if d.db != nil {
		return d.db.Query(statement, args...)
	}
	return nil, fmt.Errorf("Can not run query")
}
func (d *dbStruct) Exec(statement string, args ...interface{}) (sql.Result, error) {
	if d.tx != nil {
		return d.tx.Exec(statement, args...)
	}
	if d.db != nil {
		return d.db.Exec(statement, args...)
	}
	return nil, fmt.Errorf("Can not run Exec")
}
func (d *dbStruct) QueryRow(statement string, args ...interface{}) *sql.Row {
	if d.tx != nil {
		return d.tx.QueryRow(statement, args...)
	}
	if d.db != nil {
		return d.db.QueryRow(statement, args...)
	}
	return &sql.Row{}
}

func (d *dbStruct) GetLastNRecordsFromPostgresQueue(batchSize int) ([]index.AuditLogRecord, int, error) {
	result := []index.AuditLogRecord{}
	log.Println("GetLastNRecordsFromPostgresQueue", batchSize)

	var queue_size int
	err := d.db.QueryRow(sqlGetQueueLength).Scan(&queue_size)
	if err != nil {
		return result, 0, err
	}
	if queue_size == 0 {
		return result, queue_size, err
	}

	rows, err := d.Query(sqlGetItemsFromQueue, batchSize)
	for rows.Next() {
		line := index.AuditLogRecord{}
		err := rows.Scan(
			&line.Id,
			&line.Action,
			&line.UserID,
			&line.ContentItemID,
			&line.CreateDate,
		)
		if err != nil {
			log.Printf("ERROR: Failed to read line: %v", err)
		}
		line.CalculateFields()
		result = append(result, line)
	}
	return result, queue_size, err
}

func (d *dbStruct) RemoveUploadedRows() error {
	_, err := d.Exec(sqlDeleteVisitedQueueItems)
	return err
}

func NewDB() DB {
	return &dbStruct{}
}
