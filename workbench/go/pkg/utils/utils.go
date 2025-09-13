package utils

import (
	"bufio"
	"os"
)

type options struct {
	maxCapacity int
}

type Option func(*options)

// WithMaxCapacity returns an Option that sets the maximum token size (in
// bytes) used by the underlying bufio.Scanner within ScanStdin. This value
// controls the largest single line that can be scanned without hitting
// bufio.ErrTooLong. If a line exceeds this capacity, scanning stops and any
// lines collected up to that point are returned by ScanStdin.
//
// Use this to increase the default limit (bufio.MaxScanTokenSize) when you
// expect very long input lines. For example, pass 1<<20 to allow lines up to
// roughly 1 MiB.
func WithMaxCapacity(maxCapacity int) Option {
	return func(o *options) {
		o.maxCapacity = maxCapacity
	}
}

// ScanStdin reads all lines from standard input until EOF and returns them as a
// slice of strings. It is a convenience wrapper around bufio.Scanner.
//
// Behavior and options:
//   - By default, the maximum token size (i.e., maximum line length) is
//     bufio.MaxScanTokenSize. You can override this using the WithMaxCapacity
//     option to support longer lines.
//   - This function accumulates all lines in memory before returning; consider
//     input size when using it.
//   - Errors from bufio.Scanner (including bufio.ErrTooLong when a line exceeds
//     the configured capacity) are not returned; scanning stops and any lines
//     collected up to that point are returned. If you expect very long lines,
//     pass a larger capacity via WithMaxCapacity.
//
// Example:
//
//	// Read with defaults
//	lines := utils.ScanStdin()
//
//	// Read allowing lines up to 1 MiB
//	lines := utils.ScanStdin(utils.WithMaxCapacity(1 << 20))
func ScanStdin(opts ...Option) ([]string, error) {
	options := &options{
		maxCapacity: bufio.MaxScanTokenSize,
	}
	for _, opt := range opts {
		opt(options)
	}

	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 0, options.maxCapacity)
	scanner.Buffer(buf, options.maxCapacity)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	if err != nil {
		return lines, err
	}
	return lines, nil
}
