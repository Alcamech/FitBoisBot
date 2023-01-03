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
* Individual stat commands
* Dockerize and automate deployment
* Write test
* Create Changelog