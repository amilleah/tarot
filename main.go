package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/amilleah/tarot/internal/tarot"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	p := tea.NewProgram(tarot.New(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
