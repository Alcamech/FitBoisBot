package com.alcamech.fitboisbot.respository;

import com.alcamech.fitboisbot.model.FitBoiUser;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;

import java.util.List;

public interface UserRepository extends CrudRepository<FitBoiUser, Long> {

    @Query("SELECT DISTINCT u.name FROM FitBoiUser u WHERE u.groupId = :groupId")
    List<String> findDistinctName(Long groupId);

    List<FitBoiUser> findFitBoiUsersByGroupId(Long groupId);
}
