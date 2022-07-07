package index

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuditLogRecord struct {
	// log_record fields
	Id                  int       `json:"id"`
	Action              int       `json:"action"`
	UserID              int       `json:"userId"`
	ContentItemID       int       `json:"contentItemId"`
	CreateDate          time.Time `json:"createDate"`
	CreateDateTimestamp int64     `json:"createDateTimestamp"`
	// Algolia related
	ObjectID string `json:"objectID"`
}

func (t *AuditLogRecord) CalculateFields() {
	id := uuid.New()
	timestamp := t.CreateDate.Unix()

	t.CreateDateTimestamp = timestamp
	t.ObjectID = fmt.Sprintf(
		"%v-%v",
		timestamp,
		id.String(),
	)
}
