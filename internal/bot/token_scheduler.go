package bot

import (
	"log/slog"
	"time"
	_ "time/tzdata"

	"github.com/Alcamech/FitBoisBot/internal/constants"
)

func (s *BotService) ScheduleMonthlyTokenAwards() {
	ticker := time.NewTicker(constants.SchedulerCheckInterval) // Check every minute
	defer ticker.Stop()

	for range ticker.C {
		s.processMonthlyTokenAwards()
	}
}

func (s *BotService) processMonthlyTokenAwards() {
	groups, err := s.groupStore.GetAll()
	if err != nil {
		slog.Error("Failed to fetch groups", "error", err)
		return
	}

	for _, group := range groups {
		location, err := time.LoadLocation(group.Timezone)
		if err != nil {
			slog.Error("Failed to load timezone", "error", err, "timezone", group.Timezone, "group_id", group.ID)
			continue
		}

		now := time.Now().In(location)
		if now.Day() != constants.AwardDay || now.Hour() != constants.AwardHour || now.Minute() != 0 {
			continue
		}

		activityTime := now.AddDate(0, -1, 0)
		activityMonth := activityTime.Format("01")  // MM format
		activityYear := activityTime.Format("2006") // YYYY format

		userIDs, activityCount, err := s.activityStore.GetMostActiveUsersForMonth(group.ID, activityMonth, activityYear)
		if err != nil {
			slog.Info("No activities found for group", "group_id", group.ID, "month", activityMonth, "year", activityYear)
			continue
		}

		currentYear := now.Format("2006")
		totalRewardAmount := constants.MonthlyRewardAmount
		numWinners := len(userIDs)

		var rewardPerUser int
		if numWinners == 1 {
			rewardPerUser = totalRewardAmount
		} else {
			rewardPerUser = totalRewardAmount / numWinners
		}

		var winnerNames []string
		for _, userID := range userIDs {
			err = s.tokenStore.IncrementTokens(userID, group.ID, currentYear, rewardPerUser)
			if err != nil {
				slog.Error("Failed to award tokens", "error", err, "user_id", userID, "group_id", group.ID)
				continue
			}

			user, err := s.userStore.FindByID(userID)
			if err != nil {
				slog.Error("Failed to fetch user", "error", err, "user_id", userID, "group_id", group.ID)
				continue
			}
			winnerNames = append(winnerNames, user.Name)
		}

		s.sendMonthlyAwardMessage(group.ID, winnerNames, activityMonth, activityYear, rewardPerUser, numWinners, activityCount)
	}
}
