package main

import (
	"strings"
	"testing"
)

func TestProcessInputPlainText(t *testing.T) {
	var out, errOut strings.Builder
	ok := processInput(strings.NewReader("1.5GiB\nnotasize\n"), "size", 2, false, &out, &errOut)
	if ok {
		t.Error("expected ok = false due to bad line")
	}
	if got, want := out.String(), "1610612736\n"; got != want {
		t.Errorf("out = %q, want %q", got, want)
	}
	if !strings.Contains(errOut.String(), "notasize") {
		t.Errorf("errOut = %q, want it to mention the bad line", errOut.String())
	}
}

func TestProcessInputJSON(t *testing.T) {
	var out, errOut strings.Builder
	ok := processInput(strings.NewReader("1.5GiB\nnotasize\n"), "size", 2, true, &out, &errOut)
	if ok {
		t.Error("expected ok = false due to bad line")
	}
	want := `{"input":"1.5GiB","output":"1610612736"}` + "\n" +
		`{"input":"notasize","error":"not a recognized size: \"notasize\""}` + "\n"
	if got := out.String(); got != want {
		t.Errorf("out = %q, want %q", got, want)
	}
	if errOut.String() != "" {
		t.Errorf("errOut = %q, want empty in JSON mode", errOut.String())
	}
}

func TestProcessInputJSONDuration(t *testing.T) {
	var out, errOut strings.Builder
	ok := processInput(strings.NewReader("95400\n"), "duration", 2, true, &out, &errOut)
	if !ok {
		t.Error("expected ok = true")
	}
	want := `{"input":"95400","output":"1d2h30m"}` + "\n"
	if got := out.String(); got != want {
		t.Errorf("out = %q, want %q", got, want)
	}
}
