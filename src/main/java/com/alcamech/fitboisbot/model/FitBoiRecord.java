package com.alcamech.fitboisbot.model;

import javax.persistence.*;

@Entity
public class FitBoiRecord {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;
    private Long userId;
    private String activity;
    private String month;
    private String day;
    private String year;

    public FitBoiRecord() {}

    public FitBoiRecord(Long userId, String activity, String month, String day, String year) {
        this.userId = userId;
        this.activity = activity;
        this.month = month;
        this.day = day;
        this.year = year;
    }
}
