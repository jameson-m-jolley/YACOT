package main

import (
	"strings"
	"testing"
)

func TestParseMatrix(t *testing.T) {
	empty := "-1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 1 -1"

	tests := []struct {
		name      string
		input     string
		expectErr bool
		check     func(t *testing.T, got string)
	}{
		{
			name:  "single event",
			input: "0 2026 4 3 9 0 60 0 0 " + empty,
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, "BEGIN:VCALENDAR") {
					t.Error("missing VCALENDAR header")
				}
				if !strings.Contains(got, "BEGIN:VEVENT") {
					t.Error("missing VEVENT")
				}
				if !strings.Contains(got, "DTSTART:20260403T090000") {
					t.Error("missing DTSTART")
				}
				if !strings.Contains(got, "DURATION:PT1H") {
					t.Error("missing DURATION")
				}
				if !strings.Contains(got, "END:VCALENDAR") {
					t.Error("missing VCALENDAR footer")
				}
				if !strings.Contains(got, "UID:event-0@yacot") {
					t.Error("missing UID")
				}
			},
		},
		{
			name:  "multiple events",
			input: "0 2026 4 3 9 0 60 0 0 " + empty + " 1 2026 4 3 10 0 120 0 0 " + empty,
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, "UID:event-0@yacot") {
					t.Error("missing event-0")
				}
				if !strings.Contains(got, "UID:event-1@yacot") {
					t.Error("missing event-1")
				}
				if !strings.Contains(got, "DTSTART:20260403T090000") {
					t.Error("missing first DTSTART")
				}
				if !strings.Contains(got, "DTSTART:20260403T100000") {
					t.Error("missing second DTSTART")
				}
				if !strings.Contains(got, "DURATION:PT2H") {
					t.Error("missing PT2H")
				}
			},
		},
		{
			name:  "duration with minutes only",
			input: "0 2026 4 3 9 0 30 0 0 " + empty,
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, "DURATION:PT30M") {
					t.Error("missing PT30M")
				}
			},
		},
		{
			name:  "duration with hours and minutes",
			input: "0 2026 4 3 9 0 90 0 0 " + empty,
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, "DURATION:PT1H30M") {
					t.Error("missing PT1H30M")
				}
			},
		},
		{
			name:  "empty input",
			input: "",
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, "BEGIN:VCALENDAR") {
					t.Error("missing VCALENDAR header")
				}
				if !strings.Contains(got, "END:VCALENDAR") {
					t.Error("missing VCALENDAR footer")
				}
			},
		},
		{
			name:      "invalid field count",
			input:     "0 2026 4 3 9 0 60 0",
			expectErr: true,
		},
		{
			name:  "category dev",
			input: "0 2026 4 3 9 0 60 0 1 " + empty,
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, "CATEGORIES:dev") {
					t.Error("missing dev category")
				}
			},
		},
		{
			name:  "category personal",
			input: "0 2026 4 3 9 0 60 0 2 " + empty,
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, "CATEGORIES:personal") {
					t.Error("missing personal category")
				}
			},
		},
		{
			name:  "priority is set",
			input: "0 2026 4 3 9 0 60 5 0 " + empty,
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, "PRIORITY:5") {
					t.Error("missing priority")
				}
			},
		},
		{
			name:  "with recurrence",
			input: "0 2026 4 3 9 0 60 0 0 -1 -1 -1 -1 -1 4 1 2026 12 31 23 59 -1 -1 -1 -1 -1 -1 -1 1 -1",
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, "RRULE:FREQ=WEEKLY") {
					t.Error("missing RRULE")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMatrix(tt.input, nil)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("parseMatrix() error = %v", err)
				return
			}
			tt.check(t, got)
		})
	}
}

func TestParseMatrixOutput(t *testing.T) {
	input := "0 2026 4 3 9 0 60 0 0 " + "-1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 1 -1"
	got, err := parseMatrix(input, nil)
	if err != nil {
		t.Fatalf("parseMatrix() error = %v", err)
	}

	eventCount := strings.Count(got, "BEGIN:VEVENT")
	if eventCount != 1 {
		t.Errorf("expected 1 event, got %d", eventCount)
	}
}
