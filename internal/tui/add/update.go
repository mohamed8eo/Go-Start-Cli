package add

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case searchResultMsg:
		m.loading = false
		m.cursor = 0
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
			m.results = msg.packages
		}
		return m, nil

	case triggerSearchMsg:
		if msg.query == m.input.Value() {
			m.loading = true
			return m, searchPackages(msg.query)
		}
		return m, nil

	case installDoneMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = StateSearching
		} else {
			m.state = StateDone
			return m, tea.Quit
		}
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit

		case "j":
			if len(m.results) > 0 {
				if m.cursor < len(m.results)-1 {
					m.cursor++
				}
				return m, nil
			}

		case "k":
			if len(m.results) > 0 {
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			}

		case "tab":
			if len(m.results) > 0 {
				pkg := m.results[m.cursor]
				if !alreadyQueued(m.queued, pkg.Path) {
					m.queued = append(m.queued, pkg)
				}
				return m, nil
			}

		case "enter":
			if len(m.results) > 0 {
				m.state = StateInstalling
				return m, installPackage(m.results[m.cursor])
			}

		case "ctrl+d":
			if len(m.queued) > 0 {
				m.state = StateInstalling
				return m, installPackages(m.queued)
			}

		case "ctrl+r":
			if len(m.queued) > 0 {
				m.queued = m.queued[:len(m.queued)-1]
			}
			return m, nil
		}
	}

	prevQuery := m.input.Value()
	m.input, cmd = m.input.Update(msg)
	newQuery := m.input.Value()

	if newQuery != prevQuery {
		m.results = []Package{}
		m.cursor = 0
		m.loading = true
		return m, tea.Batch(cmd, triggerSearch(newQuery))
	}

	return m, cmd
}

func alreadyQueued(queued []Package, path string) bool {
	for _, q := range queued {
		if q.Path == path {
			return true
		}
	}
	return false
}

type installDoneMsg struct{ err error }

func installPackage(pkg Package) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("go", "get", pkg.Path+"@latest")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return installDoneMsg{err: err}
		}
		return installDoneMsg{}
	}
}

func installPackages(pkgs []Package) tea.Cmd {
	return func() tea.Msg {
		for _, pkg := range pkgs {
			fmt.Println("Installing", pkg.Path+"@latest")
			cmd := exec.Command("go", "get", pkg.Path+"@latest")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return installDoneMsg{err: err}
			}
		}
		return installDoneMsg{}
	}
}

func triggerSearch(query string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(350 * time.Millisecond)
		return triggerSearchMsg{query: query}
	}
}
