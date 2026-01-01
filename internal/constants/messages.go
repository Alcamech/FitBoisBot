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
  • ` + "`/help`" + ` or ` + "`/h`" + ` \- Show this help message\.
  • ` + "`/fastgg`" + ` or ` + "`/gg`" + ` \- Display the Fast GG leaderboard\.
  • ` + "`/tokens`" + ` or ` + "`/t`" + ` \- Display Fitboi Token balances\.
  • ` + "`/balance`" + ` or ` + "`/b`" + ` \- Show your token balance\.
  • ` + "`/leaderboard`" + ` or ` + "`/l`" + ` \- Display monthly activity leaderboard\.
  • ` + "`/timezone`" + ` or ` + "`/tz`" + ` \- Show or set timezone\.
  • ` + "`/challenge`" + ` or ` + "`/c`" + ` \- Create a fitness challenge \(use ` + "`/help challenge`" + ` for details\)\.

*🗓️ Activity Format:*
Post activities in the format:
` + "`activity\\-MM\\-DD\\-YYYY`" + ` or ` + "`activity\\-MM\\-DD\\-YY`" + `\.

*🚑 Support:*
For any issues, reach out to our support team\.

💪 *Stay fit and keep posting your progress\!*

__%s__`

// Template strings for dynamic messages
const (
	MonthlyAwardTemplate     = "Monthly counts have been reset\n\nCongratulations ⭐️ %s ⭐️ for being the most active user for %s/%s 🏆\n\nHere's your reward. 💰 You've won %d FitBoi Tokens! 💰"
	MonthlyAwardTieTemplate  = "Monthly counts have been reset\n\nCongratulations ⭐️ %s ⭐️ for tying as the most active users for %s/%s with %d activities each! 🏆\n\nThe %d FitBoi Token reward has been split %d ways. 💰 You've each won %d FitBoi Tokens! 💰"
	TokenLeaderboardTemplate = "🏆 Token Leaderboard for %s 🏆\n"
	NoTokensTemplate         = "No tokens awarded yet ! 🏆"
	FastGGLeaderboardPrefix  = "Fastest GGs 😎 "
)

// Challenge messages
const (
	MsgChallengeCreated            = "Challenge created! Waiting for participants to join within 12 hours."
	MsgChallengeJoined             = "You've joined the challenge!"
	MsgChallengeActivated          = "Challenge activated! 14 days to compete."
	MsgChallengeCompleted          = "Challenge completed!"
	MsgChallengeAutoCompleted      = "Challenge completed after 14 days! Creator: use /completechallenge to award winners."
	MsgChallengeCancelled          = "Challenge cancelled. Wager refunded."
	MsgChallengeAutoCancelled      = "Challenge auto-cancelled - no one joined within 12 hours. Wager refunded."
	MsgChallengeCannotCancel       = "Cannot cancel - challenge is already active."
	MsgChallengeAlreadyExists      = "A challenge is already active or pending in this group."
	MsgChallengeNotFound           = "No challenge found."
	MsgChallengeInsufficientTokens = "Insufficient tokens for this wager."
	MsgChallengeAlreadyParticipant = "You're already in this challenge."
	MsgChallengeInvalidDifficulty  = "Invalid difficulty. Use: easy, moderate, or hard"
	MsgChallengeOnlyCreator        = "Only the challenge creator can do that."
	MsgChallengeNoLongerActive     = "Challenge is completed. No more scoring allowed."
	MsgChallengeScoreUpdated       = "Score updated!"
	MsgChallengeScoresUpdated      = "Scores updated for %d participant(s)!"
	MsgChallengeInvalidWager       = "Wager must be greater than 0."
	MsgChallengeInvalidFormat      = "Invalid format. Use: /challenge [difficulty] [wager] [title] [description]"
	MsgChallengeNoWinners          = "Challenge completed with no winners. All wagers burned."
	MsgChallengeWinnersAwarded     = "Winners awarded!"
	MsgChallengeAlreadyAwarded     = "Winners have already been awarded for this challenge."
	MsgChallengeUserNotParticipant = "One or more mentioned users are not participants."
	MsgChallengeNoChallenges       = "No challenges found for this group."
)

// Challenge help text
const MsgChallengeHelp = `<b>Challenge System</b>

<b>Overview:</b>
Create fitness challenges with token rewards based on difficulty!

<b>Commands:</b>
• <code>/challenge</code> or <code>/c</code> [difficulty] [wager] [title] [description]
  Create a new challenge (easy/moderate/hard)

• <code>/joinchallenge</code> or <code>/jc</code> [wager]
  Join the current challenge

• <code>/viewchallenge</code> or <code>/vc</code> [id]
  View current or specific challenge

• <code>/listchallenges</code> or <code>/lc</code> [page]
  Browse all challenges with pagination

• <code>/score</code> or <code>/s</code> @user1 [points1] @user2 [points2]...
  Award points to participants (creator only)
  Examples:
    /s @john 5 @rob 6 (individual scores)
    /s @john @rob @mark 10 (shared score)

• <code>/cancelchallenge</code> or <code>/cc</code>
  Cancel pending challenge (creator only)

• <code>/completechallenge</code> or <code>/done</code> [@winner1 @winner2...]
  Complete challenge and award winners (creator only)

<b>Difficulty Multipliers:</b>
• Easy (0.5x): Winners get wager + 50%
• Moderate (1.0x): Winners get wager + 100%
• Hard (1.5x): Winners get wager + 150%

<b>Rules:</b>
• 12-hour window to get participants
• 14-day challenge duration
• Creator determines winners
• Scores are informational only`
