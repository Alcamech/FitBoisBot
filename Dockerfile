FROM openjdk:11-jdk-alpine
LABEL maintainer="alcamech@gmail.com"
VOLUME /app
ADD target/FitBoisBot-1.0.0.jar app.jar
ENTRYPOINT ["java", "-jar","/app.jar"]