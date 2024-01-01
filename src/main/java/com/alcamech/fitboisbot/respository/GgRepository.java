package com.alcamech.fitboisbot.respository;

import com.alcamech.fitboisbot.model.FitBoiGg;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;
import org.springframework.transaction.annotation.Transactional;

public interface GgRepository extends CrudRepository<FitBoiGg, Long> {
    @Modifying
    @Transactional
    @Query("UPDATE FitBoiGg ugg SET ugg.fastGgCount = ugg.fastGgCount + 1 " +
            "WHERE ugg.userId = :userId AND ugg.groupId = :groupId AND ugg.year = YEAR(CURDATE())")
    void updateGgCountForCurrentYear(Long userId, Long groupId);

    @Query("SELECT ugg.fastGgCount FROM FitBoiGg ugg WHERE ugg.userId = :userId AND ugg.year = YEAR(CURDATE())")
    Integer fetchFastGgCountByIdAndCurrentYear(Long userId);
}
