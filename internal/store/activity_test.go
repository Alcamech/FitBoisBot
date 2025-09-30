package store

import (
	"testing"

	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Activity{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestGetMostActiveUsersForMonth_SingleWinner(t *testing.T) {
	db := setupTestDB(t)
	store := NewActivityStore(db)

	groupID := int64(123)
	month := "03"
	year := "2024"

	// User 1: 3 activities, User 2: 2 activities
	activities := []models.Activity{
		{UserID: 1, GroupID: groupID, Month: month, Year: year, Activity: "running", Day: "01", MessageID: 1},
		{UserID: 1, GroupID: groupID, Month: month, Year: year, Activity: "cycling", Day: "02", MessageID: 2},
		{UserID: 1, GroupID: groupID, Month: month, Year: year, Activity: "swimming", Day: "03", MessageID: 3},
		{UserID: 2, GroupID: groupID, Month: month, Year: year, Activity: "walking", Day: "01", MessageID: 4},
		{UserID: 2, GroupID: groupID, Month: month, Year: year, Activity: "hiking", Day: "02", MessageID: 5},
	}

	for _, activity := range activities {
		if err := db.Create(&activity).Error; err != nil {
			t.Fatalf("Failed to create test activity: %v", err)
		}
	}

	userIDs, count, err := store.GetMostActiveUsersForMonth(groupID, month, year)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3 but got %d", count)
	}

	if len(userIDs) != 1 {
		t.Errorf("Expected 1 user but got %d", len(userIDs))
	}

	if len(userIDs) > 0 && userIDs[0] != 1 {
		t.Errorf("Expected user ID 1 but got %d", userIDs[0])
	}
}

func TestGetMostActiveUsersForMonth_TwoWayTie(t *testing.T) {
	db := setupTestDB(t)
	store := NewActivityStore(db)

	groupID := int64(123)
	month := "03"
	year := "2024"

	// User 1 and 2: 2 activities each, User 3: 1 activity
	activities := []models.Activity{
		{UserID: 1, GroupID: groupID, Month: month, Year: year, Activity: "running", Day: "01", MessageID: 1},
		{UserID: 1, GroupID: groupID, Month: month, Year: year, Activity: "cycling", Day: "02", MessageID: 2},
		{UserID: 2, GroupID: groupID, Month: month, Year: year, Activity: "walking", Day: "01", MessageID: 3},
		{UserID: 2, GroupID: groupID, Month: month, Year: year, Activity: "hiking", Day: "02", MessageID: 4},
		{UserID: 3, GroupID: groupID, Month: month, Year: year, Activity: "swimming", Day: "01", MessageID: 5},
	}

	for _, activity := range activities {
		if err := db.Create(&activity).Error; err != nil {
			t.Fatalf("Failed to create test activity: %v", err)
		}
	}

	userIDs, count, err := store.GetMostActiveUsersForMonth(groupID, month, year)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected count 2 but got %d", count)
	}

	if len(userIDs) != 2 {
		t.Fatalf("Expected 2 users but got %d", len(userIDs))
	}

	// Check both users are present (order doesn't matter)
	foundUser1 := false
	foundUser2 := false
	for _, id := range userIDs {
		if id == 1 {
			foundUser1 = true
		}
		if id == 2 {
			foundUser2 = true
		}
	}

	if !foundUser1 || !foundUser2 {
		t.Errorf("Expected users 1 and 2, got %v", userIDs)
	}
}

func TestGetMostActiveUsersForMonth_ThreeWayTie(t *testing.T) {
	db := setupTestDB(t)
	store := NewActivityStore(db)

	groupID := int64(123)
	month := "03"
	year := "2024"

	// All three users: 1 activity each
	activities := []models.Activity{
		{UserID: 1, GroupID: groupID, Month: month, Year: year, Activity: "running", Day: "01", MessageID: 1},
		{UserID: 2, GroupID: groupID, Month: month, Year: year, Activity: "walking", Day: "01", MessageID: 2},
		{UserID: 3, GroupID: groupID, Month: month, Year: year, Activity: "swimming", Day: "01", MessageID: 3},
	}

	for _, activity := range activities {
		if err := db.Create(&activity).Error; err != nil {
			t.Fatalf("Failed to create test activity: %v", err)
		}
	}

	userIDs, count, err := store.GetMostActiveUsersForMonth(groupID, month, year)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count 1 but got %d", count)
	}

	if len(userIDs) != 3 {
		t.Fatalf("Expected 3 users but got %d", len(userIDs))
	}
}

func TestGetMostActiveUsersForMonth_NoActivities(t *testing.T) {
	db := setupTestDB(t)
	store := NewActivityStore(db)

	groupID := int64(123)
	month := "03"
	year := "2024"

	userIDs, count, err := store.GetMostActiveUsersForMonth(groupID, month, year)
	
	if err == nil {
		t.Error("Expected error but got none")
	}

	if count != 0 {
		t.Errorf("Expected count 0 but got %d", count)
	}

	if userIDs != nil {
		t.Errorf("Expected nil userIDs but got %v", userIDs)
	}
}
