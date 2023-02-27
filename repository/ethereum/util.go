package ethereum

import (
	"fmt"
	"strconv"
	"strings"
)

func numberToHex(number int64) string {
	return fmt.Sprintf("0x%s", strconv.FormatInt(number, 16))
}

func hexToNumber(hex string) int64 {
	n, _ := strconv.ParseInt(strings.ReplaceAll(hex, "0x", ""), 16, 64)
	return n
}
