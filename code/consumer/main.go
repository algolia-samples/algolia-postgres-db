package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	dbpackage "algolia/postgres/consumer/db"
	index "algolia/postgres/consumer/index"
)

func main() {
	DELAY_SEC := os.Getenv("DELAY_SEC")
	delay_sec, err := strconv.Atoi(DELAY_SEC)
	if err != nil {
		log.Panicf("Failed to read DELAY_SEC env variable")
	}
	BATCH_SIZE := os.Getenv("BATCH_SIZE")
	batch_size, err := strconv.Atoi(BATCH_SIZE)
	if err != nil {
		log.Panicf("Failed to read BATCH_SIZE env variable")
	}

	log.Println("Init QueueConsumer")
	db := dbpackage.NewDB()
	err = db.OpenDB()
	if err != nil {
		log.Panicf("Failed to open DB", err)
	}

	for {
		delay := time.Duration(delay_sec + rand.Intn(5))
		log.Printf("Wait %d seconds...", delay)
		rand.Seed(time.Now().UnixNano())
		time.Sleep(delay * time.Second)
		err := db.BeginTx()
		if err != nil {
			log.Printf("ERROR: ", fmt.Errorf("Failed to begin transaction: %v", err))
			continue
		}
		err = processQueue(db, batch_size)
		if err != nil {
			log.Printf("Skip processing the queue: %v", err)
			db.Rollback()
			continue
		}
		db.Commit()
	}
}

func processQueue(db dbpackage.DB, batch_size int) error {
	log_entries, queue_length, err := db.GetLastNRecordsFromPostgresQueue(batch_size)
	if err != nil {
		return fmt.Errorf("Failed to get queue items, try again in the next iteration", err)
	}
	if len(log_entries) == 0 {
		return fmt.Errorf("Queue is empty")
	}
	log.Printf("%d items are in the queue, got last %d items according to batch size: %d", queue_length, len(log_entries), batch_size)
	err = index.UploadRecordsToAlgola(log_entries)
	if err != nil {
		return fmt.Errorf("Failed to upload records to Algolia", err)
	}
	err = db.RemoveUploadedRows()
	if err != nil {
		return fmt.Errorf("Failed to remove uploaded records from the queue", err)
	}
	return nil
}
