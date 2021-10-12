package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUuid(t *testing.T) {
	uuids := make([]string, 1000)
	for i := 0; i < len(uuids); i++ {
		uuid, err := CreateUuid()
		require.NoError(t, err)
		uuids[i] = uuid
	}

	for i := 0; i < len(uuids)-1; i++ {
		require.NotEqual(t, uuids[i], uuids[i+1])
	}
}
