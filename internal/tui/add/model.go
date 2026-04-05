package add

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type Package struct {
	Path        string
	Description string
	Version     string
}

type state int

const (
	StateSearching state = iota
	StateInstalling
	StateDone
)

type Model struct {
	input   textinput.Model
	results []Package
	queued  []Package
	cursor  int
	state   state
	err     error
	loading bool
	width   int
	height  int
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "search packages... (e.g. argon2)"
	ti.Focus()
	ti.CharLimit = 100
	ti.SetWidth(50)
	ti.SetVirtualCursor(false)

	return Model{
		input:   ti,
		results: []Package{},
		queued:  []Package{},
		cursor:  0,
		state:   StateSearching,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
