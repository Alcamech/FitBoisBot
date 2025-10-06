package bot

import (
	"strings"
	"testing"
)

func TestFormatMonthlyAwardMessage_SingleWinner(t *testing.T) {
	winnerNames := []string{"Alice"}
	month := "03"
	year := "2024"
	activityCount := int64(15)

	message := formatMonthlyAwardMessage(winnerNames, month, year, activityCount)

	if !strings.Contains(message, "Alice") {
		t.Errorf("Expected message to contain winner name 'Alice', got: %s", message)
	}
	if !strings.Contains(message, "03") {
		t.Errorf("Expected message to contain month '03', got: %s", message)
	}
	if !strings.Contains(message, "2024") {
		t.Errorf("Expected message to contain year '2024', got: %s", message)
	}
	if !strings.Contains(message, "100") {
		t.Errorf("Expected message to contain full reward '100', got: %s", message)
	}
}

func TestFormatMonthlyAwardMessage_TwoWayTie(t *testing.T) {
	winnerNames := []string{"Alice", "Bob"}
	month := "03"
	year := "2024"
	activityCount := int64(10)

	message := formatMonthlyAwardMessage(winnerNames, month, year, activityCount)

	if !strings.Contains(message, "Alice") {
		t.Errorf("Expected message to contain 'Alice', got: %s", message)
	}
	if !strings.Contains(message, "Bob") {
		t.Errorf("Expected message to contain 'Bob', got: %s", message)
	}
	if !strings.Contains(message, "and") {
		t.Errorf("Expected message to contain 'and' for two winners, got: %s", message)
	}
	// Total reward should be 100, split 50 each
	if !strings.Contains(message, "50") {
		t.Errorf("Expected message to contain per-user reward '50', got: %s", message)
	}
}

func TestFormatMonthlyAwardMessage_ThreeWayTie(t *testing.T) {
	winnerNames := []string{"Alice", "Bob", "Charlie"}
	month := "03"
	year := "2024"
	activityCount := int64(8)

	message := formatMonthlyAwardMessage(winnerNames, month, year, activityCount)

	if !strings.Contains(message, "Alice") {
		t.Errorf("Expected message to contain 'Alice', got: %s", message)
	}
	if !strings.Contains(message, "Bob") {
		t.Errorf("Expected message to contain 'Bob', got: %s", message)
	}
	if !strings.Contains(message, "Charlie") {
		t.Errorf("Expected message to contain 'Charlie', got: %s", message)
	}
	// 100 tokens / 3 winners = 33 tokens each
	if !strings.Contains(message, "33") {
		t.Errorf("Expected message to contain per-user reward '33', got: %s", message)
	}
}

func TestFormatMonthlyAwardMessage_FourWayTie(t *testing.T) {
	winnerNames := []string{"Alice", "Bob", "Charlie", "Diana"}
	month := "04"
	year := "2025"
	activityCount := int64(5)

	message := formatMonthlyAwardMessage(winnerNames, month, year, activityCount)

	// Verify all names present
	for _, name := range winnerNames {
		if !strings.Contains(message, name) {
			t.Errorf("Expected message to contain '%s', got: %s", name, message)
		}
	}

	// 100 tokens / 4 winners = 25 tokens each
	if !strings.Contains(message, "25") {
		t.Errorf("Expected message to contain per-user reward '25', got: %s", message)
	}
}

func TestFormatActivityCounts_SingleUser(t *testing.T) {
	userCounts := map[string]int64{
		"Alice": 5,
	}

	message := formatActivityCounts(userCounts)

	if !strings.Contains(message, "Alice") {
		t.Errorf("Expected message to contain 'Alice', got: %s", message)
	}
	if !strings.Contains(message, "5") {
		t.Errorf("Expected message to contain count '5', got: %s", message)
	}
}

func TestFormatActivityCounts_MultipleUsers_Sorted(t *testing.T) {
	userCounts := map[string]int64{
		"Alice":   3,
		"Bob":     10,
		"Charlie": 7,
	}

	message := formatActivityCounts(userCounts)

	// Verify all users are present
	if !strings.Contains(message, "Alice") {
		t.Errorf("Expected message to contain 'Alice', got: %s", message)
	}
	if !strings.Contains(message, "Bob") {
		t.Errorf("Expected message to contain 'Bob', got: %s", message)
	}
	if !strings.Contains(message, "Charlie") {
		t.Errorf("Expected message to contain 'Charlie', got: %s", message)
	}

	// Verify Bob (highest count) appears before Charlie and Alice
	bobIndex := strings.Index(message, "Bob")
	charlieIndex := strings.Index(message, "Charlie")
	aliceIndex := strings.Index(message, "Alice")

	if bobIndex == -1 || charlieIndex == -1 || aliceIndex == -1 {
		t.Fatalf("Not all users found in message")
	}

	if bobIndex > charlieIndex {
		t.Errorf("Expected Bob (10) to appear before Charlie (7)")
	}
	if charlieIndex > aliceIndex {
		t.Errorf("Expected Charlie (7) to appear before Alice (3)")
	}
}

func TestFormatActivityCounts_EmptyMap(t *testing.T) {
	userCounts := map[string]int64{}

	message := formatActivityCounts(userCounts)

	// Should return a message indicating no activities
	if message == "" {
		t.Error("Expected non-empty message for empty user counts")
	}
	// Message should indicate no activities recorded
	if !strings.Contains(message, "No activities") && !strings.Contains(message, "no activities") {
		t.Logf("Message returned: %s", message)
	}
}
