# Changelog

All notable changes to FitBoisBot will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.1.0] - [Unreleased]

### Added

- **CLI Commands**: Complete Cobra CLI with `serve`, `announce`, and `version` subcommands
- **Admin Announcements**: Send messages to all groups or specific groups via CLI  
- **Version Tracking**: Build-time version injection with git commit and build time
- **Database Seeding**: `make db-seed` command to insert mock test data
- **HTML Formatting**: Improved message formatting using HTML entities instead of Markdown
- **Timestamp Tracking**: Added created_at/updated_at columns to all database tables
- **Token Balance Updates**: Proper timestamp tracking for token balance changes

### Fixed

- **Self-GG Prevention**: Users can no longer GG their own activity posts
- **Self Correction**: Users can edit their activity posts to correct the entry
- **Message Formatting**: Switched from MarkdownV2 to HTML to fix Telegram parsing errors
- **Leaderboard Display**: Clean ranked list format for activity counts, GG rankings, and token balances

### Changed

- **Message Layout**: Simplified leaderboard formatting with bold numbers and clean structure
- **Documentation**: Updated README.md, CLAUDE.md with current features and CLI usage
- **Development Workflow**: Added database seeding for easier testing

### Technical

- Implemented transaction-safe database migrations
- Added mock data script with realistic test scenarios
- Enhanced error handling and logging
- Updated Makefile with new commands
- Updated all Go models with timestamp fields
- Enhanced production migration script with timestamp additions
- Updated mock data to include timestamps

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

