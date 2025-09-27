package constants

// Error and info messages
const (
	MsgGroupTimezoneError    = "Failed to fetch group timezone. Please set the timezone using /settimezone."
	MsgInvalidTimezone       = "Invalid timezone set for this group. Please update the timezone using /settimezone."
	MsgTimezoneSetSuccess    = "Timezone updated successfully to "
	MsgTimezoneSetFailed     = "Failed to set timezone. Please try again."
	MsgTimezoneRequired      = "Please specify a timezone. Example: /timezone America/New_York"
	MsgTimezoneInvalid       = "Invalid timezone. Please use a valid IANA timezone. Example: /timezone America/New_York"
	MsgLeaderboardFailed     = "Failed to fetch leaderboard. Please try again later."
	MsgActivityFormatError   = "Wrong activity format. use activity-MM-DD-YYYY or activity-MM-DD-YY"
	MsgFastestGG             = "Fastest gg in the west"
	MsgNoActivitiesRecorded  = "No activity recorded."
	MsgNoFastGGsRecorded     = "No Fast GGs recorded yet!"
	MsgActivityCountsUpdated = "Activity counts updated: "
)

// Help text template - %s will be replaced with version
const MsgHelpText = `*🏋️ Welcome to the FitBois Bot\!*
Here are the available commands:

*📜 Commands:*
  • ` + "`/help`" + ` \- Show this help message\.
  • ` + "`/fastgg`" + ` \- Display the Fast GG leaderboard your group\.
  • ` + "`/tokens`" + ` \- Display Fitboi Token balances\.
  • ` + "`/timezone`" + ` \- Show current timezone or set new one \(e\.g\., ` + "`/timezone America/New_York`" + `\)\.

*🗓️ Activity Format:*
Post activities in the format:
` + "`activity\\-MM\\-DD\\-YYYY`" + ` or ` + "`activity\\-MM\\-DD\\-YY`" + `\.

*🚑 Support:*
For any issues, reach out to our support team\.

💪 *Stay fit and keep posting your progress\!*

__%s__`

// Template strings for dynamic messages
const (
	MonthlyAwardTemplate     = "Month counts have been reset\n\nCongratulations ⭐️ %s ⭐️ for being the most active user for %s/%s 🏆\n\nHere's your reward. 💰 You've won %d FitBoi Tokens! 💰"
	TokenLeaderboardTemplate = "🏆 Token Leaderboard for %s 🏆\n"
	NoTokensTemplate         = "No tokens awarded yet ! 🏆"
	FastGGLeaderboardPrefix  = "Fastest GGs 😎 "
)

