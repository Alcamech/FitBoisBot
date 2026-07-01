package bot

import "testing"

func TestCanScoreEntries(t *testing.T) {
	const creatorID = int64(100)
	const memberID = int64(200)
	const otherID = int64(300)

	tests := []struct {
		name      string
		isCreator bool
		senderID  int64
		entries   []ScoreEntry
		want      bool
	}{
		{
			name:      "creator scores others",
			isCreator: true,
			senderID:  creatorID,
			entries:   []ScoreEntry{{UserID: memberID, Points: 5}, {UserID: otherID, Points: 6}},
			want:      true,
		},
		{
			name:      "creator scores self",
			isCreator: true,
			senderID:  creatorID,
			entries:   []ScoreEntry{{UserID: creatorID, Points: 5}},
			want:      true,
		},
		{
			name:      "member scores self",
			isCreator: false,
			senderID:  memberID,
			entries:   []ScoreEntry{{UserID: memberID, Points: 5}},
			want:      true,
		},
		{
			name:      "member scores someone else",
			isCreator: false,
			senderID:  memberID,
			entries:   []ScoreEntry{{UserID: otherID, Points: 5}},
			want:      false,
		},
		{
			name:      "member scores self and someone else",
			isCreator: false,
			senderID:  memberID,
			entries:   []ScoreEntry{{UserID: memberID, Points: 5}, {UserID: otherID, Points: 6}},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canScoreEntries(tt.isCreator, tt.senderID, tt.entries); got != tt.want {
				t.Errorf("canScoreEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}
