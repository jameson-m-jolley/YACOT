# YACOT - Yet Another Calendar Optimization Tool

## Project Overview

A Unix-like CLI tool for optimizing calendar schedules using a pipeline of optimization passes. Reads ICS files, transforms them through Uiua optimization scripts, and outputs optimized ICS files.

## Tech Stack & Expertise

- **Go** - Core application, I/O, parsing (compiled to native binary)
- **go-ical** - Pure Go iCalendar parsing library
- **Uiua** - Default optimization scripts in `scripts/uiua/` directory
- **Must be expert in both Go and Uiua**

## Architecture

```
ICS Input → Go (go-ical) → Flat Byte Stream → Uiua Pass 1 → ... → Uiua Pass N → Flat Byte Stream → Go → ICS Output
```

- **Go (native)**: Handles I/O, ICS parsing/writing, matrix format conversion
- **go-ical**: Pure Go iCalendar parsing library
- **Uiua**: Optimization algorithms as standalone filter scripts

## Flat Byte Stream Format

Space-separated integers on a single line - ALL VALUES MUST BE NUMERIC (no symbolic data):

```
id,year,month,day,hour,minute,duration_minutes,priority,category_id,deadline_year,deadline_month,deadline_day,deadline_hour,deadline_minute,frequency,interval,until_year,until_month,until_day,until_hour,until_minute,by_day_mask,by_month_mask,by_monthday_mask,by_hour_mask,by_minute_mask,by_second_mask,count,is_mutable,wkst
```

### Field Descriptions

| # | Field | Description | Range |
|---|-------|-------------|-------|
| 0 | id | Event index (0, 1, 2...) |
| 1 | year | Year (e.g., 2026) |
| 2 | month | Month (1-12) |
| 3 | day | Day (1-31) |
| 4 | hour | Hour (0-23) |
| 5 | minute | Minute (0-59) |
| 6 | duration_minutes | Duration in minutes |
| 7 | priority | Priority (0-9) |
| 8 | category_id | Category code (0=work, 1=dev, 2=personal) |
| 9 | deadline_year | Deadline year (-1 if none) |
| 10 | deadline_month | Deadline month |
| 11 | deadline_day | Deadline day |
| 12 | deadline_hour | Deadline hour |
| 13 | deadline_minute | Deadline minute |
| 14 | frequency | Recurrence frequency (-1 to 6, see below) |
| 15 | interval | Repeat interval (default 1) |
| 16 | until_year | Until date year |
| 17 | until_month | Until date month |
| 18 | until_day | Until date day |
| 19 | until_hour | Until date hour |
| 20 | until_minute | Until date minute |
| 21 | by_day_mask | 7-bit weekday mask (SU=bit0, MO=bit1, ..., SA=bit6) |
| 22 | by_month_mask | 12-bit month mask (Jan=bit0, Feb=bit1, ..., Dec=bit11) |
| 23 | by_monthday_mask | 31-bit day-of-month mask (day 1-31) |
| 24 | by_hour_mask | 24-bit hour mask (0-23) |
| 25 | by_minute_mask | 60-bit minute mask (0-59) |
| 26 | by_second_mask | 60-bit second mask (0-59) |
| 27 | count | Number of occurrences (0 if unset) |
| 28 | is_mutable | 1=mutable, 0=immutable |
| 29 | wkst | Week start day (0=SU, 1=MO, ..., 6=SA) |

### Bitmask Formats

- **by_day_mask**: 7-bit mask for days of week
  - Sunday = bit0 (value 1)
  - Monday = bit1 (value 2)
  - Tuesday = bit2 (value 4)
  - Wednesday = bit3 (value 8)
  - Thursday = bit4 (value 16)
  - Friday = bit5 (value 32)
  - Saturday = bit6 (value 64)
  - Every day = 127

- **by_month_mask**: 12-bit mask for months
  - January = bit0 (value 1)
  - February = bit1 (value 2)
  - ...
  - December = bit11 (value 4096)

- **by_monthday_mask**: 31-bit mask for day of month (days 1-31)
  - Day 1 = bit0, Day 2 = bit1, ..., Day 31 = bit30

- **by_hour_mask**: 24-bit mask for hours (0-23)
  - Hour 0 = bit0, Hour 1 = bit1, ..., Hour 23 = bit23

- **by_minute_mask**: 60-bit mask for minutes (0-59)
- **by_second_mask**: 60-bit mask for seconds (0-59)

### Sentinel Values
- `-1` = no value / unset field

### Frequency Codes

| Code | RRULE Value   | Description                     | Example Use Case |
|------|-------------|---------------------------------|-------------------------------|
| -1   | (none)        | No recurrence (one-time event)  | Single meeting or appointment |
|  0   | SECONDLY      | Every second                    | Rarely used |
|  1   | MINUTELY      | Every minute                    | Rarely used |
|  2   | HOURLY        | Every hour                      | Reminder every hour |
|  3   | DAILY        | Every day                       | Daily standup |
|  4   | WEEKLY       | Every week                      | Weekly team meeting |
|  5   | MONTHLY      | Every month                     | Monthly review |
|  6   | YEARLY      | Every year                      | Annual birthday |

### Category Codes
- 0 = work
- 1 = dev
- 2 = personal
- Other values can be added as needed


**Example**:
```
0 2023 1 1 0 0 60 0 0 -1 -1 -1 -1 -1 -1 -1 -1 -1 ...
```

This format allows Uiua to parse and reshape into a matrix using `parse` and array operations (this MUST BE FLAT no exceptions UIUA like reading in flat arrays).

### Sentinel Values
- `deadline_year/month/day = -1` (no deadline)
- `category` is always >= 0 | -1 means no category

## Dependencies
   uiua must be installed 

### Required
- **Go** - Install via `brew install go` or from https://go.dev/dl/
- **Uiua** - Must be in PATH for running optimization passes

### Optional (for development)
- GoLand or VS Code with Go extension

## Binaries

Native executables:

| Binary | Input | Output |
|--------|-------|--------|
| `bin/ics-to-matrix` | ICS file or stdin | flat byte stream (file or stdout) |
| `bin/matrix-to-ics` | flat byte stream file or stdin | ICS file or stdout |

### CLI Usage

```sh
# Full pipeline
./bin/ics-to-matrix --input a.ics | uiua scripts/uiua/block.ua keys=value ...| ./bin/matrix-to-ics --output out.ics

# With stdin/stdout
./bin/ics-to-matrix --stdin < a.ics | uiua scripts/uiua/block.ua keys=value ...| ./bin/matrix-to-ics --output out.ics
```
### building 

the build need to be handled by make
all go bin need to be under a make command

## CLI Reference

### ics-to-matrix

| Flag | Description | Default |
|------|-------------|---------|
| `--input` | Input ICS file | empty |
| `--stdin` | Read from stdin | false |
| `--output` | Output flat stream file | empty → stdout |
| `--metadata` | Output metadata JSON file | empty → metadata lost |

**Either `--input` or `--stdin` is required (not both)**

# important
output must be on one line the uiua logic depends on this fact (DO NOT print newline chars)


### matrix-to-ics

| Flag | Description | Default |
|------|-------------|---------|
| `--input` | Input flat stream file | empty |
| `--stdin` | Read from stdin | false |
| `--output` | Output ICS file | empty → stdout |
| `--metadata` | Input metadata JSON file | empty → error |

**Either `--input` or `--stdin` is required (not both)**

## Available Passes

### Uiua scripts in `scripts` dir:
# all fallow the fallowing convention
   -take flat numeric stream from stdin(one line separated by a spaces)
   -takes args as key=val key1=val1 ... keyn=valn
   -prints a flat numeric stream(one line separated by a spaces)

- `block.ua` - Time blocking
- `context-switch.ua` - Group similar tasks
- `constraint-lock.ua` - Enforce dependencies and deadlines
- `time-pad.ua` - Add buffer time between tasks
- `free-time.ua` - Aggregate free time slots
- `priority.ua` - Schedule by priority
- `prune.ua` - Remove redundant tasks (mostly removes past events)
- `tools.ua` - UIua scripts in `scripts`

Each pass:
1. Reads flat byte stream from stdin
2. Writes optimized flat byte stream to stdout
3. Preserves one line of output per pass as a flat numeric stream

## Implementation Phases

### Phase 1: Core Infrastructure
0. Add comprehensive tests for correctness and properties
1. Create `ics-to-matrix` (Go + go-ical)
2. Create `matrix-to-ics` (Go)

### Phase 2: Uiua Passes
3. Create Uiua optimization passes in `scripts`

### Phase 3: Polish & Testing
4. Add comprehensive tests to existing tests
5. Error handling refinements
6. Documentation updates

## Engineering Constraints
- **use property based tests**
- **use clean code**
- **~1GB RAM target** (soft cap, can exceed if needed)
- **Sub-second execution** target
- **Memory-safe language** (Go and uiua)
- **Pure functional style** preferred over imperative
- **Full test coverage** - all functions must have tests
- **Constraint comments** - document time and memory budgets for each function
- **NO compiler/interpreter warnings** - must build clean
- **NO feature creep** - stick to MVP scope defined above
- **Isolation is king** - modular passes, clear interfaces, test isolation

## human-todo.txt

Track technical debt and issues for future humans to address. Format:

```
# Issue title
- Description of the problem
- Suggested fix or interface
```

When implementing a feature that isn't fully complete:
1. Define the interface
2. Write the test case
3. Print a WIP error message explaining the feature is not yet implemented

## Writing Custom Passes

Each pass is an interpreter script that:
1. Reads flat byte stream from stdin
2. Writes optimized flat byte stream to stdout
3. Preserves one event per line format

### Uiua pass structure:
```
# Your optimization algorithm
# Input: Flat byte stream (space-separated integers, one event per line)
# Output: Optimized flat byte stream (same format)
```

### Uiua Flat Stream Handling

The flat byte stream is a space-separated list of integers. In Uiua:
- Perform optimization operations
- use tools that are in tools.ua to ensure compatibility (human maid and verified).


## Development Workflow

### Adding New Go Binaries
1. Create a new directory under `cmd/` for your binary (e.g., `cmd/new-binary`).
2. Add a `main.go` file with your Go application logic.
3. Add a new target to the `makefile` similar to `ics-to-matrix` and `matrix-to-ics` to build your binary.
4. Run `make all` to build your new binary.

### Adding New Uiua Optimization Passes
1. Create a new Uiua script in the `scripts/uiua/` directory (e.g., `scripts/uiua/my-new-pass.ua`).
2. Ensure your script reads from stdin and writes to stdout, preserving the flat byte stream format.
3. You can test your pass using the integration test pipeline:
   `./bin/ics-to-matrix --input a.ics | uiua scripts/uiua/my-new-pass.ua | ./bin/matrix-to-ics --output out.ics`

### Running Tests
- Run all Go unit and integration tests: `make test`
- Run Go fuzz tests for `ics-to-matrix`: `make fuzz-ics`
- Run Go fuzz tests for `matrix-to-ics`: `make fuzz-matrix`
- View all available make targets and their descriptions: `make help`

## Testing Strategy

- Unit tests for all Go functions
- Integration tests for ICS roundtrip (input → matrix → output)
- Validate interpreter passes produce valid output
- Test against sample ICS files in `testdata/`


