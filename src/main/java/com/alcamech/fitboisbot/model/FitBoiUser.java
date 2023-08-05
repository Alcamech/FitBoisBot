package com.alcamech.fitboisbot.model;

import javax.persistence.Entity;
import javax.persistence.Id;
import javax.persistence.Version;

@Entity
public class FitBoiUser {

    @Id
    private Long id;
    private String name;
    private Long groupId;
    private Integer fastGgCount;
    @Version
    private Integer version;

    public FitBoiUser() {}

    public FitBoiUser(Long id, String name, Long groupId) {
        this.id = id;
        this.name = name;
        this.groupId = groupId;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public Long getGroupId() {
        return groupId;
    }

    public void setGroupId(Long groupId) {
        this.groupId = groupId;
    }

    public Integer getFastGgCount() {
        return fastGgCount;
    }

    public void setFastGgCount(Integer fastGgCount) {
        this.fastGgCount = fastGgCount;
    }

    @Override
    public String toString() {
        return "FitBoiUser{" +
                "id=" + id +
                ", name='" + name + '\'' +
                ", groupId=" + groupId +
                ", fastGgCount=" + fastGgCount +
                ", version=" + version +
                '}';
    }
}
