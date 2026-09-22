package main

import "testing"

func TestParseByteSize(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"512", 512},
		{"0", 0},
		{"1B", 1},
		{"1KB", 1000},
		{"1.5KB", 1500},
		{"1KiB", 1024},
		{"1.5GiB", 1610612736},
		{"1MB", 1000000},
		{"1MiB", 1048576},
		{"1PB", 1000000000000000},
		{"1PiB", 1 << 50},
		{"  1KiB  ", 1024},
		{"1kib", 1024},
		{"1Kib", 1024},
	}
	for _, c := range cases {
		got, err := parseByteSize(c.in)
		if err != nil {
			t.Errorf("parseByteSize(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("parseByteSize(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestParseByteSizeErrors(t *testing.T) {
	cases := []string{"", "   ", "GiB", "1.5.5GiB", "abc", "1XB"}
	for _, in := range cases {
		if _, err := parseByteSize(in); err == nil {
			t.Errorf("parseByteSize(%q) expected error, got nil", in)
		}
	}
}

func TestFormatByteSize(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0B"},
		{1, "1B"},
		{1023, "1023B"},
		{1024, "1.00KiB"},
		{1610612736, "1.50GiB"},
		{1 << 50, "1.00PiB"},
		{-1610612736, "-1.50GiB"},
		{-512, "-512B"},
	}
	for _, c := range cases {
		got := formatByteSize(c.in, 2)
		if got != c.want {
			t.Errorf("formatByteSize(%d, 2) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatByteSizePrecision(t *testing.T) {
	cases := []struct {
		in        int64
		precision int
		want      string
	}{
		{1610612736, 0, "2GiB"},
		{1610612736, 1, "1.5GiB"},
		{1610612736, 4, "1.5000GiB"},
		{1024, 0, "1KiB"},
	}
	for _, c := range cases {
		got := formatByteSize(c.in, c.precision)
		if got != c.want {
			t.Errorf("formatByteSize(%d, %d) = %q, want %q", c.in, c.precision, got, c.want)
		}
	}
}

func TestByteSizeRoundTrip(t *testing.T) {
	sizes := []int64{0, 1, 1023, 1024, 1610612736, 1 << 50}
	for _, n := range sizes {
		formatted := formatByteSize(n, 2)
		got, err := parseByteSize(formatted)
		if err != nil {
			t.Errorf("parseByteSize(%q) returned error: %v", formatted, err)
			continue
		}
		if got != n {
			t.Errorf("round trip for %d: formatted as %q, parsed back as %d", n, formatted, got)
		}
	}
}
