package com.alcamech.fitboisbot.model;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
public class FitBoiGg {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private Long userId;
    private Long groupId;
    private String year;
    private Integer fastGgCount;

    public FitBoiGg() {}

    public FitBoiGg(Long userId, Long groupId, String year, Integer fastGgCount) {
        this.userId = userId;
        this.groupId = groupId;
        this.year = year;
        this.fastGgCount = fastGgCount;
    }

    public Long getUserId() {
        return userId;
    }

    public void setUserId(Long userId) {
        this.userId = userId;
    }

    public Long getGroupId() {
        return groupId;
    }

    public void setGroupId(Long groupId) {
        this.groupId = groupId;
    }

    public String getYear() {
        return year;
    }

    public void setYear(String year) {
        this.year = year;
    }

    public Integer getFastGgCount() {
        return fastGgCount;
    }

    public void setFastGgCount(Integer fastGgCount) {
        this.fastGgCount = fastGgCount;
    }
}
