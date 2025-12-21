package tarot

import (
	"strings"

	"github.com/amilleah/tarot/pkg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	gridCols = 6
)

type Tarot struct {
	deck   *pkg.Deck
	drawn  []*pkg.Card
	mouse  tea.MouseMsg
	width  int
	height int
}

func New() *Tarot {
	return &Tarot{
		deck: pkg.NewFullDeck(),
	}
}

func (t *Tarot) Init() tea.Cmd {
	return tea.ClearScreen
}

func (t *Tarot) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			print("\n")
			return t, tea.Quit
		case "r":
			t.reset()
		case "d", "enter", " ":
			t.draw()
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			if t.mouse.Type != tea.MouseLeft {
				t.mouse = msg
			}
		case tea.MouseRelease:
			if msg.X == t.mouse.X && msg.Y == t.mouse.Y {
				if t.mouse.Type == tea.MouseLeft {
					t.click(msg.X, msg.Y)
				}
			}
			t.mouse = msg
		}
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height
	}

	return t, nil
}

func (t *Tarot) View() string {
	grid := make([]*pkg.Deck, len(t.drawn))
	for i, card := range t.drawn {
		d := pkg.NewEmptyDeck()
		d.Add(card)
		grid[i] = d
	}

	cols := gridCols
	if t.width > 0 && cols > 0 {
		available := t.width
		if len(grid) > 0 {
			available -= pkg.Width
		}
		if available < pkg.Width {
			cols = 1
		} else {
			cols = available / pkg.Width
		}
	}
	if cols <= 0 {
		cols = 1
	}
	rows := (len(grid) + cols - 1) / cols

	var b strings.Builder
	for r := 0; r < rows; r++ {
		start := r * cols
		end := start + cols
		if end > len(grid) {
			end = len(grid)
		}
		row := make([]string, 0, end-start)
		for _, deck := range grid[start:end] {
			row = append(row, deck.View())
		}
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, row...))
		if r < rows-1 {
			b.WriteString("\n")
		}
	}

	gridView := b.String()
	if gridView == "" {
		return t.deck.View()
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, t.deck.View(), gridView)
}

func (t *Tarot) click(x, y int) {
	if t.deck == nil {
		return
	}
	if t.deck.Size() == 0 {
		t.reset()
		return
	}
	if ok, _ := t.deck.IsClicked(x, y); ok {
		t.draw()
	}
}

func (t *Tarot) draw() {
	if t.deck == nil || t.deck.Size() == 0 {
		if t.deck != nil && t.deck.Size() == 0 {
			t.reset()
		}
		return
	}
	card := t.deck.Pop()
	if card == nil {
		return
	}
	card.IsVisible = true
	t.drawn = append(t.drawn, card)
}

func (t *Tarot) reset() {
	t.deck = pkg.NewFullDeck()
	t.drawn = nil
}
