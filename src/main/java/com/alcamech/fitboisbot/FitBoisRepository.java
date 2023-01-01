package com.alcamech.fitboisbot;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;

import java.util.List;

// This will be AUTO IMPLEMENTED by Spring into a Bean called fitboisRepository
// CRUD refers Create, Read, Update, Delete
public interface FitBoisRepository extends CrudRepository<FitBoiRecord, Integer> {

    @Override
    List<FitBoiRecord> findAll();

    @Query("SELECT DISTINCT f.name FROM FitBoiRecord f")
    List<String> findDistinctName();

    Long countByName(String name);
}