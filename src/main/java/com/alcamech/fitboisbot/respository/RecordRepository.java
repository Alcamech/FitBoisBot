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

    @Query("SELECT COUNT(*) FROM FitBoiRecord WHERE userId = :userId AND year = YEAR(CURDATE()) AND month = :month")
    Long countByUserIdWithCurrentYearAndMonth(@Param("userId") Long userId, @Param("month") String month);

    @Query("SELECT MAX(COUNT(f)) FROM FitBoiRecord f WHERE f.year =  YEAR(CURDATE()) AND f.month = :month GROUP BY f.userId")
    Long findMaxActivityCountByYearAndMonth(@Param("month") String month);

    @Query("SELECT f.userId FROM FitBoiRecord f WHERE f.year = YEAR(CURDATE()) AND f.month = :month GROUP BY f.userId HAVING COUNT(f) = :maxCount")
    List<Long> findAllUsersWithMaxCount(@Param("month") String month, @Param("maxCount") Long maxCount);
}