package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	dbpackage "algolia/postgres/producer/db"
)

const (
	// it's best if they match with docker-entrypoint-initdb.d/004_generate_mock_data.sql
	USER_MIN = 10
	USER_MAX = 20
	CITM_MIN = 1500
	CITM_MAX = 2300
	ACTN_MIN = 1
	ACTN_MAX = 3
)

func main() {
	DELAY_SEC := os.Getenv("DELAY_SEC")
	delay, err := strconv.Atoi(DELAY_SEC)
	if err != nil {
		log.Panicf("Failed to read DELAY_SEC env variable")
	}
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	log.Print("Init Queue Producer")
	db := dbpackage.NewDB()
	err = db.OpenDB()
	if err != nil {
		log.Panicf("Failed to open DB", err)
	}
	for {
		log.Printf("Wait %d seconds...", delay)
		time.Sleep(time.Duration(delay) * time.Second)
		userID := randomizer.Intn(USER_MIN+USER_MAX) + USER_MIN
		citmID := randomizer.Intn(CITM_MIN+CITM_MAX) + CITM_MIN
		actnID := randomizer.Intn(ACTN_MIN+ACTN_MAX) + ACTN_MIN
		logID, err := db.AddAuditLogEntry(actnID, userID, citmID)
		if err != nil {
			log.Printf("ERROR: Failed to add log entry. err: %v", err)
		}
		log.Printf(
			"New audit entry (id: %v, userID: %v, contentItemID: %v, actionID: %v)",
			logID, userID, citmID, actnID,
		)
	}
}
