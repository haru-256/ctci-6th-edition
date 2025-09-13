// Package utils provides small helper utilities that are shared across the
// repository.
//
// The package currently exposes ScanStdin, a convenience function for reading
// all lines from standard input until EOF. It supports a functional option,
// WithMaxCapacity, to configure the maximum token size used by the underlying
// bufio.Scanner.
//
// Example:
//
//	// Read all lines with default capacity.
//	lines := utils.ScanStdin()
//
//	// Read with a larger buffer (e.g., 1 MiB) when lines can be long.
//	lines := utils.ScanStdin(utils.WithMaxCapacity(1 << 20))
//
// Be mindful that ScanStdin accumulates all lines in memory before returning,
// so callers should consider input size when using it.
package utils
