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

