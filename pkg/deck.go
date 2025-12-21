package pkg

import (
	"math/rand"
	"strconv"
	"strings"
)

type Deck struct {
	cards      []*Card
	isExpanded bool
}

func NewDeck(cards []*Card) *Deck {
	return &Deck{cards: cards}
}

func NewFullDeck() *Deck {
	cards := make([]*Card, 0, len(values)*len(suits)+len(majorArtKeys))
	for i := range values {
		for j := range suits {
			cards = append(cards, NewCard(i, j))
		}
	}
	for i := range majorArtKeys {
		cards = append(cards, NewCard(i, majorSuit))
	}

	deck := NewDeck(cards)
	deck.Shuffle()

	return deck
}

func NewEmptyDeck() *Deck {
	return NewDeck(make([]*Card, 0))
}

func (d *Deck) Shuffle() {
	rand.Shuffle(d.Size(), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Expand() {
	d.isExpanded = true
}

func (d *Deck) View() string {
	if d == nil {
		return strings.Repeat(" ", Width)
	}

	// Outline
	if d.Size() == 0 {
		return viewCard(" ", "")
	}

	// Expanded cards
	if d.isExpanded {
		var view string
		for i := 0; i < d.Size()-1; i++ {
			view += strings.Split(d.cards[i].View(), "\n")[0] + "\n"
		}
		return view + d.cards[d.Size()-1].View()
	}

	// Top card only
	return d.cards[d.Size()-1].View()
}

func (d *Deck) IsClicked(x, y int) (bool, int) {
	if d == nil {
		return false, 0
	}

	if d.Size() == 0 {
		return x >= 0 && x < Width && y >= 0 && y < Height, 0
	}

	if d.isExpanded {
		for i := d.Size() - 1; i >= 0; i-- {
			if x >= 0 && x < Width && y >= i && y < i+Height {
				return true, i
			}
		}
		return false, 0
	}

	return x >= 0 && x < Width && y >= 0 && y < Height, 0
}

func (d *Deck) Add(cards ...*Card) {
	d.cards = append(d.cards, cards...)
}

func (d *Deck) Top() *Card {
	return d.Get(d.Size() - 1)
}

func (d *Deck) Bottom() *Card {
	return d.Get(0)
}

func (d *Deck) Get(idx int) *Card {
	return d.cards[idx]
}

func (d *Deck) GetFrom(idx int) []*Card {
	return d.cards[idx:]
}

func (d *Deck) Pop() *Card {
	if len(d.cards) > 0 {
		return d.PopFrom(d.Size() - 1)[0]
	}

	return nil
}

func (d *Deck) PopFrom(idx int) []*Card {
	cards := d.cards[idx:]
	d.cards = d.cards[:idx]
	return cards
}

func (d *Deck) Size() int {
	return len(d.cards)
}

var valueIndexByLabel = map[string]int{
	"1":  0,
	"2":  1,
	"3":  2,
	"4":  3,
	"5":  4,
	"6":  5,
	"7":  6,
	"8":  7,
	"9":  8,
	"10": 9,
	"P":  10,
	"N":  11,
	"Q":  12,
	"K":  13,
}

var valueParseOrder = []string{"10", "1", "2", "3", "4", "5", "6", "7", "8", "9", "P", "N", "Q", "K"}

func parseValueIndex(shorthand string) (int, bool) {
	for _, label := range valueParseOrder {
		if strings.HasPrefix(shorthand, label) {
			idx, ok := valueIndexByLabel[label]
			return idx, ok
		}
	}
	return 0, false
}

func parseMajorIndex(shorthand string) (int, bool) {
	if strings.HasPrefix(shorthand, "RWS_Tarot_") {
		parts := strings.Split(shorthand, "_")
		if len(parts) >= 3 {
			if n, err := strconv.Atoi(parts[2]); err == nil {
				return n, true
			}
		}
	}
	if strings.HasPrefix(shorthand, "M") && len(shorthand) >= 3 {
		if n, err := strconv.Atoi(shorthand[1:3]); err == nil {
			return n, true
		}
	}
	return 0, false
}

// TestDeck is a helper function to simplify testing.
func TestDeck(shorthands ...string) *Deck {
	cards := make([]*Card, len(shorthands))
	for i, shorthand := range shorthands {
		cards[i] = testCard(shorthand)
	}
	return &Deck{cards: cards}
}

func testCard(shorthand string) *Card {
	card := &Card{}
	shorthand = strings.TrimSuffix(shorthand, "?")

	// Major shorthand: M00 or RWS_Tarot_00_*
	if strings.HasPrefix(shorthand, "M") || strings.HasPrefix(shorthand, "RWS_Tarot_") {
		if idx, ok := parseMajorIndex(shorthand); ok {
			card.Value = idx
			card.Suit = majorSuit
			return card
		}
	}

	if idx, ok := parseValueIndex(shorthand); ok {
		card.Value = idx
	}
	for i, suit := range suits {
		if strings.Contains(shorthand, suit) {
			card.Suit = i
			break
		}
	}

	return card
}
