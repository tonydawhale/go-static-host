package util

import (
	"math/rand"
	"time"
  )
  

func GenerateId(n int) string {
    const charset = "abcdefghijklmnopqrstuvwxyz" +
  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
	  b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}