package bot

import (
	"testing"

	"github.com/Alcamech/FitBoisBot/internal/constants"
)

func TestParseChallengeCommand_Valid(t *testing.T) {
	tests := []struct {
		name        string
		args        string
		wantDiff    string
		wantMult    float64
		wantWager   int
		wantTitle   string
		wantDesc    string
	}{
		{
			name:      "easy with description",
			args:      "easy 100 MyChallenge This is a description",
			wantDiff:  "easy",
			wantMult:  0.5,
			wantWager: 100,
			wantTitle: "MyChallenge",
			wantDesc:  "This is a description",
		},
		{
			name:      "moderate no description",
			args:      "moderate 50 FitnessChallenge",
			wantDiff:  "moderate",
			wantMult:  1.0,
			wantWager: 50,
			wantTitle: "FitnessChallenge",
			wantDesc:  "",
		},
		{
			name:      "hard uppercase",
			args:      "HARD 200 HardChallenge",
			wantDiff:  "hard",
			wantMult:  1.5,
			wantWager: 200,
			wantTitle: "HardChallenge",
			wantDesc:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := parseChallengeCommand(tt.args)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if params.Difficulty != tt.wantDiff {
				t.Errorf("Difficulty = %q, want %q", params.Difficulty, tt.wantDiff)
			}
			if params.Multiplier != tt.wantMult {
				t.Errorf("Multiplier = %f, want %f", params.Multiplier, tt.wantMult)
			}
			if params.Wager != tt.wantWager {
				t.Errorf("Wager = %d, want %d", params.Wager, tt.wantWager)
			}
			if params.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", params.Title, tt.wantTitle)
			}
			if params.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", params.Description, tt.wantDesc)
			}
		})
	}
}

func TestParseChallengeCommand_Invalid(t *testing.T) {
	tests := []struct {
		name string
		args string
	}{
		{name: "empty args", args: ""},
		{name: "too few args", args: "easy 100"},
		{name: "invalid difficulty", args: "invalid 100 Title"},
		{name: "invalid wager", args: "easy abc Title"},
		{name: "zero wager", args: "easy 0 Title"},
		{name: "negative wager", args: "easy -50 Title"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseChallengeCommand(tt.args)
			if err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}

func TestGetDifficultyMultiplier(t *testing.T) {
	tests := []struct {
		difficulty string
		want       float64
		wantErr    bool
	}{
		{"easy", constants.MultiplierEasy, false},
		{"moderate", constants.MultiplierModerate, false},
		{"hard", constants.MultiplierHard, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.difficulty, func(t *testing.T) {
			got, err := getDifficultyMultiplier(tt.difficulty)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDifficultyMultiplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getDifficultyMultiplier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJoinChallengeCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    int
		wantErr bool
	}{
		{"valid wager", "100", 100, false},
		{"valid with spaces", "  50  ", 50, false},
		{"empty", "", 0, true},
		{"zero", "0", 0, true},
		{"negative", "-10", 0, true},
		{"non-numeric", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseJoinChallengeCommand(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJoinChallengeCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseJoinChallengeCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseViewChallengeCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantID  int64
		wantHas bool
	}{
		{"empty - show current", "", 0, false},
		{"valid ID", "42", 42, true},
		{"ID with spaces", "  123  ", 123, true},
		{"invalid ID", "abc", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotHas := parseViewChallengeCommand(tt.args)
			if gotID != tt.wantID {
				t.Errorf("parseViewChallengeCommand() ID = %v, want %v", gotID, tt.wantID)
			}
			if gotHas != tt.wantHas {
				t.Errorf("parseViewChallengeCommand() hasID = %v, want %v", gotHas, tt.wantHas)
			}
		})
	}
}

func TestParseListChallengesCommand(t *testing.T) {
	tests := []struct {
		name string
		args string
		want int
	}{
		{"empty - default to 1", "", 1},
		{"valid page", "2", 2},
		{"page with spaces", "  3  ", 3},
		{"invalid page", "abc", 1},
		{"zero page", "0", 1},
		{"negative page", "-1", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseListChallengesCommand(tt.args)
			if got != tt.want {
				t.Errorf("parseListChallengesCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCallbackData(t *testing.T) {
	tests := []struct {
		data       string
		wantAction string
		wantValue  string
	}{
		{"challenge_list_1", "challenge_list", "1"},
		{"challenge_view_42", "challenge_view", "42"},
		{"challenge_join_5", "challenge_join", "5"},
		{"invalid", "", ""},
		{"single", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			action, value := parseCallbackData(tt.data)
			if action != tt.wantAction {
				t.Errorf("parseCallbackData() action = %q, want %q", action, tt.wantAction)
			}
			if value != tt.wantValue {
				t.Errorf("parseCallbackData() value = %q, want %q", value, tt.wantValue)
			}
		})
	}
}
