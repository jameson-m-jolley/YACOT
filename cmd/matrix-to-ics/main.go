package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// eventMetadata holds string data for events from JSON.
type eventMetadata struct {
	ID          int      `json:"id"`
	UID         string   `json:"uid"`
	Summary     string   `json:"summary"`
	Description string   `json:"description,omitempty"`
	Location    string   `json:"location,omitempty"`
	Organizer   string   `json:"organizer,omitempty"`
	Attendees   []string `json:"attendees,omitempty"`
}

// metadataJSON represents the JSON file containing all event metadata.
type metadataJSON struct {
	Source string          `json:"source"`
	Cached string          `json:"cached"`
	Events []eventMetadata `json:"events"`
}

func escapeICSText(text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, ";", "\\;")
	text = strings.ReplaceAll(text, ",", "\\,")
	text = strings.ReplaceAll(text, "\n", "\\n")
	text = strings.ReplaceAll(text, "\r", "")
	return text
}

func foldLine(line string) string {
	const maxLen = 75
	if len(line) <= maxLen {
		return line
	}
	var result strings.Builder
	for len(line) > maxLen {
		result.WriteString(line[:maxLen])
		result.WriteString(crlf)
		line = " " + line[maxLen:]
	}
	result.WriteString(line)
	return result.String()
}

func writeProperty(sb *strings.Builder, propName, value string) {
	fullLine := propName + ":" + value
	written := foldLine(fullLine)
	sb.WriteString(written)
	sb.WriteString(crlf)
}

func formatDTStamp() string {
	return time.Now().UTC().Format("20060102T150405Z")
}

const crlf = "\r\n"

func parseMatrix(data string, metadata []eventMetadata) (string, error) {
	metadataMap := make(map[int]eventMetadata)
	for _, m := range metadata {
		metadataMap[m.ID] = m
	}
	var sb strings.Builder

	sb.WriteString("BEGIN:VCALENDAR" + crlf)
	sb.WriteString("VERSION:2.0" + crlf)
	writeProperty(&sb, "PRODID", "-//YACOT//Calendar Optimizer//EN")
	sb.WriteString("METHOD:PUBLISH" + crlf)
	sb.WriteString("X-WR-CALNAME:YACOT Calendar" + crlf)

	normalized := strings.Join(strings.Fields(data), " ")
	numbers := strings.Fields(normalized)

	// Allow empty input (produces empty calendar)
	numEvents := 0
	if len(numbers) > 0 && len(numbers)%30 != 0 {
		return "", fmt.Errorf("invalid matrix: expected multiple of 30 fields, got %d", len(numbers))
	}
	if len(numbers) > 0 {
		numEvents = len(numbers) / 30
	}

	uidCount := make(map[string]int)

	for eventID := 0; eventID < numEvents; eventID++ {
		fields := numbers[eventID*30 : (eventID+1)*30]

		year, _ := strconv.Atoi(fields[1])
		month, _ := strconv.Atoi(fields[2])
		day, _ := strconv.Atoi(fields[3])
		hour, _ := strconv.Atoi(fields[4])
		minute, _ := strconv.Atoi(fields[5])
		duration, _ := strconv.Atoi(fields[6])
		priority, _ := strconv.Atoi(fields[7])
		category, _ := strconv.Atoi(fields[8])

		dtStart := fmt.Sprintf("%04d%02d%02dT%02d%02d00", year, month, day, hour, minute)

		durHours := duration / 60
		durMins := duration % 60
		var durStr string
		if durHours > 0 && durMins > 0 {
			durStr = fmt.Sprintf("PT%dH%dM", durHours, durMins)
		} else if durHours > 0 {
			durStr = fmt.Sprintf("PT%dH", durHours)
		} else {
			durStr = fmt.Sprintf("PT%dM", durMins)
		}

		sb.WriteString("BEGIN:VEVENT" + crlf)

		md := metadataMap[eventID]
		var uid string
		if md.UID != "" {
			uidCount[md.UID]++
			if uidCount[md.UID] > 1 {
				uid = fmt.Sprintf("%s-%d", md.UID, uidCount[md.UID])
			} else {
				uid = md.UID
			}
		} else {
			uid = fmt.Sprintf("event-%d@yacot", eventID)
		}
		writeProperty(&sb, "UID", uid)
		writeProperty(&sb, "DTSTAMP", formatDTStamp())
		writeProperty(&sb, "DTSTART", dtStart)
		writeProperty(&sb, "DURATION", durStr)

		if priority > 0 {
			writeProperty(&sb, "PRIORITY", fmt.Sprintf("%d", priority))
		}

		catName := "work"
		if category == 1 {
			catName = "dev"
		} else if category == 2 {
			catName = "personal"
		}
		writeProperty(&sb, "SUMMARY", escapeICSText(md.Summary))

		if md.Description != "" {
			writeProperty(&sb, "DESCRIPTION", escapeICSText(md.Description))
		}
		if md.Location != "" {
			writeProperty(&sb, "LOCATION", escapeICSText(md.Location))
		}
		if md.Organizer != "" {
			writeProperty(&sb, "ORGANIZER", escapeICSText(md.Organizer))
		}
		for _, att := range md.Attendees {
			writeProperty(&sb, "ATTENDEE", escapeICSText(att))
		}
		writeProperty(&sb, "CATEGORIES", catName)

		deadlineYear := fields[9]
		deadlineMonth := fields[10]
		deadlineDay := fields[11]
		deadlineHour := fields[12]
		deadlineMinute := fields[13]
		if deadlineYear != "-1" {
			dueStr := fmt.Sprintf("%s%s%sT%s%s00", deadlineYear, deadlineMonth, deadlineDay, deadlineHour, deadlineMinute)
			writeProperty(&sb, "DUE", dueStr)
		}

		frequency := fields[14]
		interval := fields[15]
		untilYear := fields[16]
		untilMonth := fields[17]
		untilDay := fields[18]

		if frequency != "-1" {
			var freqStr string
			switch frequency {
			case "0":
				freqStr = "MINUTELY"
			case "1":
				freqStr = "MINUTELY"
			case "2":
				freqStr = "HOURLY"
			case "3":
				freqStr = "DAILY"
			case "4":
				freqStr = "WEEKLY"
			case "5":
				freqStr = "MONTHLY"
			case "6":
				freqStr = "YEARLY"
			}

			rrule := fmt.Sprintf("FREQ=%s", freqStr)
			if interval != "1" && interval != "-1" {
				rrule += fmt.Sprintf(";INTERVAL=%s", interval)
			}
			if untilYear != "-1" {
				untilY, _ := strconv.Atoi(untilYear)
				untilM, _ := strconv.Atoi(untilMonth)
				untilD, _ := strconv.Atoi(untilDay)
				rrule += fmt.Sprintf(";UNTIL=%04d%02d%02dT235959Z", untilY, untilM, untilD)
			}
			writeProperty(&sb, "RRULE", rrule)
		}

		sb.WriteString("END:VEVENT" + crlf)
	}

	sb.WriteString("END:VCALENDAR" + crlf)

	return sb.String(), nil
}

func main() {
	inputFile := flag.String("input", "", "Input flat stream file")
	stdin := flag.Bool("stdin", false, "Read from stdin")
	outputFile := flag.String("output", "", "Output ICS file (empty = stdout)")
	metadataFlag := flag.String("metadata", "", "Input metadata JSON file")
	flag.Parse()

	var metadata []eventMetadata
	if *metadataFlag != "" {
		mdData, err := os.ReadFile(*metadataFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading metadata file: %v\n", err)
			os.Exit(1)
		}
		var md metadataJSON
		if err := json.Unmarshal(mdData, &md); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing metadata JSON: %v\n", err)
			os.Exit(1)
		}
		metadata = md.Events
	}

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
		fmt.Fprintf(os.Stderr, "Usage: --input <file.txt> [--output <output.ics>]\n   or: --stdin [--output <output.ics>]\n")
		os.Exit(1)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	output, err := parseMatrix(string(data), metadata)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing matrix: %v\n", err)
		os.Exit(1)
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
}
