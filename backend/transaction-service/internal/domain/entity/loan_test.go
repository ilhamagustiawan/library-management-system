package entity

import (
	"testing"
	"time"
)

func TestQuoteFineChargesEachStartedOverdueDay(t *testing.T) {
	dueAt := time.Date(2026, 7, 26, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		quotedAt time.Time
		wantDays int
	}{
		{name: "on time", quotedAt: dueAt, wantDays: 0},
		{name: "one second late", quotedAt: dueAt.Add(time.Second), wantDays: 1},
		{name: "exactly one day late", quotedAt: dueAt.Add(24 * time.Hour), wantDays: 1},
		{name: "second overdue day started", quotedAt: dueAt.Add(24*time.Hour + time.Second), wantDays: 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			quote := QuoteFine(dueAt, test.quotedAt, 5000)
			if test.wantDays == 0 {
				if quote != nil {
					t.Fatalf("QuoteFine() = %#v, want no fine", quote)
				}
				return
			}
			if quote == nil || quote.OverdueDays != test.wantDays || quote.TotalAmountMinor != int64(test.wantDays*5000) {
				t.Fatalf("QuoteFine() = %#v, want %d days", quote, test.wantDays)
			}
		})
	}
}
