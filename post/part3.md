# Index Postgres Data with Algolia - Part 3

> This post is a follow-up on [how we index](./part2.md) our data with Algolia.

In order to let our users access our audit data from an external location we need to send it to a place where they can access it. We're going to walk through the implementation details on how we achieve that.

## Implementation Details

Let's dive into our solution in detail. Our demo follows a simple pattern.

We create new audit log lines periodically with a single producer,
 and read these new lines with multiple data consumers.

For the sake of simplicity we've created a minimal example for the services.

- [Producer](#producer) puts random activity in the database
- [Consumer](#consumer) reads the queue and uploads the data into Algolia

For easier readability, and sticking to the point, I used only the important external libraries, and kept the code as simple as possible.

This example does not reflect our real codebase, but gives you an idea of how it can be implemented with minimal overhead.

## Getting started

Let's go through again what you need to get started as we [did in the earlier post](./part2.md).

All components and depenencies are described in the `docker-compose.yml` file.

- To run the example as intended you have to have [Docker installed](https://docs.docker.com/get-docker/).
- You need to create a `.env` file based on `.env.example` file
- You need to register in [Algolia](https://algolia.com) and set up your credentials in the `.env` file.
- Note that the docker-compose file uses environment variables, it will only work as intended if it's started from the same folder as where the `.env` file is located

The example can be started with:

```bash
docker-compose up --build
```

After start the following should happen.

1. Postgres container starts
1. Postgres container initializes the database with mock data
1. After Postgres container is in a healthy stage, a single producer, and 2 consumers start up
1. The applications wait and do their job periodically as set in `DELAY_PRODUCER` and `DELAY_CONSUMER` environment variables.

## Component Details

### Postgres Database

We run the [official docker image](https://hub.docker.com/_/postgres) of [postgres](./postgres).

Upon first start it:

1. Creates a db folder (or optionally a volume) in order to persist the database between runs
1. Runs the init sqls that are inside [docker-entrypoint-initdb.d](./postgres/docker-entrypoint-initdb.d).
  You can find more info on the init scripts at [dockerhub's postgres image](https://hub.docker.com/_/postgres).

#### Init Scripts

The init scripts run in alphabetical order.
The first set of sql files (`001_init_audit_log_table.sql` and `002_queue_table.sql`) create the tables that are used by the applications.

1. Create the `app` schema.
1. Create `app.audit_log` fact table
1. Create `app.queue_audit_log` queue table

After the tables are ready `003_queue_trigger.sql` creates a trigger to catch the new data inserted into the `app.audit_log` table and replicate it to the `app.queue_audit_log` queue.

```sql
create or replace function audit_insert_trigger_fnc()
  returns trigger as $$
    begin
        insert into 
            app.queue_audit_log ( 
             action
            ,user_id
            ,content_item_id
            ,create_date
            )
        values(
             new."action"
            ,new."user_id"
            ,new."content_item_id"
            ,new."create_date"
        );

        return new;
    end;
$$ language 'plpgsql';


create trigger audit_insert_trigger
  after insert on app.audit_log
  for each row
  execute procedure audit_insert_trigger_fnc();
```

When the stucture is ready, and the trigger is in place, the last script (`004_generate_mock_data.sql`) generates random data into the fact table.
Its configurable parts are extracted into variables, so we can see how it behaves for different amount of data.
The randomizer has a hard-coded init seed, so it should generate the same data across multiple recreations.

```sql
-- set random seed for repeatable random data generation
SELECT setseed(0.8);
DO $$
    DECLARE
        -- configurable parameters for data generation
        nr_lines integer := 20;
        user_min integer := 10;
        user_max integer := 20;
        citm_min integer := 1500;
        citm_max integer := 2300;
        actn_min integer := 1;
        actn_max integer := 3;
    BEGIN
        with
            -- generate user_ids
            users as (
                select generate_series(user_min, user_max) as user_id
            )
            -- generate content_ids
           ,content as (
               select generate_series(citm_min, citm_max) as content_id
            )
            -- generate action_ids
           ,actions as (
               select generate_series(actn_min, actn_max) as action_id
            )
            -- get the cartesian product of the above in a random sort
           ,limited_data as (
               select
                 random() randomizer
                 ,* 
               from users, content, actions 
               order by randomizer
               limit nr_lines
            )
        insert 
            into app.audit_log (
                action
                ,user_id
                ,content_item_id
            )
            select
                 action_id
                ,user_id
                ,content_id
            from limited_data
        ;
END $$
;

-- select * from audit_log order by content_item_id, user_id, action;
```

This mock data generation script uses a controlled random data generation by setting the initial random seed at the start of the code with `setseed`.
We generate a `random()` number for each generated lines, and we can use this to avoid adding similar lines.

We generate identifiers with `generate_series` between the configurable ranges for each values.

In order to select only the given number of items we add an upper bound of the resultset with `limit`.

To make the code better separated the different logical components are defined in their own [Common Table Expressions](https://www.postgresql.org/docs/current/queries-with.html) aka. CTE-s defined by `with` queries.

In the `limited_data` CTE we join together all generated lines for the different data types and shuffle them before limiting the results.

The official PostgreSQL docker image is written in a way that the database init scripts are only started if it's the first start of the database.
If you stop the services, and then restart it again, the initialization does not happen again, but the data shall be still there.

### Producer

The code inside `./producer` folder represents our application.
In our scenario we don't want to modify this code, but leverage the power of Algolia through postgreSQL.

- Connects to the database on start.
- Periodically generate a single new random log line into the fact table.
- Out of this application's scope the insert trigger will copy this data into the queue.

This is a straightforward application.
The `main` function is where the most of the action happens.
The `db` folder contains a PostgreSQL connector in `db.go`, and an insert statement in `sql.go`.

### Consumer

The consumer services under `./consumer` folder read the last inserted lines from the database and put them into our Algolia index.

Connects to the database on start and then periodically reads the queue for new data and does the following in a transaction:

1. Reads the last `N` lines, then marks them as visited
1. Uploads the selected lines into Algolia
1. Clears the visited lines from the queue

We assume that multiple consumers shall be available to adjust to heavy loads if neecessary.
We can not rely on the fact that these consumers are running at all times, because of this constraint we can not leverage the [notify](https://www.postgresql.org/docs/current/sql-notify.html)/[listen](https://www.postgresql.org/docs/current/sql-listen.html) pattern of postgreSQL. We're using a queue table instead.

The heart of this concept lies in the following SQL query.

```sql
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
```

Let's separate what this query does into separate steps.
All theese steps are happening all at once in a single instruction inside a transaction started by our go code.

In order to let our queue table accessed by multiple consumers we need to add `for update skip locked`.
The `for update` clause lets the engine now that the select subquery will be used for updating them.
The `skip locked` part will ignore lines that other coonsumers might have locked.

> With SKIP LOCKED, any selected rows that cannot be immediately locked are skipped. Skipping locked rows provides an inconsistent view of the data, so this is not suitable for general purpose work, but can be used to avoid lock contention with multiple consumers accessing a queue-like table. [Sourece](https://www.postgresql.org/docs/current/sql-select.html).

- The `order by create_date desc` makes sure that we get the latest lines from the queue that are available.
- The `limit $1` line makes sure that we only select a subset of lines from the queue.
- The `returning` declaration lets us get the data of each seleted line.
- In the update statement we set the `set _visited = true` only from the lines that are currently unvisited by `where _visited = false`. It's a saffety measure, in theory we shall never have `_visited = true` outside a transaction. This could be further simplified by deleting these lines inside the transaction.
