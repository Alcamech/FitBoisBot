package com.alcamech.fitboisbot.respository;

import com.alcamech.fitboisbot.model.FitBoiUser;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

public interface UserRepository extends CrudRepository<FitBoiUser, Long> {

    @Modifying
    @Transactional
    @Query("UPDATE FitBoiUser u SET u.fastGgCount = u.fastGgCount +1 WHERE u.id = :id AND u.groupId = :groupId")
    void updateGgCount(Long id, Long groupId);

    @Modifying
    @Transactional
    @Query("UPDATE FitBoiUser u SET u.fastGgCount = 0 WHERE u.id = :id AND u.groupId = :groupId")
    void initializeGgCount(Long id, Long groupId);

    @Query("SELECT DISTINCT u.name FROM FitBoiUser u WHERE u.groupId = :groupId")
    List<String> findDistinctName(Long groupId);

    List<FitBoiUser> findFitBoiUsersByGroupId(Long groupId);

    @Query("SELECT u.fastGgCount FROM FitBoiUser u WHERE u.id = :id")
    Integer fetchFastGgCountById(Long id);
}
