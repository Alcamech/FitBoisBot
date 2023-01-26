package com.alcamech.fitboisbot.respository;
import com.alcamech.fitboisbot.model.FitBoiRecord;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;

import java.util.List;

public interface FitBoisRepository extends CrudRepository<FitBoiRecord, Integer> {

    @Query("SELECT DISTINCT f.name FROM FitBoiRecord f")
    List<String> findDistinctName();

    Long countByName(String name);
}