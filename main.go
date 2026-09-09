// Command human-units converts byte sizes and durations between
// human-readable notation ("1.5GiB", "2h30m") and raw numbers (bytes,
// seconds), one value per line.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func convert(kind, line string) (string, error) {
	switch kind {
	case "size":
		if n, err := strconv.ParseInt(line, 10, 64); err == nil {
			return formatByteSize(n), nil
		}
		n, err := parseByteSize(line)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(n, 10), nil
	case "duration":
		if n, err := strconv.ParseInt(line, 10, 64); err == nil {
			return formatDuration(n), nil
		}
		n, err := parseDuration(line)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(n, 10), nil
	default:
		return "", fmt.Errorf("unknown kind %q (want \"size\" or \"duration\")", kind)
	}
}

func processInput(r io.Reader, kind string, out, errOut io.Writer) bool {
	ok := true
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		result, err := convert(kind, line)
		if err != nil {
			fmt.Fprintf(errOut, "%s: %v\n", line, err)
			ok = false
			continue
		}
		fmt.Fprintln(out, result)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(errOut, err)
		ok = false
	}
	return ok
}

func main() {
	kind := flag.String("kind", "size", `what to convert: "size" or "duration"`)
	flag.Parse()

	if *kind != "size" && *kind != "duration" {
		fmt.Fprintf(os.Stderr, "invalid -kind %q (want \"size\" or \"duration\")\n", *kind)
		os.Exit(2)
	}

	args := flag.Args()
	ok := true

	if len(args) == 0 {
		ok = processInput(os.Stdin, *kind, os.Stdout, os.Stderr)
	} else {
		for _, name := range args {
			f, err := os.Open(name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				ok = false
				continue
			}
			if !processInput(f, *kind, os.Stdout, os.Stderr) {
				ok = false
			}
			f.Close()
		}
	}

	if !ok {
		os.Exit(1)
	}
}
