package entity

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomeString(n int) string {
	const tokens = "abcdefghijklmnopqrstuvwxyz0123456789"
	var sb strings.Builder
	k := len(tokens)

	for i := 0; i < n; i++ {
		c := tokens[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func TestCryptPassword(t *testing.T) {
	plain := randomeString(8)
	hashed, err := EncryptPassword(plain)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	ok := VerifyPasswordHash(hashed, plain)
	require.True(t, ok)

	wrongPlain := randomeString(8)
	ng := VerifyPasswordHash(hashed, wrongPlain)
	require.False(t, ng)

	hashed2, err := EncryptPassword(plain)
	require.NoError(t, err)
	require.NotEmpty(t, hashed2)
	require.NotEqual(t, hashed, hashed2)
}
