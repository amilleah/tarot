package pkg

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
	reflowansi "github.com/muesli/reflow/ansi"
)

var (
	values       = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "P", "N", "Q", "K"}
	suits        = []string{"Cups", "Pents", "Swords", "Wands"}
	tarotArt     = loadArt(ansiFS, "ansi")
	majorArtKeys = loadMajorKeys(ansiFS, "ansi")
)

var suitSymbols = map[string]string{
	"Cups":   "♥",
	"Pents":  "♦",
	"Swords": "♠",
	"Wands":  "♣",
}

const (
	Width     = 6
	Height    = 5
	majorSuit = 4
)

//go:embed ansi
var ansiFS embed.FS

type Card struct {
	Value     int
	Suit      int
	IsVisible bool
}

func NewCard(value, suit int) *Card {
	return &Card{
		Value:     value,
		Suit:      suit,
		IsVisible: false,
	}
}

func (c *Card) View() string {
	if !c.IsVisible {
		return viewCard("╱", "")
	}
	if art, ok := artFor(c); ok {
		return viewArtCard(art, c.String())
	}
	return viewCard(" ", c.String())
}

func (c *Card) String() string {
	if c.Suit == majorSuit {
		return romanNumeral(c.Value)
	}
	if c.Value < 0 || c.Value >= len(values) || c.Suit < 0 || c.Suit >= len(suits) {
		return ""
	}
	symbol := suitSymbols[suits[c.Suit]]
	return values[c.Value] + symbol
}

func romanNumeral(n int) string {
	if n == 0 {
		return "0"
	}
	lookup := []struct {
		Value  int
		Symbol string
	}{
		// {1000, "M"},
		// {900, "CM"},
		// {500, "D"},
		// {400, "CD"},
		// {100, "C"},
		// {90, "XC"},
		// {50, "L"},
		// {40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}
	var b strings.Builder
	for _, item := range lookup {
		for n >= item.Value {
			b.WriteString(item.Symbol)
			n -= item.Value
		}
	}
	return b.String()
}

func artFor(c *Card) ([]string, bool) {
	key := artKey(c)
	if key == "" {
		return nil, false
	}
	lines, ok := tarotArt[key]
	return lines, ok
}

func artKey(c *Card) string {
	if c.Suit == majorSuit {
		if key, ok := majorArtKeys[c.Value]; ok {
			return key
		}
		return ""
	}
	if c.Value < 0 || c.Value >= len(values) {
		return ""
	}
	if c.Suit < 0 || c.Suit >= len(suits) {
		return ""
	}
	return fmt.Sprintf("%s%02d", suits[c.Suit], c.Value+1)
}

func viewArtCard(lines []string, shorthand string) string {
	innerW := Width - 2
	innerH := Height - 2

	shorthand = trimToWidth(shorthand, innerW)
	padding := strings.Repeat("─", innerW-reflowansi.PrintableRuneWidth(shorthand))

	inner := make([]string, innerH)
	for i := range inner {
		if i < len(lines) {
			inner[i] = padRight(lines[i], innerW)
		} else {
			inner[i] = strings.Repeat(" ", innerW)
		}
	}

	var b strings.Builder
	b.WriteString("╭" + shorthand + padding + "╮\n")
	for i := 0; i < innerH; i++ {
		b.WriteString("│")
		b.WriteString(inner[i])
		b.WriteString("│\n")
	}
	b.WriteString("╰" + padding + shorthand + "╯")
	return b.String()
}

func viewCard(design string, shorthand string) string {
	innerW := Width - 2
	innerH := Height - 2

	shorthand = trimToWidth(shorthand, innerW)
	padding := strings.Repeat("─", innerW-reflowansi.PrintableRuneWidth(shorthand))

	var b strings.Builder
	b.WriteString("╭" + shorthand + padding + "╮\n")
	for i := 0; i < innerH; i++ {
		b.WriteString("│")
		b.WriteString(strings.Repeat(design, innerW))
		b.WriteString("│\n")
	}
	b.WriteString("╰" + padding + shorthand + "╯")
	return b.String()
}

func trimToWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if reflowansi.PrintableRuneWidth(s) <= width {
		return s
	}
	var b strings.Builder
	w := 0
	for _, r := range s {
		rw := runewidth.RuneWidth(r)
		if w+rw > width {
			break
		}
		b.WriteRune(r)
		w += rw
	}
	return b.String()
}

func padRight(s string, width int) string {
	pad := width - reflowansi.PrintableRuneWidth(s)
	if pad <= 0 {
		return s
	}
	return s + strings.Repeat(" ", pad)
}

func loadArt(fsys fs.FS, dir string) map[string][]string {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		return map[string][]string{}
	}
	pattern := path.Join(".", "*.ans")
	names, err := fs.Glob(sub, pattern)
	if err != nil || len(names) == 0 {
		return map[string][]string{}
	}
	sort.Strings(names)

	out := make(map[string][]string, len(names))
	for _, name := range names {
		data, readErr := fs.ReadFile(sub, name)
		if readErr != nil {
			continue
		}
		lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
		base := path.Base(name)
		key := strings.TrimSuffix(base, path.Ext(base))
		out[key] = lines
	}
	return out
}

func loadMajorKeys(fsys fs.FS, dir string) map[int]string {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		return map[int]string{}
	}
	pattern := path.Join(".", "RWS_Tarot_*.ans")
	names, err := fs.Glob(sub, pattern)
	if err != nil || len(names) == 0 {
		return map[int]string{}
	}

	out := make(map[int]string)
	for _, name := range names {
		base := path.Base(name)
		stem := strings.TrimSuffix(base, path.Ext(base))
		parts := strings.Split(stem, "_")
		if len(parts) < 3 {
			continue
		}
		idx, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}
		out[idx] = stem
	}
	return out
}
