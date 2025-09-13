package utils

import (
	"bufio"
	"os"
)

// ScanStdin reads all lines from standard input and returns them as a slice of strings.
func ScanStdin(maxCapacity int) []string {
	scanner := bufio.NewScanner(os.Stdin)
	// 新しいバッファの最大容量を定義 (ここでは 1MB = 1024 * 1024 bytes)
	// デフォルトは bufio.MaxScanTokenSize (64KB)
	// スキャナが使用するバッファ用のバイトスライスを作成
	// 長さ0、容量maxCapacityのスライスを用意
	buf := make([]byte, 0, maxCapacity)
	// 作成したバッファと最大容量をスキャナに設定
	scanner.Buffer(buf, maxCapacity)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
