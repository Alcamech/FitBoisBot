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

    @Query(value = "SELECT MAX(activity_count) FROM (SELECT COUNT(*) AS activity_count FROM fit_boi_record WHERE year = YEAR(CURDATE()) " +
            "AND month = :month GROUP BY user_id) AS subquery", nativeQuery = true)
    Long findMaxActivityCountByYearAndMonth(@Param("month") String month);

    @Query(value = "SELECT user_id FROM fit_boi_record WHERE year = YEAR(CURDATE()) AND month = :month GROUP BY user_id HAVING COUNT(*) = :maxCount", nativeQuery = true)
    List<Long> findAllUsersWithMaxCount(@Param("month") String month, @Param("maxCount") Long maxCount);
}