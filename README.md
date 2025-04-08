# FitBois Telegram Bot

Stay accountable. Get fit. Dominate the leaderboard.

[**Launch the Bot**](https://t.me/FitBoisBot)

---

## 📜 Current Commands

Use `/help` in Telegram to see these:

- `/help` – Show available commands and usage.
- `/fastgg` – Display the Fast GG leaderboard for your group.
- `/tokens` – Display Fitboi Token balances.
- `/timezone` – Set your group’s timezone (e.g., `/timezone America/New_York`).

### 🗓️ Usage

Post activities with an image of your tracked workout and a caption using these formats:

```
activity-MM-DD-YYYY
or
activity-MM-DD-YY
```

Be the first to reply `gg` to your friends

---

## 🛠️ Roadmap & TODOs

### General Improvements

- [ ] Improve and update documentation
- [x] Update dependencies
- [x] Small refactor
- [x] TODOs in code
- [x] Break out methods
- [x] Improve formatting of bot messages
- [x] Error handling
- [x] Allow users to set timezone
- [x] Scale for other groups
- [x] Fit Boi of the Year Awards EOY
- [ ] Write test
- [ ] Create changelog
- [ ] Create release notes
- [ ] Improved logging
- [ ] Monthly Check-Ins
- [ ] Monthly rewards
- [ ] Intelligent Q/A with SQL data

---

## 🧠 Planned Features

### `/stats` – Personal Stats Summary

Returns an individual’s activity summary:

- Total activity count
- Count by month
- Count by activity type
- Fastest GG count

### `/challenge` – Challenge Completion (Draft)

Submit a challenge attempt:

```
/challenge 1 100
```

- `1` = challenge complete (optional)
- `100` = token wager amount
