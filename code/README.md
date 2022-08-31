# Algolia application - Go backend

## Requirements

- Docker / Go 1.18
- [Configured .env file](../../README.md)

## How to run

```bash
# navigate to root
cp .env.example .env
# edit env file with your data
# docker-compoes will read ths env file but only if it's starrted from root
docker-compose up --build --remove-orphans
```

## Architecture

### Postgres Database

We run the [official docker image](https://hub.docker.com/_/postgres) of [postgres](./postgres).

Upon first start it:

1. creates a db folder (or optionally a volume) in ordder to persist the data
1. Runs the init sqls that are inside [docker-entrypoint-initdb.d](./postgres/docker-entrypoint-initdb.d).

#### Init Scripts

- Create `app` schema.
- Create `app.audit_log` fact table
- Create `app.queue_audit_log` queue table
- Create `audit_insert_trigger` to replicate data sent to fact table to the queue
- Generate random data into the fact table. Its changing parts are extracted into variables. The randomizer has an init seed, it generates the same data across multiple recreations.

### Producer

- Periodically generate a new random log line into the fact table.
- The insert trigger will copy this data into the queue.

### Consumer

- Periodically reads the queue for new data, and then does the following:

  1. marks the last lines as visited and reads them
  1. Uploads the latest lines into Algolia
  1. Clears the visited lines from the queue

## Impovement ideas

- Use listen/notify for realtime communication
  - Drawback: if the consumer is not connected and new lines arrive they won't be considered
- Support multiple parallel consumers:
  - the `_visited` flag shall be an identifier instead of a bool so that each worker acn work on different lines
  - we need to handle what happens with the lines that are still held by consumers that have been stopped
  - the visited mark shall lock or apply immediately so that others can not claim it
