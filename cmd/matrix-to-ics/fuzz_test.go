package main

import (
	"strings"
	"testing"
)

func FuzzParseMatrix(f *testing.F) {
	// Seed corpus with valid flat single-line matrix (30 fields per event)
	f.Add("0 2026 4 3 9 0 60 0 0 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 1")
	f.Add("0 2026 4 3 9 0 60 0 0 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 1 1 2026 4 3 10 0 120 0 0 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 1")

	f.Fuzz(func(t *testing.T, data string) {
		// parseMatrix expects flat single-line input - normalize whitespace
		normalized := strings.Join(strings.Fields(data), " ")

		// Should not panic - parseMatrix should handle gracefully
		got, err := parseMatrix(normalized, nil)

		// If there's no error, validate output structure
		if err == nil && len(got) > 0 {
			if !strings.Contains(got, "BEGIN:VCALENDAR") {
				t.Error("missing VCALENDAR header")
			}
			if !strings.Contains(got, "END:VCALENDAR") {
				t.Error("missing VCALENDAR footer")
			}

			eventCount := strings.Count(got, "BEGIN:VEVENT")
			endCount := strings.Count(got, "END:VEVENT")
			if eventCount != endCount {
				t.Error("mismatched VEVENT BEGIN/END")
			}
		}
	})
}
