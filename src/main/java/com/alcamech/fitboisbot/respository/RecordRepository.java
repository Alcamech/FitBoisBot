package com.alcamech.fitboisbot.respository;
import com.alcamech.fitboisbot.model.FitBoiRecord;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;
import org.springframework.data.repository.query.Param;

import java.util.List;

public interface RecordRepository extends CrudRepository<FitBoiRecord, Integer> {

    @Query("SELECT DISTINCT f.userId FROM FitBoiRecord f")
    List<Long> findDistinctRecords();

    @Query("SELECT COUNT(*) FROM FitBoiRecord WHERE userId = :userId AND year = YEAR(CURDATE())")
    Long countByUserIdAndCurrentYear(Long userId);

    @Query("SELECT COUNT(*) FROM FitBoiRecord WHERE userId = :userId AND year = YEAR(CONVERT_TZ(CURDATE(), '+00:00', '-05:00')) AND month = MONTH(CONVERT_TZ(CURDATE(), '+00:00', '-05:00'))")
    Long countByUserIdWithCurrentYearAndMonth(Long userId);

    @Query("SELECT f.userId, COUNT(*) as recordCount FROM FitBoiRecord f WHERE f.year = :year AND f.month = :month GROUP BY f.userId ORDER BY recordCount DESC LIMIT 1")
    Long findMostActiveUserByYearAndMonth(@Param("year") int year, @Param("month") int month);
}