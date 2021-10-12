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

const (
	passwordTokens  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	accountIdTokens = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

func randomeString(tokens string, n int) string {
	var sb strings.Builder
	k := len(tokens)

	for i := 0; i < n; i++ {
		c := tokens[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func TestCryptPassword(t *testing.T) {
	plain := randomeString(passwordTokens, 8)
	hashed, err := EncryptPassword(plain)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	ok := VerifyPasswordHash(hashed, plain)
	require.True(t, ok)

	wrongPlain := randomeString(passwordTokens, 8)
	ng := VerifyPasswordHash(hashed, wrongPlain)
	require.False(t, ng)

	hashed2, err := EncryptPassword(plain)
	require.NoError(t, err)
	require.NotEmpty(t, hashed2)
	require.NotEqual(t, hashed, hashed2)
}

func TestVerifyPlainPassword(t *testing.T) {
	for i := 0; i < 8; i++ {
		plain := randomeString(passwordTokens, i)
		res := VerifyPlainPassword(plain)
		require.False(t, res)
	}

	for i := 8; i <= 255; i++ {
		plain := randomeString(passwordTokens, i)
		res := VerifyPlainPassword(plain)
		require.True(t, res)
	}

	plain := randomeString(passwordTokens, 256)
	res := VerifyPlainPassword(plain)
	require.False(t, res)
}

func TestVerifyAccountId(t *testing.T) {
	plain := randomeString(accountIdTokens, 0)
	res := VerifyAccountId(plain)
	require.False(t, res)

	plain = randomeString(accountIdTokens, 256)
	res = VerifyAccountId(plain)
	require.False(t, res)

	for i := 1; i <= 255; i++ {
		plain := randomeString(accountIdTokens, i)
		res := VerifyAccountId(plain)
		require.True(t, res)
	}
}

func TestVerifyAccountName(t *testing.T) {
	name := randomeString(accountIdTokens, 0)
	res := VerifyAccountName(name)
	require.False(t, res)

	name = randomeString(accountIdTokens, 256)
	res = VerifyAccountName(name)
	require.False(t, res)

	for i := 1; i <= 255; i++ {
		name := randomeString(accountIdTokens, i)
		res := VerifyAccountName(name)
		require.True(t, res)
	}
}
