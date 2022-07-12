# Index Postgres Data with Algolia

## Prerequisites

- In order to run the application as we walk through, you need to have docker-compose installed. You can follow their [installation guide](https://docs.docker.com/compose/install/).
- You need to have an Algolia account with at least free tier.
- The default postgreSQL exposed port is 25432, it should not be alloacted when starting the app, or the compose file has to be changed.

## Configure Algolia

Create a user if you don't have it already at [Algolia](https://www.algolia.com/).

You can use their [Quickstart guide](https://www.algolia.com/doc/guides/getting-started/quick-start/#sign-up-for-an-algolia-account) to get started.

This demo application was designed to stay below the free tier usage. You can get more info on the limits at Algolia's [pricing page](https://www.algolia.com/pricing/).

With the default settings it will add 20 lines on start, and generate one new line every 10 seeconds.
In the free tier currenlty 10000 lines are free.
It'd take at least 27 hours to fill up the limits.

To set up the Application correctly to your app, you need to go to the [API Keys seection](https://www.algolia.com/account/api-keys/all).
Take note on the following:

- Application ID
- Admin API Key

> Make sure not to commit these values in any public place, and only use the API Key from the backend.

## Getting Started with the Code

### Requirements

The application pack consists of two separate go services and a postgreSQL database.
It's encapsulated in Docker images, and can be started with docker-compose.
Moving forward the post will describe how to use the app with Docker.

If you wish to start it with your environment you'll need to have a running postgres, and a [golang environment](https://go.dev/learn/) to build the services.

### Set up .env

Before you start the application, you have to fill some environment variables that the demo environment can use.
Docker-compose is smart in the way, that it reads a file called `.env` if it's available in the same folder as where it's started. You can read more about this behaviour over [here](https://docs.docker.com/compose/environment-variables/#the-env-file).

The key takeaway is that you need to create a file called `.env` in the same folder as `docker-compose.yml`.

To ease the onboaring you can copy the `.env.example` file and fill out the **...** parts with your data.

The necessary details are:

- `POSTGRES_PASSWORD`: the password to use with the newly started postgreSQL instance
- `ALGOLIA_APP_ID`: the *Application ID* that you noted from the settings page
- `ALGOLIA_API_KEY`: the *Admin API Key*
- `ALGOLIA_INDEX_NAME`: the name of the index that this demo applicaction shall use. If no index exists with the name that you enter it will be created upon first start

### Start the application with Docker-compose

You need to point your terminal to the location of the code and run the following command:

```bash
docker-compose up
```

If you make changes in the application code, don't forget to rebuild the applications.
If you make changes in the service names don't forget to remove orphaned images.
To handle such cases you can start the app with the following command.

```bash
docker-compose up --build --remove-orphans
```

### Troubleshooting

If you follow the instructions above, the demo application *should* start up without any problem.

Although I'll list some common errors that might arise.

In case you've started the docker-compose from a different directory, the application will mostly run, but the environment variables might not set up correctly.
The go services will fail to start, and the postgres image will initialize to the default user.
The solution is to stop the app with `docker-compose down` create the `.env` file, and remove the postgres volume data folder from `./postgres/db`.

In case of a weird python timeout error shows up upon running docker-compose.
The docker daemon service is not running, and the docker-compose program can not connect to it to start the demo application.

### Set up Searchable Indices in Algolia

After starting the application a following actions should happen.

The docker-compose log shall show that a db and a single procuder started with multiple consumers.
The `app.audit_log` table shall have a few auto-generated lines.
The `app.queue_audit_log` table shall not be empty yet.

After the `$DELAY_CONSUMER` amount of seconds, which is 55 by default, the consumer shall load some data to Algolia. And the values shall appear in the given index set by `$ALGOLIA_INDEX_NAME`.

You can search through these lines in the Algolia Dashboard's Search Explorer part.

In order to find the attributes we've just uploaded you need to make them *searchable*.

If you already have some lines the possible values will be listed in the UI.

You need to add:

- action
- contentItemId
- createDateTimestamp
- userId

You can read more about Searchable Attributes in the [docs](https://www.algolia.com/doc/guides/managing-results/must-do/searchable-attributes/).

### Run Queries in Algolia

After you've added the searchable indices you can find the data that you'd expect.

You can even run custom queries like this:

```json
{
  "filters": "createDateTimestamp > 1655127131 AND userId=12 AND action=2"
}
```

It searches or all actions with an ID:2 by user:12 that were added AFTER 1655127131 (Monday, June 13, 2022 1:32:11 PM).
The [epochconverter](https://www.epochconverter.com/) is a handy tool to convert between timestamps and date.

> We're going to deep dive into our simplified implementation in the [next part](./part3.md).
