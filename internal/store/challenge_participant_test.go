package store

import (
	"testing"

	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupParticipantTestDB(t *testing.T) *gorm.DB {
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

func createTestChallenge(db *gorm.DB) *models.Challenge {
	challenge := &models.Challenge{
		GroupID:    123,
		CreatorID:  456,
		Title:      "Test Challenge",
		Difficulty: "moderate",
		Multiplier: 1.0,
		Status:     "active",
	}
	db.Create(challenge)
	return challenge
}

func TestCreateParticipant(t *testing.T) {
	db := setupParticipantTestDB(t)
	store := NewParticipantStore(db)

	challenge := createTestChallenge(db)

	participant := &models.ChallengeParticipant{
		ChallengeID: challenge.ID,
		UserID:      789,
		WagerAmount: 100,
	}

	err := store.CreateParticipant(participant)
	if err != nil {
		t.Fatalf("Failed to create participant: %v", err)
	}

	if participant.ID == 0 {
		t.Error("Expected participant ID to be set")
	}

	// Verify record was created
	var found models.ChallengeParticipant
	err = db.First(&found, participant.ID).Error
	if err != nil {
		t.Fatalf("Failed to find created participant: %v", err)
	}

	if found.WagerAmount != 100 {
		t.Errorf("Expected wager 100 but got %d", found.WagerAmount)
	}
}

func TestGetParticipantCount(t *testing.T) {
	db := setupParticipantTestDB(t)
	store := NewParticipantStore(db)

	challenge := createTestChallenge(db)

	// Add 3 participants
	for i := 0; i < 3; i++ {
		db.Create(&models.ChallengeParticipant{
			ChallengeID: challenge.ID,
			UserID:      int64(i + 1),
			WagerAmount: 50,
		})
	}

	count, err := store.GetParticipantCount(challenge.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3 but got %d", count)
	}
}

func TestIsUserParticipant(t *testing.T) {
	db := setupParticipantTestDB(t)
	store := NewParticipantStore(db)

	challenge := createTestChallenge(db)

	// Add a participant
	db.Create(&models.ChallengeParticipant{
		ChallengeID: challenge.ID,
		UserID:      789,
		WagerAmount: 50,
	})

	// User 789 should be a participant
	isParticipant, err := store.IsUserParticipant(challenge.ID, 789)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !isParticipant {
		t.Error("Expected user 789 to be a participant")
	}

	// User 999 should not be a participant
	isParticipant, err = store.IsUserParticipant(challenge.ID, 999)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if isParticipant {
		t.Error("Expected user 999 to not be a participant")
	}
}

func TestUpdateScores(t *testing.T) {
	db := setupParticipantTestDB(t)
	store := NewParticipantStore(db)

	challenge := createTestChallenge(db)

	// Add participants
	p1 := &models.ChallengeParticipant{ChallengeID: challenge.ID, UserID: 1, WagerAmount: 50, Score: 0}
	p2 := &models.ChallengeParticipant{ChallengeID: challenge.ID, UserID: 2, WagerAmount: 50, Score: 5}
	db.Create(p1)
	db.Create(p2)

	// Update scores
	scores := map[int64]int{
		1: 10,
		2: -3,
	}

	err := store.UpdateScores(challenge.ID, scores)
	if err != nil {
		t.Fatalf("Failed to update scores: %v", err)
	}

	// Verify scores
	var found1, found2 models.ChallengeParticipant
	db.First(&found1, p1.ID)
	db.First(&found2, p2.ID)

	if found1.Score != 10 {
		t.Errorf("Expected user 1 score 10 but got %d", found1.Score)
	}

	if found2.Score != 2 { // 5 + (-3) = 2
		t.Errorf("Expected user 2 score 2 but got %d", found2.Score)
	}
}

func TestSetWinners(t *testing.T) {
	db := setupParticipantTestDB(t)
	store := NewParticipantStore(db)

	challenge := createTestChallenge(db)

	// Add 3 participants
	p1 := &models.ChallengeParticipant{ChallengeID: challenge.ID, UserID: 1, WagerAmount: 50}
	p2 := &models.ChallengeParticipant{ChallengeID: challenge.ID, UserID: 2, WagerAmount: 50}
	p3 := &models.ChallengeParticipant{ChallengeID: challenge.ID, UserID: 3, WagerAmount: 50}
	db.Create(p1)
	db.Create(p2)
	db.Create(p3)

	// Set users 1 and 2 as winners
	err := store.SetWinners(challenge.ID, []int64{1, 2})
	if err != nil {
		t.Fatalf("Failed to set winners: %v", err)
	}

	// Verify
	var found1, found2, found3 models.ChallengeParticipant
	db.First(&found1, p1.ID)
	db.First(&found2, p2.ID)
	db.First(&found3, p3.ID)

	if found1.IsWinner == nil || !*found1.IsWinner {
		t.Error("Expected user 1 to be a winner")
	}

	if found2.IsWinner == nil || !*found2.IsWinner {
		t.Error("Expected user 2 to be a winner")
	}

	if found3.IsWinner == nil || *found3.IsWinner {
		t.Error("Expected user 3 to be a loser")
	}
}

func TestSetWinners_NoWinners(t *testing.T) {
	db := setupParticipantTestDB(t)
	store := NewParticipantStore(db)

	challenge := createTestChallenge(db)

	// Add 2 participants
	p1 := &models.ChallengeParticipant{ChallengeID: challenge.ID, UserID: 1, WagerAmount: 50}
	p2 := &models.ChallengeParticipant{ChallengeID: challenge.ID, UserID: 2, WagerAmount: 50}
	db.Create(p1)
	db.Create(p2)

	// Set no winners
	err := store.SetWinners(challenge.ID, []int64{})
	if err != nil {
		t.Fatalf("Failed to set no winners: %v", err)
	}

	// Verify all are losers
	var found1, found2 models.ChallengeParticipant
	db.First(&found1, p1.ID)
	db.First(&found2, p2.ID)

	if found1.IsWinner == nil || *found1.IsWinner {
		t.Error("Expected user 1 to be a loser")
	}

	if found2.IsWinner == nil || *found2.IsWinner {
		t.Error("Expected user 2 to be a loser")
	}
}

func TestGetParticipantsByUserIDs(t *testing.T) {
	db := setupParticipantTestDB(t)
	store := NewParticipantStore(db)

	challenge := createTestChallenge(db)

	// Add 3 participants
	db.Create(&models.ChallengeParticipant{ChallengeID: challenge.ID, UserID: 1, WagerAmount: 50})
	db.Create(&models.ChallengeParticipant{ChallengeID: challenge.ID, UserID: 2, WagerAmount: 75})
	db.Create(&models.ChallengeParticipant{ChallengeID: challenge.ID, UserID: 3, WagerAmount: 100})

	// Get users 1 and 3
	participants, err := store.GetParticipantsByUserIDs(challenge.ID, []int64{1, 3})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(participants) != 2 {
		t.Fatalf("Expected 2 participants but got %d", len(participants))
	}

	// Verify we got the right ones
	found1, found3 := false, false
	for _, p := range participants {
		if p.UserID == 1 && p.WagerAmount == 50 {
			found1 = true
		}
		if p.UserID == 3 && p.WagerAmount == 100 {
			found3 = true
		}
	}

	if !found1 {
		t.Error("Expected to find user 1")
	}
	if !found3 {
		t.Error("Expected to find user 3")
	}
}
