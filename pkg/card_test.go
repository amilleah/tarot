package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArtLookup(t *testing.T) {
	card := &Card{Value: 0, Suit: majorSuit}
	art, ok := artFor(card)

	require.True(t, ok)
	require.NotEmpty(t, art)
}
