package util

import (
	"crypto/rand"
	"encoding/hex"
)

func PrefixedID(prefix string) string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return prefix + "_" + hex.EncodeToString(buf)
}
