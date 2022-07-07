package db

const (
	sqlInsertLogRecord = `
INSERT INTO app.audit_log (
	  action
	, user_id
	, content_item_id
)
VALUES ($1, $2, $3)
RETURNING id
`
)
