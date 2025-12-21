package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDeck(t *testing.T) {
	deck := NewFullDeck()

	expected := make([]*Card, 0, len(values)*len(suits)+len(majorArtKeys))
	for i := range values {
		for j := range suits {
			expected = append(expected, NewCard(i, j))
		}
	}
	for i := range majorArtKeys {
		expected = append(expected, NewCard(i, majorSuit))
	}

	require.ElementsMatch(t, expected, deck.cards)
}

func TestShuffle(t *testing.T) {
	deck := NewDeck([]*Card{NewCard(0, 0), NewCard(1, 0)})
	deck.Shuffle()

	require.ElementsMatch(t, []*Card{NewCard(0, 0), NewCard(1, 0)}, deck.cards)
}
