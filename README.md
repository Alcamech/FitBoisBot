# FitBoisBot

Telegram bot for fitness accountability with activity tracking, GG competition, and token rewards.

## Features

- **Activity Tracking**: Post `activity-MM-DD-YYYY` or `activity-MM-DD-YY` format activities
- **GG Competition**: First to reply `gg` earns fastest GG points
- **Monthly FitBoi Awards**: Monthly rewards for most active users with automatic tie-handling
- **Token System**: Earn and track FitBoi Tokens with leaderboards
- **Challenge System**: Create fitness challenges with token wagers and difficulty multipliers
- **Multi-timezone**: Group-specific timezone support
- **HTML Leaderboards**: Clean, sorted activity counts, GG rankings, and token balances
- **Admin Tools**: CLI for announcements and version management
- **Self-Correction**: Edit your own activity posts to fix mistakes

## Commands

| Command | Shorthand | Description |
|---------|-----------|-------------|
| `/help` | `/h` | Show bot info and version |
| `/fastgg` | `/gg` | Display fastest GG leaderboard |
| `/tokens` | `/t` | Show token leaderboard |
| `/balance` | `/b` | Show your token balance |
| `/leaderboard` | `/l` | Display monthly activity leaderboard |
| `/timezone` | `/tz` | View/set group timezone |
| `/challenge` | `/c` | Create a fitness challenge |
| `/joinchallenge` | `/jc` | Join current challenge with wager |
| `/viewchallenge` | `/vc` | View challenge details |
| `/listchallenges` | `/lc` | Browse all challenges |
| `/score` | `/s` | Award points to participants |
| `/cancelchallenge` | `/cc` | Cancel pending challenge |
| `/completechallenge` | `/done` | Complete challenge and award winners |

Use `/help challenge` or `/help c` for detailed challenge system documentation.

## Quick Start

```bash
cp .env.example .env    # Add your TOKEN
make db/up             # Start database
make db/seed           # Add test data (optional)
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
