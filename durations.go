package main

import (
	"fmt"
	"strconv"
	"strings"
)

// durationUnits maps a unit suffix to its length in seconds. Days, weeks,
// and years are included because time.ParseDuration deliberately leaves
// them out, but they show up constantly in log retention and job
// scheduling. A year is treated as a fixed 365 days rather than a
// calendar year, since there's no calendar context to resolve leap years
// against here.
var durationUnits = map[string]float64{
	"ns": 1e-9,
	"us": 1e-6,
	"ms": 1e-3,
	"s":  1,
	"m":  60,
	"h":  3600,
	"d":  86400,
	"w":  7 * 86400,
	"y":  365 * 86400,
}

func isDigitOrDot(b byte) bool {
	return (b >= '0' && b <= '9') || b == '.'
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// parseDuration turns a human duration like "1d2h30m" or a plain number
// (seconds) into a count of seconds.
func parseDuration(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty duration")
	}
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n, nil
	}

	orig := s
	var total float64
	for len(s) > 0 {
		i := 0
		for i < len(s) && isDigitOrDot(s[i]) {
			i++
		}
		if i == 0 {
			return 0, fmt.Errorf("invalid duration %q", orig)
		}
		numPart := s[:i]
		s = s[i:]

		j := 0
		for j < len(s) && isLetter(s[j]) {
			j++
		}
		if j == 0 {
			return 0, fmt.Errorf("missing unit after %q in %q", numPart, orig)
		}
		unit := s[:j]
		s = s[j:]

		n, err := strconv.ParseFloat(numPart, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number %q in %q", numPart, orig)
		}
		mult, ok := durationUnits[unit]
		if !ok {
			return 0, fmt.Errorf("unknown duration unit %q in %q", unit, orig)
		}
		total += n * mult
	}
	return int64(total), nil
}

// formatDuration renders a count of seconds as "1d2h3m4s", dropping any
// leading components that are zero.
func formatDuration(totalSeconds int64) string {
	if totalSeconds == 0 {
		return "0s"
	}
	neg := totalSeconds < 0
	n := totalSeconds
	if neg {
		n = -n
	}

	days := n / 86400
	n %= 86400
	hours := n / 3600
	n %= 3600
	minutes := n / 60
	seconds := n % 60

	var b strings.Builder
	if neg {
		b.WriteByte('-')
	}
	if days > 0 {
		fmt.Fprintf(&b, "%dd", days)
	}
	if hours > 0 {
		fmt.Fprintf(&b, "%dh", hours)
	}
	if minutes > 0 {
		fmt.Fprintf(&b, "%dm", minutes)
	}
	if seconds > 0 || b.Len() == 0 || (neg && b.Len() == 1) {
		fmt.Fprintf(&b, "%ds", seconds)
	}
	return b.String()
}
