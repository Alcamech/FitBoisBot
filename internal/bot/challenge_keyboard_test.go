package bot

import (
	"testing"

	"github.com/Alcamech/FitBoisBot/internal/constants"
	"github.com/Alcamech/FitBoisBot/internal/database/models"
)

func TestCreateChallengeListKeyboard_FirstPage(t *testing.T) {
	challenges := []models.Challenge{
		{ID: 10},
		{ID: 9},
		{ID: 8},
	}

	keyboard := createChallengeListKeyboard(challenges, 1, 3)

	// Should have navigation row (only Next since page 1) and view buttons
	if len(keyboard.InlineKeyboard) < 1 {
		t.Fatalf("Expected at least 1 row, got %d", len(keyboard.InlineKeyboard))
	}

	// First row should be navigation with only "Next"
	navRow := keyboard.InlineKeyboard[0]
	if len(navRow) != 1 {
		t.Errorf("Expected 1 nav button on first page, got %d", len(navRow))
	}
	if navRow[0].Text != constants.ButtonNext {
		t.Errorf("Expected %q button, got %q", constants.ButtonNext, navRow[0].Text)
	}
	if navRow[0].CallbackData == nil || *navRow[0].CallbackData != "challenge_list_2" {
		t.Errorf("Expected callback 'challenge_list_2'")
	}
}

func TestCreateChallengeListKeyboard_MiddlePage(t *testing.T) {
	challenges := []models.Challenge{
		{ID: 5},
		{ID: 4},
	}

	keyboard := createChallengeListKeyboard(challenges, 2, 3)

	// Should have both Previous and Next
	navRow := keyboard.InlineKeyboard[0]
	if len(navRow) != 2 {
		t.Errorf("Expected 2 nav buttons on middle page, got %d", len(navRow))
	}

	if navRow[0].Text != constants.ButtonPrevious {
		t.Errorf("Expected %q button, got %q", constants.ButtonPrevious, navRow[0].Text)
	}
	if navRow[0].CallbackData == nil || *navRow[0].CallbackData != "challenge_list_1" {
		t.Errorf("Expected callback 'challenge_list_1'")
	}

	if navRow[1].Text != constants.ButtonNext {
		t.Errorf("Expected %q button, got %q", constants.ButtonNext, navRow[1].Text)
	}
	if navRow[1].CallbackData == nil || *navRow[1].CallbackData != "challenge_list_3" {
		t.Errorf("Expected callback 'challenge_list_3'")
	}
}

func TestCreateChallengeListKeyboard_LastPage(t *testing.T) {
	challenges := []models.Challenge{
		{ID: 2},
		{ID: 1},
	}

	keyboard := createChallengeListKeyboard(challenges, 3, 3)

	// Should have only Previous
	navRow := keyboard.InlineKeyboard[0]
	if len(navRow) != 1 {
		t.Errorf("Expected 1 nav button on last page, got %d", len(navRow))
	}

	if navRow[0].Text != constants.ButtonPrevious {
		t.Errorf("Expected %q button, got %q", constants.ButtonPrevious, navRow[0].Text)
	}
}

func TestCreateChallengeListKeyboard_SinglePage(t *testing.T) {
	challenges := []models.Challenge{
		{ID: 3},
		{ID: 2},
		{ID: 1},
	}

	keyboard := createChallengeListKeyboard(challenges, 1, 1)

	// No navigation row needed, only view buttons
	// First row should be view buttons since no nav needed
	if len(keyboard.InlineKeyboard) == 0 {
		t.Fatal("Expected at least 1 row for view buttons")
	}

	// All buttons should be view buttons (no nav)
	for _, row := range keyboard.InlineKeyboard {
		for _, btn := range row {
			if btn.Text == constants.ButtonPrevious || btn.Text == constants.ButtonNext {
				t.Error("Should not have navigation buttons on single page")
			}
		}
	}
}

func TestCreateChallengeListKeyboard_ViewButtons(t *testing.T) {
	challenges := []models.Challenge{
		{ID: 42},
		{ID: 15},
		{ID: 7},
	}

	keyboard := createChallengeListKeyboard(challenges, 1, 1)

	// Should have view buttons for each challenge
	foundCallbacks := make(map[string]bool)
	foundLabels := make(map[string]bool)
	for _, row := range keyboard.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil {
				foundCallbacks[*btn.CallbackData] = true
			}
			foundLabels[btn.Text] = true
		}
	}

	expectedCallbacks := []string{
		"challenge_view_42",
		"challenge_view_15",
		"challenge_view_7",
	}

	for _, expected := range expectedCallbacks {
		if !foundCallbacks[expected] {
			t.Errorf("Expected callback %q not found", expected)
		}
	}

	// Check labels are numbered 1, 2, 3 (not IDs)
	expectedLabels := []string{"1", "2", "3"}
	for _, expected := range expectedLabels {
		if !foundLabels[expected] {
			t.Errorf("Expected label %q not found", expected)
		}
	}
}

func TestCreateChallengeListKeyboard_Empty(t *testing.T) {
	challenges := []models.Challenge{}

	keyboard := createChallengeListKeyboard(challenges, 1, 1)

	// Should still work with empty list
	// No view buttons, no nav (single page)
	for _, row := range keyboard.InlineKeyboard {
		for _, btn := range row {
			if btn.Text == constants.ButtonPrevious || btn.Text == constants.ButtonNext {
				t.Error("Should not have navigation buttons with empty list")
			}
		}
	}
}
