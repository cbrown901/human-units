package main

import "testing"

func TestParseDuration(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"0", 0},
		{"95400", 95400},
		{"1s", 1},
		{"1m", 60},
		{"1h", 3600},
		{"1d", 86400},
		{"1d2h30m", 95400},
		{"2h30m", 9000},
		{"1.5h", 5400},
		{"500ms", 0},
		{"1s500ms", 1},
		{"  1d2h30m  ", 95400},
	}
	for _, c := range cases {
		got, err := parseDuration(c.in)
		if err != nil {
			t.Errorf("parseDuration(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("parseDuration(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestParseDurationErrors(t *testing.T) {
	cases := []string{"", "   ", "d", "1x", "1d2", "abc"}
	for _, in := range cases {
		if _, err := parseDuration(in); err == nil {
			t.Errorf("parseDuration(%q) expected error, got nil", in)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0s"},
		{1, "1s"},
		{60, "1m"},
		{3600, "1h"},
		{86400, "1d"},
		{95400, "1d2h30m"},
		{9000, "2h30m"},
		{-95400, "-1d2h30m"},
		{-1, "-1s"},
	}
	for _, c := range cases {
		got := formatDuration(c.in)
		if got != c.want {
			t.Errorf("formatDuration(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestDurationRoundTrip(t *testing.T) {
	seconds := []int64{0, 1, 60, 3600, 86400, 95400, -95400}
	for _, n := range seconds {
		formatted := formatDuration(n)
		got, err := parseDuration(formatted)
		if err != nil {
			t.Errorf("parseDuration(%q) returned error: %v", formatted, err)
			continue
		}
		if got != n {
			t.Errorf("round trip for %d: formatted as %q, parsed back as %d", n, formatted, got)
		}
	}
}
