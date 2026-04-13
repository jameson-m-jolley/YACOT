package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	ical "github.com/arran4/golang-ical"
	ics "github.com/arran4/golang-ical"
)

// eventMetadata holds string data that cannot be represented in the numeric matrix.
type eventMetadata struct {
	ID          int      `json:"id"`
	UID         string   `json:"uid"`
	Summary     string   `json:"summary"`
	Description string   `json:"description,omitempty"`
	Location    string   `json:"location,omitempty"`
	Organizer   string   `json:"organizer,omitempty"`
	Attendees   []string `json:"attendees,omitempty"`
}

// metadataFile represents the JSON file containing all event metadata.
type metadataFile struct {
	Source string          `json:"source"`
	Cached string          `json:"cached"`
	Events []eventMetadata `json:"events"`
}

// weekdayToInt converts an RFC 5545 weekday abbreviation to a numeric code.
// Returns 0 for Sunday through 6 for Saturday, -1 for invalid values.
//
// RFC 5545 Weekday Codes:
//   - SU (Sunday)   = 0
//   - MO (Monday)   = 1
//   - TU (Tuesday)  = 2
//   - WE (Wednesday) = 3
//   - TH (Thursday) = 4
//   - FR (Friday)   = 5
//   - SA (Saturday) = 6
func weekdayToInt(wd string) int {
	switch strings.ToUpper(wd) {
	case "SU":
		return 0
	case "MO":
		return 1
	case "TU":
		return 2
	case "WE":
		return 3
	case "TH":
		return 4
	case "FR":
		return 5
	case "SA":
		return 6
	}
	return -1
}

// parseDateTime parses an ICS date-time string (YYYYMMDDTHHMMSS format).
// Returns year, month, day, hour, minute.
// Default values are used if parsing fails.
func parseDateTime(dt string) (year, month, day, hour, minute int) {
	year, month, day, hour, minute = 2026, 1, 1, 0, 0
	if len(dt) >= 4 {
		fmt.Sscanf(dt[:4], "%d", &year)
	}
	if len(dt) >= 8 {
		fmt.Sscanf(dt[4:6], "%d", &month)
		fmt.Sscanf(dt[6:8], "%d", &day)
	}
	if len(dt) >= 13 && dt[8] == 'T' {
		fmt.Sscanf(dt[9:11], "%d", &hour)
		fmt.Sscanf(dt[11:13], "%d", &minute)
	}
	return year, month, day, hour, minute
}

// parseDuration parses an ICS duration string (ISO 8601 duration format).
// Handles PT1H30M format (1 hour, 30 minutes).
// Returns duration in minutes. Default is 60 minutes.
func parseDuration(durStr string) int {
	if !strings.HasPrefix(durStr, "PT") {
		return 60
	}
	durStr = durStr[2:] // Remove "PT" prefix
	hours, minutes := 0, 0

	if idx := strings.Index(durStr, "H"); idx != -1 {
		fmt.Sscanf(durStr[:idx], "%d", &hours)
		durStr = durStr[idx+1:]
	}
	if idx := strings.Index(durStr, "M"); idx != -1 {
		fmt.Sscanf(durStr[:idx], "%d", &minutes)
	}

	return hours*60 + minutes
}

// parseCategory converts an ICS categories string to a numeric code.
//   - "work"       = 0
//   - "dev"        = 1
//   - "personal"  = 2
//   - Other       = 0 (default to work)
func parseCategory(catStr string) int {
	catStr = strings.ToLower(catStr)
	switch {
	case strings.Contains(catStr, "dev"):
		return 1
	case strings.Contains(catStr, "personal"):
		return 2
	default:
		return 0
	}
}

// eventFields holds all parsed fields for an event in the flat byte stream format.
// 30 fields total as defined in AGENTS.md Flat Byte Stream Format.
type eventFields struct {
	id             int
	year           int
	month          int
	day            int
	hour           int
	minute         int
	duration       int
	priority       int
	category       int
	deadlineYear   int
	deadlineMonth  int
	deadlineDay    int
	deadlineHour   int
	deadlineMinute int
	frequency      int
	interval       int
	untilYear      int
	untilMonth     int
	untilDay       int
	untilHour      int
	untilMinute    int
	byDayMask      int
	byMonthMask    int
	byMonthDayMask int
	byHourMask     int
	byMinuteMask   int
	bySecondMask   int
	count          int
	isMutable      int
	wkst           int
}

// parseRRULE parses a recurrence rule and returns all recurrence-related fields.
// Handles all RRULE components defined in RFC 5545 Section 3.3.10.
func parseRRULE(rules []*ical.RecurrenceRule) eventFields {
	// Default values for events without recurrence
	fields := eventFields{
		frequency:      -1,
		interval:       1,
		untilYear:      -1,
		untilMonth:     -1,
		untilDay:       -1,
		untilHour:      -1,
		untilMinute:    -1,
		byDayMask:      -1,
		byMonthMask:    -1,
		byMonthDayMask: -1,
		byHourMask:     -1,
		byMinuteMask:   -1,
		bySecondMask:   -1,
		count:          -1,
		isMutable:      1,
		wkst:           -1,
	}

	// No recurrence rule = return defaults
	if len(rules) == 0 {
		return fields
	}

	r := rules[0]

	// Parse frequency (FREQ field)
	switch r.Freq {
	case ical.FrequencyMinutely:
		fields.frequency = 0
	case ical.FrequencyHourly:
		fields.frequency = 2
	case ical.FrequencyDaily:
		fields.frequency = 3
	case ical.FrequencyWeekly:
		fields.frequency = 4
	case ical.FrequencyMonthly:
		fields.frequency = 5
	case ical.FrequencyYearly:
		fields.frequency = 6
	}

	// Parse interval (INTERVAL field)
	if r.Interval > 0 {
		fields.interval = r.Interval
	}

	// Parse count (COUNT field)
	if r.Count > 0 {
		fields.count = r.Count
	}

	// Parse until date (UNTIL field)
	if !r.Until.IsZero() {
		fields.untilYear = r.Until.Year()
		fields.untilMonth = int(r.Until.Month())
		fields.untilDay = r.Until.Day()
		fields.untilHour = r.Until.Hour()
		fields.untilMinute = r.Until.Minute()
	}

	// Parse BYDAY to 7-bit weekday mask
	// Bit 0 = Sunday, Bit 1 = Monday, ..., Bit 6 = Saturday
	for _, d := range r.ByDay {
		wd := weekdayToInt(string(d.Day))
		if wd >= 0 && wd <= 6 {
			fields.byDayMask |= 1 << wd
		}
	}

	// Parse BYMONTH to 12-bit month mask
	// Bit 0 = January, Bit 1 = February, ..., Bit 11 = December
	for _, m := range r.ByMonth {
		if m >= 1 && m <= 12 {
			fields.byMonthMask |= 1 << (m - 1)
		}
	}

	// Parse BYMONTHDAY to 31-bit day-of-month mask
	// Bit 0 = Day 1, Bit 1 = Day 2, ..., Bit 30 = Day 31
	for _, md := range r.ByMonthDay {
		if md >= 1 && md <= 31 {
			fields.byMonthDayMask |= 1 << (md - 1)
		}
	}

	// Parse BYHOUR to 24-bit hour mask
	// Bit 0 = Hour 0, Bit 1 = Hour 1, ..., Bit 23 = Hour 23
	for _, h := range r.ByHour {
		if h >= 0 && h <= 23 {
			fields.byHourMask |= 1 << h
		}
	}

	// Parse BYMINUTE to 60-bit minute mask
	// Bit 0 = Minute 0, Bit 1 = Minute 1, ..., Bit 59 = Minute 59
	for _, mn := range r.ByMinute {
		if mn >= 0 && mn <= 59 {
			fields.byMinuteMask |= 1 << mn
		}
	}

	// Parse BYSECOND to 60-bit second mask
	// Bit 0 = Second 0, Bit 1 = Second 1, ..., Bit 59 = Second 59
	for _, s := range r.BySecond {
		if s >= 0 && s <= 59 {
			fields.bySecondMask |= 1 << s
		}
	}

	// Parse WKST (week start day)
	if r.Wkst != "" {
		fields.wkst = weekdayToInt(string(r.Wkst))
	}

	return fields
}

// parseICS converts an ICS calendar to the flat byte stream format.
// Each event becomes one line of space-separated integers (30 fields per event).
// Returns the flat byte stream string, metadata, and any error.
func parseICS(data []byte) (string, []eventMetadata, error) {
	cal, err := ical.ParseCalendar(strings.NewReader(string(data)))
	if err != nil {
		return "", nil, err
	}

	var sb strings.Builder
	var metadata []eventMetadata
	eventID := 0

	for _, event := range cal.Events() {
		// Parse DTSTART (event start time)
		start := event.GetProperty(ical.ComponentPropertyDtStart)
		year, month, day, hour, minute := 2026, 1, 1, 0, 0
		if start != nil {
			year, month, day, hour, minute = parseDateTime(start.Value)
		}

		// Parse DURATION (event duration)
		duration := 60 // Default 1 hour
		if durProp := event.GetProperty(ical.ComponentPropertyDuration); durProp != nil {
			duration = parseDuration(durProp.Value)
		}

		// Parse PRIORITY (0-9, lower = higher priority)
		priority := 0
		if priProp := event.GetProperty(ical.ComponentPropertyPriority); priProp != nil {
			fmt.Sscanf(priProp.Value, "%d", &priority)
		}

		// Parse CATEGORIES
		category := 0
		if catProp := event.GetProperty(ical.ComponentPropertyCategories); catProp != nil {
			category = parseCategory(catProp.Value)
		}

		// Parse DUE (deadline for tasks)
		deadlineYear, deadlineMonth, deadlineDay, deadlineHour, deadlineMinute := -1, -1, -1, -1, -1
		if dueProps := event.GetProperties(ical.ComponentPropertyDue); len(dueProps) > 0 {
			deadlineYear, deadlineMonth, deadlineDay, deadlineHour, deadlineMinute = parseDateTime(dueProps[0].Value)
		}

		// Parse RRULE (recurrence rule)
		rules, _ := event.GetRRules()
		rruleFields := parseRRULE(rules)

		// Extract metadata (string fields that can't go in numeric matrix)
		md := eventMetadata{
			ID:      eventID,
			Summary: "Event " + fmt.Sprint(eventID),
		}

		// Get UID
		if uidProp := event.GetProperty("UID"); uidProp != nil {
			md.UID = uidProp.Value
		}

		// Get Summary
		if summaryProp := event.GetProperty("SUMMARY"); summaryProp != nil {
			md.Summary = summaryProp.Value
		}

		// Get Description
		if descProp := event.GetProperty("DESCRIPTION"); descProp != nil {
			md.Description = descProp.Value
		}

		// Get Location
		if locProp := event.GetProperty("LOCATION"); locProp != nil {
			md.Location = locProp.Value
		}

		// Get Organizer
		if orgProp := event.GetProperty("ORGANIZER"); orgProp != nil {
			md.Organizer = orgProp.Value
		}

		// Get Attendees
		for _, att := range event.GetProperties("ATTENDEE") {
			md.Attendees = append(md.Attendees, att.Value)
		}

		metadata = append(metadata, md)

		// Write the event as a flat byte stream line
		// Format: 30 space-separated integers per line
		sb.WriteString(fmt.Sprintf(
			"%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d ",
			eventID, year, month, day, hour, minute,
			duration, priority, category,
			deadlineYear, deadlineMonth, deadlineDay, deadlineHour, deadlineMinute,
			rruleFields.frequency, rruleFields.interval,
			rruleFields.untilYear, rruleFields.untilMonth, rruleFields.untilDay, rruleFields.untilHour, rruleFields.untilMinute,
			rruleFields.byDayMask, rruleFields.byMonthMask, rruleFields.byMonthDayMask,
			rruleFields.byHourMask, rruleFields.byMinuteMask, rruleFields.bySecondMask,
			rruleFields.count, rruleFields.isMutable, rruleFields.wkst,
		))
		eventID++
	}

	return sb.String(), metadata, nil
}

// getDependentID looks up a dependent event ID by the RELATED-TO property.
// Returns the dependent event's ID or -1 if not found.
func getDependentID(event *ics.VEvent, uidToID map[string]int) int {
	props := event.GetProperties(ics.ComponentPropertyRelatedTo)
	if len(props) == 0 {
		return -1
	}

	uid := strings.TrimSpace(props[0].Value)
	if uid == "" {
		return -1
	}

	if id, ok := uidToID[uid]; ok {
		return id
	}
	return -1
}

func main() {
	inputFile := flag.String("input", "", "Input ICS file")
	stdin := flag.Bool("stdin", false, "Read from stdin")
	outputFile := flag.String("output", "", "Output flat stream file (empty = stdout)")
	metadataFlag := flag.String("metadata", "", "Output metadata JSON file")
	flag.Parse()

	var r io.Reader
	if *stdin {
		r = os.Stdin
	} else if *inputFile != "" {
		f, err := os.Open(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	} else {
		fmt.Fprintf(os.Stderr, "Usage: --input <file.ics> [--output <output.txt>] [--metadata <metadata.json>]\n   or: --stdin [--output <output.txt>]\n")
		os.Exit(1)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	output, metadata, err := parseICS(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ICS: %v\n", err)
		os.Exit(1)
	}

	// Save metadata JSON if requested
	if *metadataFlag != "" && len(metadata) > 0 {
		// Get source filename for the metadata
		source := "stdin"
		if *inputFile != "" {
			source = filepath.Base(*inputFile)
		}

		md := metadataFile{
			Source: source,
			Cached: time.Now().Format("2006-01-02"),
			Events: metadata,
		}

		mdJSON, err := json.MarshalIndent(md, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling metadata: %v\n", err)
			os.Exit(1)
		}

		err = os.WriteFile(*metadataFlag, mdJSON, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing metadata file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Metadata saved to %s\n", *metadataFlag)
	}

	var out io.Writer
	if *outputFile != "" {
		f, err := os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}

	io.WriteString(out, output)
	io.WriteString(out, "\n")
}
