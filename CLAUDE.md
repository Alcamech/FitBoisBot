# CLAUDE.md

FitBoisBot is a Go Telegram bot for fitness accountability and gamification.

## Commands

### Development

- `make dev` - Run with live reload
- `make build` - Build the application  
- `make test` - Run tests
- `make db-up` - Start database container
- `make db-seed` - Insert mock test data
- `make db-login` - Connect to database
- `make help` - Show all available commands

### CLI Usage

- `./bin/fitboisbot serve` - Start bot server
- `./bin/fitboisbot announce "message"` - Send to all groups
- `./bin/fitboisbot announce --file msg.md` - Send from file
- `./bin/fitboisbot announce --group -123456789 "text"` - Send to specific group
- `./bin/fitboisbot announce --preview "text"` - Preview without sending
- `./bin/fitboisbot version` - Show version information

### Environment Variables

- `TOKEN` - Telegram bot token (required)
- `DEBUG` - Debug mode (optional, defaults to false)  
- `DB_*` - Database connection settings

## Architecture

### Structure

- `cmd/fitboisbot/` - Application entry point with Cobra CLI
- `internal/cli/` - CLI commands (serve, announce, version)
- `internal/bot/` - Bot logic (handlers, parsing, formatting, validation, time utilities)
- `internal/store/` - Data access layer with GORM
- `internal/constants/` - Configuration constants and messages
- `internal/database/models/` - Database entities
- `internal/version/` - Version tracking and build info
- `config/` - Viper configuration
- `scripts/` - Database scripts and migrations

### Key Files

- `bot.go` - Main bot service and message processing
- `handler.go` - Command handlers (/help, /fastgg, /tokens, /balance, /timezone)
- `challenge_handler.go` - Challenge system command handlers
- `challenge_formatting.go` - Challenge display formatting
- `challenge_keyboard.go` - Inline keyboard pagination for challenges
- `challenge_parsing.go` - Challenge command parsing
- `challenge_scheduler.go` - Auto-cancel and auto-complete scheduler
- `parsing.go` - Activity message parsing
- `validation.go` - Input validation
- `time.go` - Timezone operations
- `formatting.go` - HTML display formatting for Telegram
- `message.go` - Message sending utilities
- `token_scheduler.go` - Monthly reward scheduler

### Features

- Activity tracking (activity-MM-DD-YYYY or activity-MM-DD-YY format)
- Self-correction support for editing activity posts
- Fastest GG competition system (prevents self-GGing)
- Monthly FitBoi Awards with automatic tie-handling and token distribution
- Token reward system with leaderboards and individual balance lookup
- Challenge system with difficulty tiers (easy/moderate/hard), token wagers, and multiplier payouts
- Inline keyboard pagination for challenge browsing
- Auto-cancel (12h) and auto-complete (14d) challenge scheduling
- Multi-group timezone support
- HTML-formatted sorted leaderboards (activity counts, GG rankings, token balances)
- Monthly activity leaderboard formatting
- Admin CLI for announcements
- Version tracking with build-time injection
- Full timestamp tracking (created_at/updated_at) across all entities
- Comprehensive test coverage (formatting, keyboard, parsing, store operations)
- Command shorthands for faster interaction (/c, /jc, /vc, /lc, /s, /cc, /done, /b, /t, /l, /gg, /tz, /h)
