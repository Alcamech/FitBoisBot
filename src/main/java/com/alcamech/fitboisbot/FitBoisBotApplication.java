package com.alcamech.fitboisbot;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.PropertySource;
import org.springframework.context.annotation.PropertySources;
import org.springframework.scheduling.annotation.EnableScheduling;

@SpringBootApplication
@PropertySources({
        @PropertySource("classpath:application.properties"),
        @PropertySource("classpath:dev.properties")
})
@EnableScheduling
public class FitBoisBotApplication {

    public static void main(String[] args) {
        SpringApplication.run(FitBoisBotApplication.class, args);
    }

}
