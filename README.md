# YACOT - Yet Another Calendar Optimization Tool

A Unix-like CLI tool for optimizing calendar schedules using a pipeline of optimization passes. Reads ICS files, transforms them through Uiua optimization scripts, and outputs optimized ICS files.

## Quick Start

### Installation

```sh
# Clone the repository
git clone https://github.com/yourorg/yacot.git
cd yacot

# Build binaries
make all

# Verify installation
./bin/ics-to-matrix --help
```

### Basic Usage

```sh
# Convert ICS to flat byte stream
./bin/ics-to-matrix --input events.ics --output matrix.txt

# Convert flat byte stream back to ICS
./bin/matrix-to-ics --input matrix.txt --output optimized.ics
```

### Full Pipeline with Uiua

```sh
# ICS → matrix → Uiua pass → ICS
./bin/ics-to-matrix --input events.ics | \
  uiua scripts/pass.ua | \
  ./bin/matrix-to-ics --stdin > optimized.ics
```

## Architecture

```
┌─────────────┐    ┌────────────────┐    ┌─────────────┐
│  ICS Input  │───▶│ ics-to-matrix  │───▶│ Flat Stream │
└─────────────┘    └────────────────┘    └─────────────┘
                                              │
                                              ▼
┌─────────────┐    ┌────────────────┐    ┌─────────────┐
│ ICS Output │◀───│ matrix-to-ics  │◀───│  Uiua Pass │
└─────────────┘    └────────────────┘    └─────────────┘
```

- **Go (native)**: Handles I/O, ICS parsing/writing, matrix format conversion
- **go-ical**: Pure Go iCalendar parsing library
- **Uiua**: Optimization algorithms as standalone filter scripts

## CLI Reference

### ics-to-matrix

| Flag | Description | Default |
|------|------------|---------|
| `--input` | Input ICS file | (required) |
| `--stdin` | Read from stdin | false |
| `--output` | Output file | stdout |

### matrix-to-ics

| Flag | Description | Default |
|------|------------|---------|
| `--input` | Input flat stream file | (required) |
| `--stdin` | Read from stdin | false |
| `--output` | Output ICS file | stdout |

## Flat Byte Stream Format

The flat byte stream format is space-separated integers on a single line. See [AGENTS.md](./AGENTS.md) for the complete specification.

**30 fields per event:**
```
id year month day hour minute duration priority category deadline_year deadline_month deadline_day deadline_hour deadline_minute frequency interval until_year until_month until_day until_hour until_minute by_day_mask by_month_mask by_monthday_mask by_hour_mask by_minute_mask by_second_mask count is_mutable wkst
```

## Examples

### Convert and Modify with Uiua

```sh
# Remove past events, keep only future ones
./bin/ics-to-matrix --input calendar.ics | \
  uiua scripts/prune.ua | \
  ./bin/matrix-to-ics --stdin > future.ics
```

### View Event Data

```sh
# See raw event data
./bin/ics-to-matrix --input calendar.ics > matrix.txt
cat matrix.txt
```

## Requirements

- **Go** 1.21 or higher
- **Uiua** (for running optimization passes)
- **go-ical** v0.3.5 (included via go.mod)

## Building

```sh
# Build all binaries
make all

# Build specific binary
make bin/ics-to-matrix

# Clean build artifacts
make clean
```

## Testing

```sh
# Run all tests (unit + integration)
make test

# Run Go fuzz tests
make fuzz
```

## Scripts

Optimization passes are written in Uiua and located in `scripts/`:

| Script | Purpose |
|-------|---------|
| `pass.ua` | Identity pass (no changes) |
| `prune.ua` | Remove past events |
| `remove_dupes.ua` | Remove duplicate events |
| `tools.ua` | Helper functions |

## Project Structure

```
.
├── bin/                    # Compiled binaries
│   ├── ics-to-matrix
│   └── matrix-to-ics
├── cmd/                    # Go source code
│   ├── ics-to-matrix/
│   └── matrix-to-ics/
├── scripts/                 # Uiua optimization passes
│   ├── pass.ua
│   ├── prune.ua
│   └── tools.ua
├── testdata/               # Test fixtures
├── AGENTS.md               # Technical specification
└── makefile               # Build automation
```

## License

See repository for license details.

## Further Reading

- [AGENTS.md](./AGENTS.md) - Technical specification and field definitions
- [RFC 5545](https://tools.ietf.org/html/rfc5545) - iCalendar specification