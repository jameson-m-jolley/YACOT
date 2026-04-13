package main

import (
	"strings"
	"testing"
)

func TestParseICS(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "basic event",
			input: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20260403T090000
DURATION:PT1H
PRIORITY:5
CATEGORIES:work
END:VEVENT
END:VCALENDAR`,
			wantErr: false,
		},
		{
			name: "multiple events",
			input: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20260403T090000
DURATION:PT1H
END:VEVENT
BEGIN:VEVENT
DTSTART:20260403T100000
DURATION:PT2H
END:VEVENT
END:VCALENDAR`,
			wantErr: false,
		},
		{
			name: "dev category",
			input: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20260403T090000
CATEGORIES:dev
END:VEVENT
END:VCALENDAR`,
			wantErr: false,
		},
		{
			name: "personal category",
			input: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20260403T090000
CATEGORIES:personal
END:VEVENT
END:VCALENDAR`,
			wantErr: false,
		},
		{
			name: "duration with hours and minutes",
			input: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20260403T090000
DURATION:PT1H30M
END:VEVENT
END:VCALENDAR`,
			wantErr: false,
		},
		{
			name: "duration with only minutes",
			input: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20260403T090000
DURATION:PT30M
END:VEVENT
END:VCALENDAR`,
			wantErr: false,
		},
		{
			name: "empty calendar",
			input: `BEGIN:VCALENDAR
VERSION:2.0
END:VCALENDAR`,
			wantErr: false,
		},
		{
			name:    "invalid ics",
			input:   `not ics`,
			wantErr: true,
		},
		{
			name: "event without DTSTART uses defaults",
			input: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
END:VEVENT
END:VCALENDAR`,
			wantErr: false,
		},
		{
			name: "duration only hours",
			input: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20260403T090000
DURATION:PT2H
END:VEVENT
END:VCALENDAR`,
			wantErr: false,
		},
		{
			name:    "event with CRLF line endings",
			input:   "BEGIN:VCALENDAR\r\nVERSION:2.0\r\nBEGIN:VEVENT\r\nDTSTART:20260403T090000\r\nEND:VEVENT\r\nEND:VCALENDAR",
			wantErr: false,
		},
		{
			name: "event without duration uses default 60",
			input: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20260403T090000
END:VEVENT
END:VCALENDAR`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := parseICS([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseICS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			lines := strings.Split(strings.TrimSpace(got), "\n")
			expectedLines := 0
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					expectedLines++
				}
			}
			if expectedLines == 0 && strings.TrimSpace(got) != "" {
				t.Error("expected lines but got none")
			}
		})
	}
}

func TestParseICSOutputFormat(t *testing.T) {
	input := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:20260403T090000
DURATION:PT1H
END:VEVENT
END:VCALENDAR`

	got, _, err := parseICS([]byte(input))
	if err != nil {
		t.Fatalf("parseICS() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	fields := strings.Fields(lines[0])
	if len(fields) != 30 {
		t.Errorf("expected 30 fields, got %d", len(fields))
	}

	if fields[0] != "0" {
		t.Errorf("first field should be 0, got %s", fields[0])
	}
	if fields[1] != "2026" {
		t.Errorf("year should be 2026, got %s", fields[1])
	}
	if fields[2] != "4" {
		t.Errorf("month should be 4, got %s", fields[2])
	}
	if fields[3] != "3" {
		t.Errorf("day should be 3, got %s", fields[3])
	}
	if fields[4] != "9" {
		t.Errorf("hour should be 9, got %s", fields[4])
	}
	if fields[5] != "0" {
		t.Errorf("minute should be 0, got %s", fields[5])
	}
}

func TestParseICSFuzz(t *testing.T) {
	fuzzCases := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"whitespace only", "   \n\t  "},
		{"random bytes", "BEGIN:VCALENDAR\n\x00\x01\x02"},
		{"very long line", "BEGIN:VCALENDAR\n" + strings.Repeat("X", 10000) + "\nEND:VCALENDAR"},
		{"just BEGIN", "BEGIN"},
		{"just END", "END"},
		{"truncated DTSTART", "BEGIN:VCALENDAR\nBEGIN:VEVENT\nDTSTART:2026\nEND:VEVENT"},
		{"truncated DTSTART 4 chars", "BEGIN:VCALENDAR\nBEGIN:VEVENT\nDTSTART:202\nEND:VEVENT"},
		{"truncated DTSTART 6 chars", "BEGIN:VCALENDAR\nBEGIN:VEVENT\nDTSTART:20260\nEND:VEVENT"},
	}

	for _, tc := range fuzzCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panicked: %v", r)
				}
			}()
			_, _, _ = parseICS([]byte(tc.input))
		})
	}
}
