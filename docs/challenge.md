# Challenge System Design

## Overview

Challenges can be started by anyone in the group. The user who starts it wagers tokens, users who join also wager. You need at least one person to join to activate the challenge. Only one challenge can be active at a time per group.

### Core Concepts

- **Creator-managed scoring**: Only creator can award points to participants
- **Flexible rules**: Title/description are informational, actual rules and details are expected to be determined by the group.
- **Creator determines winners**: After 14 days, creator should mention winners in `/completechallenge` command
- **Difficulty multipliers**: Each winner gets back wager + (wager × multiplier)
- **Cancellation**: Creator can cancel pending challenges (before anyone joins). Active challenges cannot be cancelled.

### Difficulty System

**Total payout formula**: `wager + (wager × multiplier)`

- **Easy** (0.5x): Winners receive wager + (wager × 0.5) = 1.5x their wager
  - Example: 100 wagered → 100 + 50 = 150 tokens received (50 profit)
- **Moderate** (1.0x): Winners receive wager + (wager × 1.0) = 2x their wager
  - Example: 100 wagered → 100 + 100 = 200 tokens received (100 profit)
- **Hard** (1.5x): Winners receive wager + (wager × 1.5) = 2.5x their wager
  - Example: 100 wagered → 100 + 150 = 250 tokens received (150 profit)
- **Difficult** (2.0x): Winners receive wager + (wager × 2.0) = 3x their wager
  - Example: 100 wagered → 100 + 200 = 300 tokens received (200 profit)

### Winner Selection

After 14 days, the challenge ends and the **creator determines winners** by mentioning them:

**Single Winner:**

```
/completechallenge @john
```

John gets his payout (wager + (wager x multiplier)), Everyone else loses the wager.

**Multiple winners:**

```
/completechallenge @john @sarah @mike alice
```

All four receive their individual payouts (wager + wager × multiplier each). Everyone else loses their wager.

**Everyone who met goal:**
Creator can reward any number of participants who met the challenge criteria.

**No winners:**

```
/completechallenge
```

No one mentioned = no one gets tokens, all wagers are lost/burned.

### Token Economics

- **Losers**: Lose their entire wager (tokens burned)
- **Winners**: Receive back original wager + bonus (wager × multiplier)
- **All difficulties**: Winners always profit, just different amounts
- **Minting**: System mints (wager × multiplier) for each winner
- **Net effect**: Losers' burned tokens offset some/all of winners' minted tokens

### Time Constraints

- **Creator auto-joins**: Creator starts challenge with wager and is automatically a participant
- **12-hour join window**: Challenge auto-expires (refund creator) if no other participant joins within 12 hours
- **14-day duration**: Challenges run for exactly 14 days, then creator awards winners. Challenges can no longer be scored after 14 days, its up to the creator to do a completion for awards.
- **Activation**: Challenge becomes "active" when at least one other person joins (minimum 2 participants total)

## SQL Schema

### Main Challenges Table

```sql
CREATE TABLE challenges (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL REFERENCES groups(id),
    creator_id BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    difficulty VARCHAR(20) NOT NULL CHECK (difficulty IN ('easy', 'moderate', 'hard', 'difficult')),
    multiplier DECIMAL(3,1) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'completed', 'cancelled')),
    start_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_group_status (group_id, status),
    INDEX idx_creator (creator_id)
);
```

### Challenge Participants Table

```sql
CREATE TABLE challenge_participants (
    id BIGSERIAL PRIMARY KEY,
    challenge_id BIGINT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    wager_amount INT NOT NULL CHECK (wager_amount > 0),
    score INT NOT NULL DEFAULT 0,
    is_winner BOOLEAN DEFAULT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY unique_challenge_user (challenge_id, user_id),
    INDEX idx_user (user_id)
);
```

## Go Models

```go
type Challenge struct {
    ID          int64      `gorm:"primaryKey;autoIncrement"`
    GroupID     int64      `gorm:"not null;index:idx_group_status"`
    CreatorID   int64      `gorm:"not null;index:idx_creator"`
    Title       string     `gorm:"size:255;not null"`
    Description string     `gorm:"type:text"`
    Difficulty  string     `gorm:"size:20;not null;check:difficulty IN ('easy','moderate','hard','difficult')"`
    Multiplier  float64    `gorm:"not null"`
    Status      string     `gorm:"size:20;not null;default:'pending';check:status IN ('pending','active','completed','cancelled');index:idx_group_status"`
    StartDate   *time.Time // When first participant joins (activates challenge)
    CreatedAt   time.Time  `gorm:"autoCreateTime"`
    UpdatedAt   time.Time  `gorm:"autoUpdateTime"`

    Group        Group                   `gorm:"foreignKey:GroupID"`
    Creator      User                    `gorm:"foreignKey:CreatorID"`
    Participants []ChallengeParticipant  `gorm:"foreignKey:ChallengeID"`
}

type ChallengeParticipant struct {
    ID          int64     `gorm:"primaryKey;autoIncrement"`
    ChallengeID int64     `gorm:"not null;uniqueIndex:unique_challenge_user"`
    UserID      int64     `gorm:"not null;index:idx_user;uniqueIndex:unique_challenge_user"`
    WagerAmount int       `gorm:"not null;check:wager_amount > 0"`
    Score       int       `gorm:"not null;default:0"`
    IsWinner    *bool     `gorm:"default:null"`
    JoinedAt    time.Time `gorm:"not null;autoCreateTime"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`

    Challenge Challenge `gorm:"foreignKey:ChallengeID"`
    User      User      `gorm:"foreignKey:UserID"`
}
```

## Key Design Decisions

### Status Flow

- **pending**: Challenge created by user, waiting for second participant (12-hour window). Can be cancelled by creator.
- **active**: At least one other participant has joined (14-day duration starts). Cannot be cancelled.
- **completed**: Challenge auto-completed after 14 days OR creator manually completed via `/completechallenge`. No more scoring allowed. Creator can still award winners if not already done.
- **cancelled**: Challenge manually cancelled by creator while pending OR auto-cancelled after 12 hours with no second participant

### One Challenge Per Group (Active or Pending)

Enforced via application logic:

```go
// Check before creating - only one active OR pending challenge allowed
// Completed challenges don't block creation of new challenges
WHERE group_id = ? AND status IN ('active', 'pending')
```

Notes:

- All challenges are stored permanently and accessible via `/viewchallenge [id]`
- `/viewchallenge` (no args) shows the current challenge (most recent pending/active/completed)
- `/listchallenges` shows paginated history of all challenges
- Multiple completed challenges can exist; only active/pending block new creation

### Wager Tracking

- Each participant's wager stored in `challenge_participants.wager_amount`
- Allows flexibility for different wager amounts per participant
- Creator is automatically a participant with their own wager

### Winner Determination

- **Creator decides winners** after 14 days using `/completechallenge @user1 @user2...`
- Participants can view final scores using `/viewchallenge [id]` to help creator determine winners
- Creators can access completed challenges anytime via `/viewchallenge [id]` or `/listchallenges`
- Even if new challenges are created, old challenges remain accessible by ID
- Flexible award system: winner-take-all, multiple winners, or everyone who met goals
- Each winner receives: `wager + (wager × multiplier)`
- `is_winner` field in `challenge_participants` set to `TRUE` for winners, `FALSE` for losers
- `NULL` until challenge completes

### Scoring System (Optional/Informational)

- `score` field in `challenge_participants` for tracking progress
- **Only creator can award points** using `/score @user1 [points1] @user2 [points2]...`
- Flexible scoring formats:
  - Individual scores: `/score @john 5 @rob 6` (different points per user)
  - Shared score: `/score @john @rob @mark 10` (same points for all)
- Points can be positive or negative
- Scores are informational only - creator determines winners independently
- Visible via `/viewchallenge` command
- Useful for tracking progress during challenge

### Relationships

- Challenge belongs to a `Group` and has a `Creator` (User)
- Participants link `Users` to `Challenges` with wager, score, and winner status
- Cascade delete on participants when challenge is deleted
- Creator is also a participant (can win their own challenge)

### Timestamps

Consistent with existing schema:

- `created_at`: When challenge was created (used for 12-hour expiration check)
- `updated_at`: Auto-updated on record modification
- `joined_at`: When participant joined challenge
- `start_date`: When first non-creator participant joins (activates challenge, 14-day countdown)

## Implementation Plan

### Phase 1: Database Layer

#### 1. Create Model Files (`internal/database/models/`)

**`challenge.go`**:

- Challenge struct with difficulty, multiplier, status fields
- Status values: pending, active, completed, cancelled
- Difficulty values: easy, moderate, hard, difficult
- Multiplier values: 0.5, 1.0, 1.5, 2.0

**`challenge_participant.go`**:

- ChallengeParticipant struct with wager, score, is_winner fields
- Score field for optional progress tracking (honor system)

#### 2. Create Store Files (`internal/store/`)

**`challenge.go`** - ChallengeStore methods:

- `CreateChallenge()` - Create new challenge with difficulty
- `GetCurrentChallenge()` - Get most recent challenge for group (any status)
- `GetActiveOrPendingChallenge()` - Check if active/pending challenge exists (for creation validation)
- `GetChallengeByID()` - Retrieve specific challenge with participants by ID
- `GetChallengesByGroup()` - Get paginated list of challenges for group, sorted by created_at DESC
- `CountChallengesByGroup()` - Get total count of challenges for pagination
- `UpdateChallengeStatus()` - Change challenge status
- `ActivateChallenge()` - Mark as active and set start_date
- `CancelChallenge()` - Mark as cancelled (manual or 12-hour timeout)
- `CompleteChallenge()` - Mark as completed
- `GetPendingChallengesForCancellation()` - Find pending >12 hours
- `GetActiveChallengesForCompletion()` - Find active >14 days

**`challenge_participant.go`** - ParticipantStore methods:

- `CreateParticipant()` - Add participant with wager
- `GetParticipants()` - Get all participants with scores
- `GetParticipantCount()` - Count participants
- `UpdateScore()` - Award/subtract points for a single user (only creator can call)
- `UpdateScores()` - Award/subtract points for multiple users with individual values (only creator can call)
  - Takes map of userID -> points to add/subtract
  - Supports different point values per user
- `GetWinnersByUserIDs()` - Get participants by user IDs
- `SetWinners()` - Mark winners after completion
- `IsUserParticipant()` - Check if user already joined

### Phase 2: Bot Logic

#### 3. Command Handlers (`internal/bot/challenge_handler.go`)

**`/challenge [difficulty] [wager] [title] [description]`**:

- Parse difficulty (easy/moderate/hard/difficult)
- Validate no existing active/pending challenge (completed challenges don't block creation)
- Validate creator has sufficient tokens
- Deduct wager from creator
- Create challenge with multiplier
- Auto-create participant entry for creator
- Send confirmation with 12-hour join countdown

**`/joinchallenge [wager]`**:

- Check for pending or active challenge
- Validate user has sufficient tokens
- Validate user not already in challenge
- Deduct wager from user
- Add participant
- If first non-creator join: activate challenge, set start_date
- Send confirmation with 14-day countdown (if activated), or confirmation of join if already active

**`/viewchallenge [id]`**:

- **No argument**: Display current challenge details (most recent pending/active/completed)
- **With ID**: Display specific challenge by ID (e.g., `/viewchallenge 5`)
- Show challenge status, difficulty, and multiplier
- List all participants with wagers and current scores
- Show potential winnings (wager × multiplier) per participant
- Display time remaining (12h to join or 14d to complete) if not completed
- For completed challenges: show final scores and whether winners have been awarded
- Allows participants to check progress during challenge and review historical challenges

**`/listchallenges [page]`**:

- Display paginated list of challenges for the group, sorted by created_at DESC
- Show challenge ID, title, status, difficulty, and created date
- Default page size: 10 challenges per page
- Include Telegram inline keyboard buttons for pagination:
  - `[« Previous] [Next »]` buttons
  - `[View #ID]` button for each challenge in the list
- Page argument optional (defaults to page 1)
- Clicking `[View #ID]` button triggers `/viewchallenge [id]`

**`/score @user1 [points1] @user2 [points2]...`**:

- Validate only creator can call this
- Parse multiple mentioned users with individual or shared point values
- Validate challenge is active (not completed - no scoring after 14 days)
- Two supported formats:
  - **Individual scores**: `/score @john 5 @rob 6` - john gets 5, rob gets 6
  - **Shared score**: `/score @john @rob @mark 10` - all three get 10 points
  - **Mixed**: Can combine both patterns in parsing logic
- Points can be positive or negative
- Update all mentioned participants' scores accordingly
- Send confirmation with list of updated users and their point changes

**`/cancelchallenge`**:

- Validate only creator can call this
- Validate challenge is pending (not active or completed)
- Mark challenge as cancelled
- Refund creator's wager
- Send cancellation confirmation

**`/completechallenge [@user1 @user2...]`**:

- Validate only creator can call
- Validate challenge is active or completed (allow awarding after auto-completion)
- Parse mentioned users (winners)
- Calculate payouts: `wager + (wager × multiplier)` per winner
- Mint bonus tokens: `(wager × multiplier)` per winner
- Award tokens to winners
- Mark winners in database
- Mark challenge as completed (if not already)
- Send results message with final scores and payouts

**`/help [topic]`** (update existing handler in `handler.go`):

- **No argument**: Show general bot help (existing behavior)
- **`/help challenge`**: Show challenge-specific help (`MsgChallengeHelp`) with:
  - Overview of challenge system
  - List of all challenge commands with examples
  - Difficulty tier explanations (easy/moderate/hard/difficult with multipliers)
  - How scoring works (including flexible individual/shared scoring)
  - How to view and list challenges
  - Rules and time constraints

#### 4. Formatting Functions (`internal/bot/challenge_formatting.go`)

- `formatChallengeDetails()` - Challenge info with status, difficulty badge, time remaining
- `formatParticipantList()` - Participants with wagers, scores, potential winnings
- `formatChallengeResult()` - Completion results with winners and payouts
- `formatScoreboard()` - Live scoreboard with current scores (works for active and completed)
- `formatTimeRemaining()` - Countdown for join window or completion (only for pending/active)
- `formatDifficulty()` - Difficulty badge with emoji
- `formatAwardStatus()` - Show whether winners have been awarded (for completed challenges)
- `formatChallengeList()` - Format list of challenges with ID, title, status, difficulty, date
- `formatChallengeHelp()` - Format challenge-specific help message with commands and examples
- `createPaginationKeyboard()` - Create Telegram inline keyboard with Previous/Next buttons and View buttons

#### 5. Telegram Inline Keyboards (`internal/bot/challenge_keyboard.go`)

**Pagination keyboard structure:**

```go
// Example for page 2 of 5:
InlineKeyboard: [][]InlineKeyboardButton{
    { // Row 1: Navigation
        {Text: "« Previous", CallbackData: "challenge_list_1"},
        {Text: "Next »", CallbackData: "challenge_list_3"},
    },
    { // Row 2: View buttons (challenges on current page)
        {Text: "View #15", CallbackData: "challenge_view_15"},
        {Text: "View #14", CallbackData: "challenge_view_14"},
    },
    // ... more rows for additional challenges
}
```

**Callback handlers:**

- `challenge_list_{page}` - Handle pagination
- `challenge_view_{id}` - Handle viewing specific challenge

### Phase 3: Scheduler & Automation

#### 6. Challenge Scheduler (`internal/bot/challenge_scheduler.go`)

**Background tasks (runs hourly)**:

- **Auto-cancellation**: Find pending challenges >12 hours with only 1 participant
  - Mark as cancelled
  - Refund creator's wager
  - Post cancellation message
- **Auto-completion**: Find active challenges >14 days from start_date
  - Mark as completed (no automatic token distribution)
  - No scoring allowed after this point
  - Post completion message reminding creator to use `/completechallenge` to award winners
  - Creator must manually award based on participant wagers

### Phase 4: Business Logic

#### 7. Validation Rules

- Creator has sufficient tokens before creating
- User has sufficient tokens before joining
- One active OR pending challenge per group (completed challenges don't block new ones)
- No duplicate joins
- Wager amounts > 0
- Valid difficulty: easy, moderate, hard, difficult
- Only creator can award scores (/score) - only while active (before 14 days)
- Score command supports flexible formats:
  - Individual scores: `/score @user1 [points1] @user2 [points2]`
  - Shared score: `/score @user1 @user2... [points]`
- Only creator can cancel challenge (/cancelchallenge) - only while pending (before anyone joins)
- Only creator can complete challenge (/completechallenge) - allowed on active or completed status
- Anyone can view challenge (/viewchallenge [id]) - shows current or specific challenge by ID
- Anyone can list challenges (/listchallenges [page]) - paginated history with inline keyboard navigation
- Use created_at for 12-hour window and auto-cancellation
- Use start_date for 14-day duration and auto-completion

#### 8. Token Integration

**Challenge Creation**:

- Deduct wager from creator
- Auto-create participant for creator

**Joining Challenge**:

- Deduct wager from participant
- If first non-creator: activate and set start_date

**Challenge Completion** (by creator via `/completechallenge`):

- Can be called on active or completed challenges
- For each winner: `payout = wager + (wager × multiplier)`
- Mint `(wager × multiplier)` tokens per winner
- Award tokens to winners
- Non-winners: wagers remain burned
- Mark as completed (if not already)
- If no winners mentioned: all wagers burned

**Challenge Cancellation** (manual by creator OR auto after 12 hours):

- Refund creator's wager
- Mark as cancelled

**Challenge Auto-Completion** (auto, 14 days):

- Mark as completed
- No automatic token distribution
- Creator uses `/completechallenge` to award winners based on their wagers
- No more scoring allowed after auto-completion

**Multiplier Token Economics**:

- Easy (0.5x): Mint 0.5 × wager per winner
- Moderate (1.0x): Mint 1.0 × wager per winner
- Hard (1.5x): Mint 1.5 × wager per winner
- Difficult (2.0x): Mint 2.0 × wager per winner

#### 9. Scoring System

- Only creator can use `/score` command
- Flexible scoring formats:
  - **Individual scores**: `/score @john 5 @rob 6` - john gets 5, rob gets 6
  - **Shared score**: `/score @john @rob @mark 10` - all three get 10 points
- Parsing logic:
  - Parse mentions and numbers in sequence
  - If number follows mention immediately: assign to that user
  - If mentions are followed by a single number at end: assign to all mentioned users
- Points can be positive or negative
- Purely informational - doesn't affect winner determination
- Visible via `/viewchallenge`
- Creator-managed tracking of progress

### Phase 5: Testing & Documentation

#### 10. Tests

**`challenge_test.go`**:

- Challenge creation with difficulties
- Multiplier validation
- Status transitions
- Cancellation/completion queries
- Pagination queries (offset/limit)
- Challenge retrieval by ID

**`challenge_participant_test.go`**:

- Participant creation
- Score updates (single user)
- Batch score updates (multiple users with same points)
- Individual score updates (multiple users with different points)
- Winner selection
- Duplicate prevention

**`challenge_handler_test.go`**:

- Command parsing (including optional ID argument for /viewchallenge)
- Scoring command parsing:
  - Individual scores: `/score @user1 [points1] @user2 [points2]`
  - Shared scores: `/score @user1 @user2 [points]`
- Validation logic
- Payout calculations
- Token minting
- Pagination logic
- Inline keyboard callback handling
- /help challenge command

**`challenge_scheduler_test.go`**:

- 12-hour auto-cancellation detection
- 14-day auto-completion detection
- Refund logic on cancellation
- Auto-completion logic

#### 11. Constants (`internal/constants/messages.go`)

```go
// Challenge messages
const (
    MsgChallengeCreated            = "Challenge created! 12 hours to get participants..."
    MsgChallengeJoined             = "You've joined the challenge!"
    MsgChallengeActivated          = "Challenge activated! 14 days to compete."
    MsgChallengeCompleted          = "Challenge completed!"
    MsgChallengeAutoCompleted      = "Challenge completed after 14 days! Creator: use /completechallenge to award winners."
    MsgChallengeCancelled          = "Challenge cancelled."
    MsgChallengeAutoCancelled      = "Challenge auto-cancelled - no one joined within 12 hours."
    MsgChallengeCannotCancel       = "Cannot cancel - challenge is already active."
    MsgChallengeAlreadyExists      = "A challenge is already active in this group."
    MsgChallengeNotFound           = "No challenge found in this group."
    MsgChallengeInsufficientTokens = "Insufficient tokens for this wager."
    MsgChallengeAlreadyParticipant = "You're already in this challenge."
    MsgChallengeInvalidDifficulty  = "Invalid difficulty. Use: easy, moderate, hard, or difficult"
    MsgChallengeOnlyCreator        = "Only the challenge creator can do that."
    MsgChallengeNoLongerActive     = "Challenge is completed. No more scoring allowed."
    MsgChallengeScoreUpdated       = "Score updated!"
    MsgChallengeScoresUpdated      = "Scores updated for %d participants!"
    MsgChallengeInvalidWager       = "Wager must be greater than 0."
)

// Challenge help text
const MsgChallengeHelp = `🏆 <b>Challenge System</b>

<b>Overview:</b>
Create fitness challenges with token rewards based on difficulty!

<b>Commands:</b>
• <code>/challenge [difficulty] [wager] [title] [description]</code>
  Create a new challenge (easy/moderate/hard/difficult)

• <code>/joinchallenge [wager]</code>
  Join the current challenge

• <code>/viewchallenge [id]</code>
  View current or specific challenge

• <code>/listchallenges [page]</code>
  Browse all challenges with pagination

• <code>/score @user1 [points1] @user2 [points2]...</code>
  Award points to participants (creator only)
  Examples:
    /score @john 5 @rob 6 (individual scores)
    /score @john @rob @mark 10 (shared score)

• <code>/cancelchallenge</code>
  Cancel pending challenge (creator only)

• <code>/completechallenge [@winner1 @winner2...]</code>
  Complete challenge and award winners (creator only)

<b>Difficulty Multipliers:</b>
• Easy (0.5x): Winners get wager + 50%
• Moderate (1.0x): Winners get wager + 100%
• Hard (1.5x): Winners get wager + 150%
• Difficult (2.0x): Winners get wager + 200%

<b>Rules:</b>
• 12-hour window to get participants
• 14-day challenge duration
• Creator determines winners
• Scores are informational only`
)

// Time constants
const (
    ChallengeJoinWindow = 12 * time.Hour
    ChallengeDuration   = 14 * 24 * time.Hour
)

// Difficulty multipliers
const (
    MultiplierEasy        = 0.5
    MultiplierModerate = 1.0
    MultiplierHard        = 1.5
    MultiplierDifficult   = 2.0
)

// Pagination
const (
    ChallengePageSize = 10  // Challenges per page in /listchallenges
)
```

#### 12. Documentation Updates

**README.md** - Add challenge commands section:

- `/challenge [difficulty] [wager] [title] [description]` - Create challenge
- `/joinchallenge [wager]` - Join challenge
- `/viewchallenge [id]` - View current or specific challenge
- `/listchallenges [page]` - List all challenges with pagination
- `/score @user1 [points1] @user2 [points2]...` - Award points (creator only)
  - Supports individual scores: `/score @john 5 @rob 6`
  - Supports shared scores: `/score @john @rob @mark 10`
- `/cancelchallenge` - Cancel pending challenge (creator only)
- `/completechallenge [@winners...]` - Complete and award (creator only)
- `/help challenge` - Show challenge-specific help

**CLAUDE.md** - Add challenge system to features:

- Self-governed challenge system with difficulty tiers
- Inline keyboard pagination for challenge history
- Persistent challenge storage and retrieval by ID

**CHANGELOG.md** - Document challenge system addition
