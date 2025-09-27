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
- `handler.go` - Command handlers (/help, /fastgg, /tokens, /timezone)
- `parsing.go` - Activity message parsing
- `validation.go` - Input validation
- `time.go` - Timezone operations
- `formatting.go` - HTML display formatting for Telegram
- `message.go` - Message sending utilities
- `token_scheduler.go` - Monthly reward scheduler

### Features

- Activity tracking (activity-MM-DD-YYYY format)
- Fastest GG competition system
- Monthly token rewards
- Multi-group timezone support
- HTML-formatted leaderboards (activity counts, GG rankings, token balances)
- Admin CLI for announcements
- Version tracking with build-time injection
- Full timestamp tracking (created_at/updated_at) across all entities
