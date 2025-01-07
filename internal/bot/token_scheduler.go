package bot

import (
	"log"
	"time"
	_ "time/tzdata"
)

func ScheduleMonthlyTokenAwards() {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for range ticker.C {
		processMonthlyTokenAwards()
	}
}

func awardTokensForGroup(groupID int64, previousMonth string) {
	users, err := activityRepo.GetUsersWithActivities(groupID)
	if err != nil {
		log.Printf("Failed to fetch users with activities for group %d: %v", groupID, err)
		return
	}

	var topUserID int64
	var maxActivities int64

	for _, userID := range users {
		count, err := activityRepo.GetActivityCountByUserIdAndMonth(userID, groupID, previousMonth)
		if err != nil {
			log.Printf("Failed to fetch activity count for user %d: %v", userID, err)
			continue
		}

		if count > maxActivities {
			maxActivities = count
			topUserID = userID
		}
	}

	if topUserID == 0 {
		log.Printf("No activities found for group %d in %s", groupID, previousMonth)
		return
	}

	year := previousMonth[:4]
	if err := tokenRepo.IncrementTokens(topUserID, groupID, year, 10); err != nil {
		log.Printf("Failed to award tokens to user %d: %v", topUserID, err)
		return
	}

	log.Printf("Awarded 10 tokens to user %d for group %d in %s", topUserID, groupID, previousMonth)
}

func processMonthlyTokenAwards() {
	groups, err := groupRepo.GetAllGroups()
	if err != nil {
		log.Printf("Failed to fetch groups: %v", err)
		return
	}

	for _, group := range groups {
		location, err := time.LoadLocation(group.Timezone)
		if err != nil {
			log.Printf("Failed to load timezone %s for group %d: %v", group.Timezone, group.ID, err)
			continue
		}

		now := time.Now().In(location)
		if now.Day() != 1 || now.Hour() != 9 || now.Minute() != 0 {
			continue
		}
		// log.Printf("Testing: Current time in timezone %s for group %d: %v", group.Timezone, group.ID, now)

		activityTime := now.AddDate(0, -1, 0)
		activityMonth := activityTime.Format("01")  // MM format
		activityYear := activityTime.Format("2006") // YYYY format

		userID, _, err := activityRepo.GetMostActiveUserForMonth(group.ID, activityMonth, activityYear)
		if err != nil {
			log.Printf("No activities found for group %d in %s/%s: %v", group.ID, activityMonth, activityYear, err)
			continue
		}

		currentYear := now.Format("2006")
		rewardAmount := 100
		err = tokenRepo.IncrementTokens(userID, group.ID, currentYear, rewardAmount)
		if err != nil {
			log.Printf("Failed to award tokens to user %d for group %d: %v", userID, group.ID, err)
			continue
		}

		user, err := userRepo.FindByID(userID)
		if err != nil {
			log.Printf("Failed to fetch user %d for group %d: %v", userID, group.ID, err)
			continue
		}

		sendMonthlyAwardMessage(bot, group.ID, user.Name, activityMonth, activityYear, rewardAmount)
	}
}
