# FitBoisBot

Telegram bot for fitness accountability with activity tracking, GG competition, and token rewards.

## Features

- **Activity Tracking**: Post `activity-MM-DD-YYYY` format activities
- **GG Competition**: First to reply `gg` earns fastest GG points  
- **Token Rewards**: Monthly rewards for active users
- **Multi-timezone**: Group-specific timezone support
- **Leaderboards**: Activity counts, GG rankings, token balances

## Commands

- `/help` - Show bot info and version
- `/fastgg` - Display fastest GG leaderboard  
- `/tokens` - Show token leaderboard
- `/timezone [zone]` - View/set group timezone

## Quick Start

```bash
cp .env.example .env    # Add your TOKEN
make db-up             # Start database
make db-seed           # Add test data (optional)
make dev               # Run bot
```

## CLI Commands

```bash
make build             # Build binary
./bin/fitboisbot serve # Start bot server

# Admin announcements
./bin/fitboisbot announce "Message"
./bin/fitboisbot announce --file msg.md
./bin/fitboisbot announce --group -123 "Group message"

./bin/fitboisbot version # Show version info
```

## Make Commands

```bash
make help      # Show all commands
make dev       # Development mode
make test      # Run tests  
make db-up     # Start database
make db-seed   # Insert mock data
make db-login  # Database console
```
