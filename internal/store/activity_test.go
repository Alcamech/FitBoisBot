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

func TestCreateRecord(t *testing.T) {
	db := setupTestDB(t)
	store := NewActivityStore(db)

	post := ActivityPost{
		UserID:    123,
		GroupID:   456,
		MessageID: 789,
		Activity:  "running",
		Month:     "03",
		Day:       "15",
		Year:      "2024",
	}

	err := store.CreateRecord(post)
	if err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}

	// Verify record was created
	var activity models.Activity
	err = db.Where("user_id = ? AND group_id = ? AND message_id = ?", post.UserID, post.GroupID, post.MessageID).First(&activity).Error
	if err != nil {
		t.Fatalf("Failed to find created record: %v", err)
	}

	if activity.Activity != post.Activity {
		t.Errorf("Expected activity %q but got %q", post.Activity, activity.Activity)
	}
	if activity.Month != post.Month {
		t.Errorf("Expected month %q but got %q", post.Month, activity.Month)
	}
	if activity.Day != post.Day {
		t.Errorf("Expected day %q but got %q", post.Day, activity.Day)
	}
	if activity.Year != post.Year {
		t.Errorf("Expected year %q but got %q", post.Year, activity.Year)
	}
}

func TestUpdateRecord_Success(t *testing.T) {
	db := setupTestDB(t)
	store := NewActivityStore(db)

	// Create initial record
	initial := models.Activity{
		UserID:    123,
		GroupID:   456,
		MessageID: 789,
		Activity:  "running",
		Month:     "03",
		Day:       "15",
		Year:      "2024",
	}
	if err := db.Create(&initial).Error; err != nil {
		t.Fatalf("Failed to create initial record: %v", err)
	}

	// Update with new data
	post := ActivityPost{
		UserID:    123,
		GroupID:   456,
		MessageID: 789,
		Activity:  "cycling", // Changed
		Month:     "03",
		Day:       "16", // Changed
		Year:      "2024",
	}

	err := store.UpdateRecord(post)
	if err != nil {
		t.Fatalf("Failed to update record: %v", err)
	}

	// Verify record was updated
	var activity models.Activity
	err = db.Where("user_id = ? AND group_id = ? AND message_id = ?", post.UserID, post.GroupID, post.MessageID).First(&activity).Error
	if err != nil {
		t.Fatalf("Failed to find updated record: %v", err)
	}

	if activity.Activity != "cycling" {
		t.Errorf("Expected activity %q but got %q", "cycling", activity.Activity)
	}
	if activity.Day != "16" {
		t.Errorf("Expected day %q but got %q", "16", activity.Day)
	}
}

func TestUpdateRecord_NotFound(t *testing.T) {
	db := setupTestDB(t)
	store := NewActivityStore(db)

	post := ActivityPost{
		UserID:    123,
		GroupID:   456,
		MessageID: 999, // Doesn't exist
		Activity:  "running",
		Month:     "03",
		Day:       "15",
		Year:      "2024",
	}

	err := store.UpdateRecord(post)
	if err == nil {
		t.Error("Expected error when updating non-existent record, but got none")
	}
}

func TestGetCountByUserMonthYear(t *testing.T) {
	db := setupTestDB(t)
	store := NewActivityStore(db)

	userID := int64(123)
	groupID := int64(456)
	month := "03"
	year := "2024"

	// Create 3 activities for the user in March 2024
	activities := []models.Activity{
		{UserID: userID, GroupID: groupID, Month: month, Year: year, Activity: "running", Day: "01", MessageID: 1},
		{UserID: userID, GroupID: groupID, Month: month, Year: year, Activity: "cycling", Day: "02", MessageID: 2},
		{UserID: userID, GroupID: groupID, Month: month, Year: year, Activity: "swimming", Day: "03", MessageID: 3},
		// Different month - should not count
		{UserID: userID, GroupID: groupID, Month: "04", Year: year, Activity: "walking", Day: "01", MessageID: 4},
		// Different year - should not count
		{UserID: userID, GroupID: groupID, Month: month, Year: "2023", Activity: "hiking", Day: "01", MessageID: 5},
	}

	for _, activity := range activities {
		if err := db.Create(&activity).Error; err != nil {
			t.Fatalf("Failed to create test activity: %v", err)
		}
	}

	count, err := store.GetCountByUserMonthYear(userID, groupID, month, year)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3 but got %d", count)
	}
}

func TestGetCountByUserMonthYear_NoActivities(t *testing.T) {
	db := setupTestDB(t)
	store := NewActivityStore(db)

	count, err := store.GetCountByUserMonthYear(999, 888, "03", "2024")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0 but got %d", count)
	}
}
