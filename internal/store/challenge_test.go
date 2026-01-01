package store

import (
	"testing"
	"time"

	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupChallengeTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Challenge{}, &models.ChallengeParticipant{}, &models.User{}, &models.Group{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestCreateChallenge(t *testing.T) {
	db := setupChallengeTestDB(t)
	store := NewChallengeStore(db)

	challenge := &models.Challenge{
		GroupID:     123,
		CreatorID:   456,
		Title:       "30 Day Challenge",
		Description: "Exercise every day",
		Difficulty:  "moderate",
		Multiplier:  1.0,
		Status:      "pending",
	}

	err := store.CreateChallenge(challenge)
	if err != nil {
		t.Fatalf("Failed to create challenge: %v", err)
	}

	if challenge.ID == 0 {
		t.Error("Expected challenge ID to be set")
	}

	// Verify record was created
	var found models.Challenge
	err = db.First(&found, challenge.ID).Error
	if err != nil {
		t.Fatalf("Failed to find created challenge: %v", err)
	}

	if found.Title != "30 Day Challenge" {
		t.Errorf("Expected title %q but got %q", "30 Day Challenge", found.Title)
	}
}

func TestGetActiveOrPendingChallenge(t *testing.T) {
	db := setupChallengeTestDB(t)
	store := NewChallengeStore(db)

	groupID := int64(123)

	// Create a pending challenge
	pending := &models.Challenge{
		GroupID:    groupID,
		CreatorID:  456,
		Title:      "Pending Challenge",
		Difficulty: "easy",
		Multiplier: 0.5,
		Status:     "pending",
	}
	db.Create(pending)

	// Should find the pending challenge
	found, err := store.GetActiveOrPendingChallenge(groupID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if found.ID != pending.ID {
		t.Errorf("Expected challenge ID %d but got %d", pending.ID, found.ID)
	}

	// Update to active
	pending.Status = "active"
	db.Save(pending)

	// Should still find it
	found, err = store.GetActiveOrPendingChallenge(groupID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if found.ID != pending.ID {
		t.Errorf("Expected challenge ID %d but got %d", pending.ID, found.ID)
	}

	// Update to completed
	pending.Status = "completed"
	db.Save(pending)

	// Should not find it
	_, err = store.GetActiveOrPendingChallenge(groupID)
	if err == nil {
		t.Error("Expected error when no active/pending challenge exists")
	}
}

func TestGetChallengesByGroup_Pagination(t *testing.T) {
	db := setupChallengeTestDB(t)
	store := NewChallengeStore(db)

	groupID := int64(123)

	// Create 5 challenges
	for i := 0; i < 5; i++ {
		challenge := &models.Challenge{
			GroupID:    groupID,
			CreatorID:  456,
			Title:      "Challenge",
			Difficulty: "easy",
			Multiplier: 0.5,
			Status:     "completed",
		}
		db.Create(challenge)
		time.Sleep(10 * time.Millisecond) // Ensure different created_at
	}

	// Get first page
	challenges, err := store.GetChallengesByGroup(groupID, 2, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(challenges) != 2 {
		t.Errorf("Expected 2 challenges but got %d", len(challenges))
	}

	// Get second page
	challenges, err = store.GetChallengesByGroup(groupID, 2, 2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(challenges) != 2 {
		t.Errorf("Expected 2 challenges but got %d", len(challenges))
	}

	// Get third page (only 1 left)
	challenges, err = store.GetChallengesByGroup(groupID, 2, 4)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(challenges) != 1 {
		t.Errorf("Expected 1 challenge but got %d", len(challenges))
	}
}

func TestCountChallengesByGroup(t *testing.T) {
	db := setupChallengeTestDB(t)
	store := NewChallengeStore(db)

	groupID := int64(123)

	// Create 3 challenges for target group
	for i := 0; i < 3; i++ {
		db.Create(&models.Challenge{
			GroupID:    groupID,
			CreatorID:  456,
			Title:      "Challenge",
			Difficulty: "easy",
			Multiplier: 0.5,
			Status:     "completed",
		})
	}

	// Create 2 challenges for different group
	for i := 0; i < 2; i++ {
		db.Create(&models.Challenge{
			GroupID:    999,
			CreatorID:  456,
			Title:      "Other Challenge",
			Difficulty: "easy",
			Multiplier: 0.5,
			Status:     "completed",
		})
	}

	count, err := store.CountChallengesByGroup(groupID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3 but got %d", count)
	}
}

func TestActivateChallenge(t *testing.T) {
	db := setupChallengeTestDB(t)
	store := NewChallengeStore(db)

	challenge := &models.Challenge{
		GroupID:    123,
		CreatorID:  456,
		Title:      "Test Challenge",
		Difficulty: "hard",
		Multiplier: 1.5,
		Status:     "pending",
	}
	db.Create(challenge)

	err := store.ActivateChallenge(challenge.ID)
	if err != nil {
		t.Fatalf("Failed to activate challenge: %v", err)
	}

	// Verify status and start_date
	var found models.Challenge
	db.First(&found, challenge.ID)

	if found.Status != "active" {
		t.Errorf("Expected status 'active' but got %q", found.Status)
	}

	if found.StartDate == nil {
		t.Error("Expected start_date to be set")
	}
}

func TestGetPendingChallengesForCancellation(t *testing.T) {
	db := setupChallengeTestDB(t)
	store := NewChallengeStore(db)

	// Create an old pending challenge (should be found)
	old := &models.Challenge{
		GroupID:    123,
		CreatorID:  456,
		Title:      "Old Pending",
		Difficulty: "easy",
		Multiplier: 0.5,
		Status:     "pending",
	}
	db.Create(old)

	// Manually set created_at to 13 hours ago
	db.Model(old).Update("created_at", time.Now().Add(-13*time.Hour))

	// Create a new pending challenge (should not be found)
	newChallenge := &models.Challenge{
		GroupID:    124,
		CreatorID:  456,
		Title:      "New Pending",
		Difficulty: "easy",
		Multiplier: 0.5,
		Status:     "pending",
	}
	db.Create(newChallenge)

	challenges, err := store.GetPendingChallengesForCancellation(12 * time.Hour)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(challenges) != 1 {
		t.Fatalf("Expected 1 challenge but got %d", len(challenges))
	}

	if challenges[0].ID != old.ID {
		t.Errorf("Expected challenge ID %d but got %d", old.ID, challenges[0].ID)
	}
}
