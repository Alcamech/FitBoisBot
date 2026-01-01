# Changelog

All notable changes to FitBoisBot will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [3.0.0] - 2026-01-01

### Added

- **Challenge System**: Complete fitness challenge system with token wagers
  - Create challenges with difficulty tiers (easy/moderate/hard)
  - Multiplier payouts (0.5x, 1.0x, 1.5x)
  - Creator-managed scoring and winner selection
  - 12-hour join window, 14-day challenge duration
  - Auto-cancel and auto-complete scheduling
  - Inline keyboard pagination for challenge browsing
- **Lifetime Token Balance**: `/balance` command shows total tokens across all years
- **Command Shorthands**: Quick aliases for all commands (/c, /jc, /vc, /lc, /s, /cc, /done, /b, /t, /l, /gg, /tz, /h)

### Changed

- **Token System**: Tokens now tracked per year with lifetime totals for spending
- **Help Text**: Updated with all new commands and shorthands
- **View Challenge**: Improved formatting with labeled Status and Difficulty fields
- **Documentation**: Updated README.md and CLAUDE.md with challenge system details

### Technical

- Added challenge and challenge_participant database models
- Added challenge store and participant store with full CRUD operations
- Implemented callback query handling for inline keyboards
- Added constants for badges, button labels, and status values
- Comprehensive test coverage for formatting, keyboard, and parsing functions

## [2.1.0] - 2025-10-12

### Added

- **CLI Commands**: Complete Cobra CLI with `serve`, `announce`, and `version` subcommands
- **Admin Announcements**: Send messages to all groups or specific groups via CLI
- **Version Tracking**: Build-time version injection with git commit and build time
- **Database Seeding**: `make db-seed` command to insert mock test data
- **HTML Formatting**: Improved message formatting using HTML entities instead of Markdown
- **Timestamp Tracking**: Added created_at/updated_at columns to all database tables
- **Token Balance Updates**: Proper timestamp tracking for token balance changes
- **Monthly FitBoi Awards**: Automatic monthly awards with tie-handling and split token distribution
- **Activity Leaderboard Sorting**: Monthly activity leaderboard now displays sorted by count
- **Test Coverage**: Comprehensive tests for formatting and activity store operations

### Fixed

- **Self-GG Prevention**: Users can no longer GG their own activity posts
- **Self Correction**: Users can edit their activity posts to correct the entry
- **Message Formatting**: Switched from MarkdownV2 to HTML to fix Telegram parsing errors
- **Leaderboard Display**: Clean ranked list format for activity counts, GG rankings, and token balances
- **Activity Parsing**: Improved activity date parsing logic
- **Token Distribution**: Fixed edge cases in monthly reward distribution
- **Message Sending**: Enhanced retry logic and error handling

### Changed

- **Message Layout**: Simplified leaderboard formatting with bold numbers and clean structure
- **Documentation**: Updated README.md, CLAUDE.md, CHANGELOG.md with current features
- **Development Workflow**: Added database seeding for easier testing
- **Mock Data**: Improved test data scenarios for better development testing

### Technical

- Implemented transaction-safe database migrations
- Added mock data script with realistic test scenarios
- Enhanced error handling and logging throughout codebase
- Updated Makefile with new commands
- Updated all Go models with timestamp fields
- Enhanced production migration script with timestamp additions
- Updated mock data to include timestamps
- Added `formatting_test.go` for leaderboard formatting tests
- Enhanced `activity_test.go` with comprehensive store operation tests
- Refactored activity store operations for better testability

## [2.0.0] - 2025-01-07

### Added

- Timezone management with `/timezone` command
- Token reward system with monthly distribution
- `/help` command with comprehensive bot information
- Support for 2-digit year format (YY) in activities
- Token leaderboard with `/tokens` command

### Changed

- Activity format now supports both YY and YYYY years
- Improved message formatting and display

### Fixed

- Timezone handling for different groups
- GG system now properly handles current year activities

## [1.9.0] - 2024-12-31

### Added

- Complete feature set implementation
- Monthly and yearly activity tracking
- Configuration management with environment variables
- Database schema improvements

### Changed

- Migrated to Go from Java/Spring Boot
- Dockerized MariaDB setup
- GORM ORM integration
- Improved project structure

### Fixed

- Database configuration and connection handling
- Schema optimizations

## [1.8.0] - 2024-09-28

### Added

- Activity detection and parsing system
- GG detection and tracking
- User and activity data persistence
- Repository pattern for data access

### Changed

- Major refactoring of bot logic
- Improved message processing pipeline

## [1.7.0] - 2024-06-09

### Added

- Go project structure and foundation
- Database connection management
- Configuration initialization system
- Telegram bot polling mechanism

### Changed

- Complete rewrite from Java to Go
- New project architecture and organization

## [1.0.4] - 2023

### Changed

- Code refactoring and improvements
- Fastest GG sorting functionality

## [1.0.3] - 2023

### Added

- `/fastgg` command for GG leaderboard

## [1.0.2] - 2023

### Added

- "Fastest GG in the west" feature

### Changed

- General refactoring and improvements

## [1.0.0] - 2023

### Added

- Initial bot functionality
- Basic activity tracking
- GG competition system
- Core Telegram bot features

