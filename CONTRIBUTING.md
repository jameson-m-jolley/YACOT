# Contributing to YACOT

Thank you for your interest in contributing to YACOT. This document provides guidelines for contributing to the project.

## Getting Started

### Prerequisites

- **Go** 1.21 or higher
- **Uiua** (for running optimization passes during testing)
- **git** for version control

### Setup

```sh
# Clone the repository
git clone https://github.com/yourorg/yacot.git
cd yacot

# Install dependencies
go mod download

# Build binaries
make all

# Verify setup
make test
```

## Development Workflow

### Running Tests

```sh
# Run all tests (unit + integration)
make test

# Run specific test package
go test -v ./cmd/ics-to-matrix

# Run fuzz tests
make fuzz
```

### Building

```sh
# Build all binaries
make all

# Build specific binary
go build -o bin/ics-to-matrix ./cmd/ics-to-matrix
go build -o bin/matrix-to-ics ./cmd/matrix-to-ics
```

### Code Style

#### Go Conventions

- **Standard Go formatting**: Use `go fmt` before committing
- **Error handling**: Always check and handle errors (no `_ == nil` checks)
- **Naming**: Use camelCase for functions/variables, PascalCase for exported
- **Documentation**: Add docstrings for all exported functions

```go
// Good
func parseDateTime(dt string) (year, month, day, hour, minute int)

// Bad
func ParseDt(dt string) (int, int, int, int, int)
```

#### Uiua Script Conventions

- Read from stdin, write to stdout
- Preserve flat byte stream format (space-separated integers)
- One event per line (tools.ua has the reshaping logic)
- args fallow the same logic as https query parameters ex: key1=val1 key2=val2 (look in tools.ua)
- No trailing newlines in output (keeps format consistent for users)

## Submitting Changes
   - Ensure your changes are well-documented and include tests
   - Ensure the changes are compatible with the existing codebase
   - if posable test on vm or more than one system (eg windows vm or open bsd vm)
   - submit PR with a clear description of the changes and what they do

### Pull Request Process

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/my-feature`
3. **Make** your changes
4. **Run** tests: `make test`
5. **Commit** with a clear message:
   ```
   Add parsing for BYDAY RRULE field
   
   Parses the BYDAY component of RFC 5545 recurrence rules
   and converts to a 7-bit weekday bitmask.
   ```
6. **Push** to your fork
7. **Create** a pull request

### Commit Message Guidelines

- Use present tense: "Add feature" not "Added feature"
- First line: 50 characters or less
- Detailed description after blank line if needed
- Reference issues: "Fixes #123" or "Closes #456"

## Testing Requirements

All code changes must pass existing tests. New functionality requires tests.

### Test Types

| Test Type | Description | Required |
|----------|------------|----------|
| Unit tests | Test individual functions | Yes |
| Integration | Test full pipeline | Yes |
| Fuzz tests | Property-based testing | Preferred |

### Running Specific Tests

```sh
# Unit tests only
go test ./cmd/...

# Fuzz tests
go test -fuzz=FuzzParseICS ./cmd/ics-to-matrix

# Verbose output
go test -v ./cmd/ics-to-matrix
```

## Architecture Notes

### Flat Byte Stream Format

The flat byte stream format is critical. When modifying:
- Always maintain 30 fields per event
- Use `-1` as sentinel for unset values
- Keep numeric format (no symbolic data)

See [AGENTS.md](./AGENTS.md) for complete specification.

### RRULE Parsing

RFC 5545 recurrence rules are mapped as:
- `Frequency`: -1 to 6 (see AGENTS.md for codes)
- `Interval`: positive integer (default 1)
- `Until`: datetime components (fields 17-21)
- `BYDAY`: 7-bit bitmask (bit 0 = Sunday)
- `BYMONTH`: 12-bit bitmask (bit 0 = January)
- `WKST`: 0-6 (week start day)

### Bitmask Formats

| Field | Bits | Range | Format |
|-------|------|-------|--------|
| by_day_mask | 7 | SU-SA | bit 0 = SU, bit 6 = SA |
| by_month_mask | 12 | Jan-Dec | bit 0 = Jan, bit 11 = Dec |
| by_monthday_mask | 31 | Day 1-31 | bit 0 = Day 1 |
| by_hour_mask | 24 | Hour 0-23 | bit 0 = Hour 0 |
| by_minute_mask | 60 | Min 0-59 | bit 0 = Min 0 |
| by_second_mask | 60 | Sec 0-59 | bit 0 = Sec 0 |

## Common Tasks

### Adding a New RRULE Field

1. Add field to `eventFields` struct in `cmd/ics-to-matrix/main.go`
2. Parse in `parseRRULE()` function
3. Add to format string (maintain 30 field order)
4. Add test case
5. Update AGENTS.md documentation

### Adding a New Uiua Pass

1. Create `scripts/new-pass.ua`
2. Read flat stream from stdin
3. Transform data
4. Write flat stream to stdout
5. Add to test pipeline

## Questions

For questions or discussions:
- Open an issue on GitHub
- Check existing issues before creating new ones

## Recognition

Contributors will be listed in the project (with permission).