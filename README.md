# FitBois Telegram Bot

## System Requirements 

* Java 11
* Maven
* MySQL

## Dependencies

To fetch dependencies `mvn clean install`

## Contributing

`dev.properties` file needed

FitBoisDevBot is used as a test, dev, and stage environment.

## Build, Run, Test

`dev.properties` file needed

To build `maven compile`

To run `mvn spring-boot:run` or click Run from `FitBoisBotApplication.java`

Refer to `dev.properties` for database credentials. 
Please create a databaes called `ftibois` and run the script `resources/fitbois.sql`

## TODO

* see TODOs in code
* Error handling
* Dockerize and automate deployment
* Write test
* Create Changelog

# Planned Features

### Bot command to show individual stats

`/stats` - should show the individual their personal stats that are being tracked

* Activity count
* count by month
* count by activity type
* fastest gg count

### Bot command to show individual "fastest gg" count

`/fastgg` - should show everyone's count of "fastest gg"

### Bot command to submit challenge completion [draft]

`/challenge` - should allow you to submit a challenge

first optional parameter is `1=complete`

second parameter is wager amount 

message should be in the format of `/challenge 1 100`