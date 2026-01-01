package bot

import (
	"testing"
	"time"

	"github.com/Alcamech/FitBoisBot/internal/constants"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"days and hours", 3*24*time.Hour + 5*time.Hour, "3d 5h"},
		{"just days", 7 * 24 * time.Hour, "7d 0h"},
		{"hours and minutes", 5*time.Hour + 30*time.Minute, "5h 30m"},
		{"just hours", 2 * time.Hour, "2h 0m"},
		{"just minutes", 45 * time.Minute, "45m"},
		{"zero", 0, "0m"},
		{"14 days", 14 * 24 * time.Hour, "14d 0h"},
		{"12 hours", 12 * time.Hour, "12h 0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.duration)
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
			}
		})
	}
}

func TestFormatStatusBadge(t *testing.T) {
	tests := []struct {
		status string
		want   string
	}{
		{constants.ChallengeStatusPending, constants.StatusBadgePending},
		{constants.ChallengeStatusActive, constants.StatusBadgeActive},
		{constants.ChallengeStatusCompleted, constants.StatusBadgeCompleted},
		{constants.ChallengeStatusCancelled, constants.StatusBadgeCancelled},
		{"unknown", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := formatStatusBadge(tt.status)
			if got != tt.want {
				t.Errorf("formatStatusBadge(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestFormatDifficultyBadge(t *testing.T) {
	tests := []struct {
		difficulty string
		want       string
	}{
		{constants.DifficultyEasy, constants.DifficultyBadgeEasy},
		{constants.DifficultyModerate, constants.DifficultyBadgeModerate},
		{constants.DifficultyHard, constants.DifficultyBadgeHard},
		{"unknown", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.difficulty, func(t *testing.T) {
			got := formatDifficultyBadge(tt.difficulty)
			if got != tt.want {
				t.Errorf("formatDifficultyBadge(%q) = %q, want %q", tt.difficulty, got, tt.want)
			}
		})
	}
}

func TestFormatScoreUpdate(t *testing.T) {
	tests := []struct {
		name    string
		entries []ScoreEntry
		names   map[int64]string
		wantContains []string
	}{
		{
			name: "single positive score",
			entries: []ScoreEntry{
				{UserID: 1, Points: 10},
			},
			names: map[int64]string{1: "Alice"},
			wantContains: []string{"Scores Updated", "Alice", "+10"},
		},
		{
			name: "single negative score",
			entries: []ScoreEntry{
				{UserID: 1, Points: -5},
			},
			names: map[int64]string{1: "Bob"},
			wantContains: []string{"Scores Updated", "Bob", "-5"},
		},
		{
			name: "multiple scores",
			entries: []ScoreEntry{
				{UserID: 1, Points: 10},
				{UserID: 2, Points: 5},
			},
			names: map[int64]string{1: "Alice", 2: "Bob"},
			wantContains: []string{"Alice", "+10", "Bob", "+5"},
		},
		{
			name: "unknown user",
			entries: []ScoreEntry{
				{UserID: 999, Points: 3},
			},
			names: map[int64]string{},
			wantContains: []string{"Unknown", "+3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatScoreUpdate(tt.entries, tt.names)
			for _, want := range tt.wantContains {
				if !contains(got, want) {
					t.Errorf("formatScoreUpdate() = %q, want to contain %q", got, want)
				}
			}
		})
	}
}

// helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchString(s, substr)))
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
