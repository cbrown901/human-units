# human-units

Log lines and config files are full of numbers like `2.5GiB` or `3d4h` that
are easy for a person to read and annoying to do math on. `human-units`
converts them to and from plain numbers (bytes, seconds) so you can pipe them
into `awk`, `sort -n`, or a spreadsheet.

It reads one value per line, from a file, multiple files, or stdin, and
prints the converted value. Direction is automatic: numbers go to
human-readable notation, human-readable notation goes to numbers.

## Build

```
go build -o human-units .
```

## Usage

Convert sizes (default mode):

```
$ echo "1.5GiB" | human-units
1610612736

$ echo "1610612736" | human-units
1.50GiB

$ human-units sizes.txt
```

Convert durations:

```
$ echo "1d2h30m" | human-units -kind duration
95400

$ echo "95400" | human-units -kind duration
1d2h30m
```

Multiple files are read in order, and stdin is used only when no file
arguments are given:

```
$ human-units -kind duration logs/*.txt < /dev/null
```

## Supported units

Sizes accept decimal (`KB`, `MB`, `GB`, `TB`, `PB`, base 1000) and binary
(`KiB`, `MiB`, `GiB`, `TiB`, `PiB`, base 1024) suffixes, plus a bare `B` or
no suffix at all for a raw byte count.

Durations accept `ns`, `us`, `ms`, `s`, `m`, `h`, and `d` (days), combined
left to right like `1d2h30m`, or a bare integer for seconds.

## Status

Early skeleton. Parsing is line-oriented and unit-tested manually so far;
see the issue tracker for what's missing (weeks/years, fractional output
precision control, JSON output).

## License

MIT, see [LICENSE](LICENSE).
