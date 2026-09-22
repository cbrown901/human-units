package main

import (
	"fmt"
	"strconv"
	"strings"
)

// sizeUnits is ordered longest-suffix-first so "KiB" is tried before "B"
// matches on the same string.
var sizeUnits = []struct {
	suffix     string
	multiplier float64
}{
	{"PiB", 1 << 50},
	{"TiB", 1 << 40},
	{"GiB", 1 << 30},
	{"MiB", 1 << 20},
	{"KiB", 1 << 10},
	{"PB", 1e15},
	{"TB", 1e12},
	{"GB", 1e9},
	{"MB", 1e6},
	{"KB", 1e3},
	{"B", 1},
}

// parseByteSize turns a human size like "1.5GiB" or "512" (bytes) into a
// byte count. Unit matching is case-insensitive.
func parseByteSize(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty size")
	}
	lower := strings.ToLower(s)
	for _, u := range sizeUnits {
		if !strings.HasSuffix(lower, strings.ToLower(u.suffix)) {
			continue
		}
		numPart := strings.TrimSpace(s[:len(s)-len(u.suffix)])
		if numPart == "" {
			continue
		}
		n, err := strconv.ParseFloat(numPart, 64)
		if err != nil {
			continue
		}
		return int64(n * u.multiplier), nil
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("not a recognized size: %q", s)
	}
	return int64(n), nil
}

// formatByteSize renders a byte count using the largest IEC unit that keeps
// the number at or above 1, with the given number of decimal places.
func formatByteSize(n int64, precision int) string {
	abs := float64(n)
	if n < 0 {
		abs = -abs
	}
	units := []struct {
		suffix string
		size   float64
	}{
		{"PiB", 1 << 50},
		{"TiB", 1 << 40},
		{"GiB", 1 << 30},
		{"MiB", 1 << 20},
		{"KiB", 1 << 10},
	}
	for _, u := range units {
		if abs >= u.size {
			value := float64(n) / u.size
			return fmt.Sprintf("%.*f%s", precision, value, u.suffix)
		}
	}
	return fmt.Sprintf("%dB", n)
}
