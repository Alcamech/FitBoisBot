// Package constants provides centralized constants for the FitBoisBot application.
package constants

import "time"

// Default configuration values
const (
	// DefaultTimezone is the default timezone used when none is set for a group
	DefaultTimezone = "America/New_York"

	// ActivityNameMaxLength is the maximum allowed length for activity names
	ActivityNameMaxLength = 50

	// DatePartLength is the required length for date parts (MM, DD)
	DatePartLength = 2
)

// Token and reward amounts
const (
	// MonthlyRewardAmount is the number of tokens awarded to the most active user each month
	MonthlyRewardAmount = 100
)

// Scheduler configuration
const (
	// SchedulerCheckInterval is how often the scheduler checks for monthly awards
	SchedulerCheckInterval = 1 * time.Minute

	// AwardHour is the hour (24h format) when monthly awards are processed
	AwardHour = 9

	// AwardDay is the day of the month when monthly awards are processed
	AwardDay = 1
)

// Date validation ranges
const (
	MinMonth = 1
	MaxMonth = 12
	MinDay   = 1
	MaxDay   = 31
)

// Challenge configuration
const (
	// ChallengeJoinWindow is the time allowed for participants to join before auto-cancel
	ChallengeJoinWindow = 12 * time.Hour

	// ChallengeDuration is the total duration of an active challenge
	ChallengeDuration = 14 * 24 * time.Hour

	// ChallengeSchedulerInterval is how often the challenge scheduler checks for auto-actions
	ChallengeSchedulerInterval = 1 * time.Hour

	// ChallengePageSize is the number of challenges shown per page in /listchallenges
	ChallengePageSize = 5
)

// Challenge difficulty multipliers
const (
	MultiplierEasy     = 0.5
	MultiplierModerate = 1.0
	MultiplierHard     = 1.5
)

// Challenge status values
const (
	ChallengeStatusPending   = "pending"
	ChallengeStatusActive    = "active"
	ChallengeStatusCompleted = "completed"
	ChallengeStatusCancelled = "cancelled"
)

// Challenge difficulty values
const (
	DifficultyEasy     = "easy"
	DifficultyModerate = "moderate"
	DifficultyHard     = "hard"
)

// Challenge status badges (emojis)
const (
	StatusBadgePending   = "🕐"
	StatusBadgeActive    = "🔥"
	StatusBadgeCompleted = "✅"
	StatusBadgeCancelled = "❌"
)

// Challenge difficulty badges (emojis)
const (
	DifficultyBadgeEasy     = "⭐"
	DifficultyBadgeModerate = "⭐⭐"
	DifficultyBadgeHard     = "⭐⭐⭐"
)

// Button labels
const (
	ButtonPrevious   = "« Previous"
	ButtonNext       = "Next »"
	ButtonBackToList = "« Back to List"
)

