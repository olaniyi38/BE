package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandom(t *testing.T) {
	random := RandomString(5)
	t.Log(random)
	require.Len(t, random, 5)
}
