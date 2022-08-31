package db

const (
	sqlGetQueueLength = `
select 
	count(*) as queue_length
from app.queue_audit_log
`
	sqlGetItemsFromQueue = `
with get_lines as (
	select
		  id
		, action
		, user_id
		, content_item_id
		, create_date
		, _visited
	from app.queue_audit_log
	where _visited = false
	order by create_date desc
	limit $1
	for update skip locked -- add concurent consumers
)
update 
	app.queue_audit_log new
set _visited = true
from get_lines old
where old.id = new.id
returning 
	  new.id
	, new.action
	, new.user_id
	, new.content_item_id
	, new.create_date
	-- , new._visited
;
`
	sqlDeleteVisitedQueueItems = `
delete 
from app.queue_audit_log
where _visited = true
`
)
