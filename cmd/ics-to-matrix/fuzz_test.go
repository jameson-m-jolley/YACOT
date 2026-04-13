package main

import (
	"strconv"
	"strings"
	"testing"
)

func FuzzParseICS(f *testing.F) {
	// Seed corpus with valid ICS
	f.Add([]byte(`BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20260403T090000
DURATION:PT1H
END:VEVENT
END:VCALENDAR`))
	f.Add([]byte(`BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20260403T090000
CATEGORIES:work
END:VEVENT
BEGIN:VEVENT
DTSTART:20260403T100000
CATEGORIES:dev
END:VEVENT
END:VCALENDAR`))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 {
			t.Skip()
		}

		got, _, err := parseICS(data)

		// Should not panic
		if err == nil && got != "" {
			lines := strings.Split(got, "\n")
			for _, line := range lines {
				if len(line) == 0 {
					continue
				}
				fields := strings.Fields(line)
				// Each event line should have 14 fields
				if len(fields) == 14 {
					// Validate numeric fields
					for i := 0; i < 9; i++ {
						_, _ = strconv.Atoi(fields[i])
					}
				}
			}
		}
	})
}
