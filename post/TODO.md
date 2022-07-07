# TODO

> This folder shall be removed from the history upon publishing by rewriting git history.


- ? UI
- ? measurements
- why not notify/listen
- further enhancements
- https://www.algolia.com/account/api-keys/all?applicationId=85S5IGY41E

- need to set  up searhable attributes
- need to create an index
- what happens on rror

- queries in algolia

1. Usecase architecture, Planned solution, small demo poc and not the actual solution
1. Start it at the user
   setup docker + Getting started in Algolia, summary
    - objectID, timestamp
    - free tier 10000
   extra implementations???
1. Detailed summary, db structure, go code, docker
1. UI Extra featurees, measurements

## Part1

<!-- We store different kind of actions that users make in an audit_log table. It would be beneficial to have an admin screen where we could look into this data and quickly check the events were made without generating traffic on the original database.
For example looking into a single user's contribution (anonimized of course), or what changes were made in a particular item e.g shown in a timeline

Our admin users would like to find out how the users use this application
We need to now what happens and when to certain items in the database
We need to be able to extract this database and look it up quickly wiithout impacting the original database
We need a few kind of reports
    - what happened with an exact record through its lifetime to investigate potential mishandling
    - What reports were impacted by a user
    - What happened between two points in time
We developed a cloud service to speed up the life of everyone who works with BI reports.
 -->

### Architeture

<!-- 
We have a go process.
We store our data in postgresql.
AuditData is not modified after creation, we only send it to the db, depending on active usebase and usage it can grow large -->

<!-- ### Constraints

we have a project that needs minimal moving parts
we can not expose the database to external auditors
We need to have a simple deployment pipeline -->

<!-- ## Current challenges

We don't want to put much stress on the database when querying these items
We need to separate the audit process from the database -->

### Proposed solution

<!-- 
We index our data to Algolia.
We can use postgres to feed data into Algolia
We need to have a simple persisted task queue without external dependencies.
We wanted to focus on postgres capabilities to make the backend replacable. -->

## Part3


<!-- ## Impovement ideas

- Use listen/notify for realtime communication
  - Drawback: if the consumer is not connected and new lines arrive they won't be considered
- Support multiple parallel consumers:
  - the `_visited` flag shall be an identifier instead of a bool so that each worker acn work on different lines
  - we need to handle what happens with the lines that are still held by consumers that have been stopped
  - the visited mark shall lock or apply immediately so that others can not claim it

## TODO

- Create human-readable and better searchable data for the ID-s that we can upload into Algolia
- Connect to the db on app init not on every request
- Use transaction in the consumer if we failed to upload data to index we shall roll back
- Handle what happens if we use a single instance and we have visited lines in the db on start -->

<!-- 
For every new line a trigger copies the data into an otheer table that acts as a queue.
A separate application reads the last N items of that table periodically.
We can not ensure the availability off this application, listen and notify are not feeasible.

The data is anonimized, we track them by a generated ID

I'll present a simplified example as an example with a github reepo, and go through how it works.

I use plain docker images
For Go services I use 2stagee build to create the build atifact in the first step and run the application in the second.
For the sake of simplicity I ddon't use any other extensions other than the Algolia connector and libpq to conneect to postgres.
I use the official postgres docker images 
-->
<!-- 
## Current table structure

## Solutions supported by PostgreSQL

## Algolia -->