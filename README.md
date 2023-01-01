# FitBois Telegram Bot

## System Requirements 

* Java 11
* Maven
* MySQL

## Dependencies

To fetch dependencies `mvn clean install`

## Build, Run, Test

`dev.properties` file needed

To build `maven compile`

To run `mvn spring-boot:run` or click Run from `FitBoisBotApplication.java`

Refer to `dev.properties` for database credentials. 
Please create a databaes called `ftibois` and run the script `resources/fitbois.sql`

## TODO

* Error handling
* Individual stat commands
* Dockerize 
* Clean up