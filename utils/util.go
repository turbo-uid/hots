package utils

import (
	"fmt"
	"strings"
	"time"
)

func GetHotCacheKey(flag string) string {
	now := time.Now()
	fmtTime := now.Format("200601021504")

	k := fmt.Sprintf("hot_%s_%s", flag, fmtTime)
	return k
}

func GetTimestamp() int64 {
	t := time.Now().Unix()

	return t
}

func RemoveChar(str, charToRemove string) string {
	var result strings.Builder
	for _, r := range str {
		if string(r) != charToRemove {
			result.WriteRune(r)
		}
	}
	return result.String()
}
