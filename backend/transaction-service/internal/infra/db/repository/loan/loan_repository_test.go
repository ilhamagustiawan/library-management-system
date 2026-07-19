package loan

import (
	"testing"
	"time"
)

func TestCalculateOverdueDaysChargesStartedDays(t *testing.T) {
	due := time.Date(2026, 7, 26, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		returned time.Time
		want     int
	}{
		{name: "on time", returned: due, want: 0},
		{name: "one second late", returned: due.Add(time.Second), want: 1},
		{name: "exact day", returned: due.Add(24 * time.Hour), want: 1},
		{name: "started second day", returned: due.Add(24*time.Hour + time.Second), want: 2},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := calculateOverdueDays(due, test.returned); got != test.want {
				t.Fatalf("calculateOverdueDays() = %d, want %d", got, test.want)
			}
		})
	}
}
