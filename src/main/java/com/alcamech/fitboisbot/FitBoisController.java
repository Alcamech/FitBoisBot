package com.alcamech.fitboisbot;

import org.springframework.stereotype.Controller;

import java.util.List;


@Controller
public class FitBoisController{
    private final FitBoisRepository fitBoisRepository;

    public FitBoisController(FitBoisRepository fitBoisRepository) {
        this.fitBoisRepository = fitBoisRepository;
    }

    public void addNewRecord(String name, String activity, String month, String day, String year) {
        FitBoiRecord FitBoiRecord = new FitBoiRecord(name, activity, month, day, year);
        fitBoisRepository.save(FitBoiRecord);
    }

    public List<FitBoiRecord> getAllFitBoiRecords() {
        return fitBoisRepository.findAll();
    }

    public List<String> getFitBois() {
        return fitBoisRepository.findDistinctName();
    }

    public Long getCountByName(String name) {
        return fitBoisRepository.countByName(name);
    }
}