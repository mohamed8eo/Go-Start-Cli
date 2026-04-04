package tui

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit

		case "up":
			switch m.currentStep {
			case StepProjectName:
				m.currentStep = StepSQL
			case StepFramework:
				m.currentStep = StepProjectName
			case StepDatabase:
				m.currentStep = StepFramework
			case StepSQL:
				m.currentStep = StepDatabase
			}
			m.cursor = 0

		case "down":
			switch m.currentStep {
			case StepProjectName:
				m.currentStep = StepFramework
			case StepFramework:
				m.currentStep = StepDatabase
			case StepDatabase:
				m.currentStep = StepSQL
			case StepSQL:
				m.currentStep = StepProjectName
			}
			m.cursor = 0

		case "j":
			if m.currentStep != StepProjectName && m.listLen() > 0 {
				m.cursor = (m.cursor + 1) % m.listLen()
			}

		case "k":
			if m.currentStep != StepProjectName && m.listLen() > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = m.listLen() - 1
				}
			}

		case "space", "y":
			switch m.currentStep {
			case StepFramework:
				m.SelectedFramework = m.cursor
			case StepDatabase:
				m.SelectedDatabase = m.cursor
			case StepSQL:
				m.SelectedSQL = m.cursor
			}

		case "enter":
			switch m.currentStep {
			case StepProjectName:
				if m.ProjectName.Value() != "" {
					m.currentStep = StepFramework
					m.cursor = 0
				}
			case StepFramework:
				m.currentStep = StepDatabase
				m.cursor = 0
			case StepDatabase:
				m.currentStep = StepSQL
				m.cursor = 0
			case StepSQL:
				m.currentStep = StepDone
			case StepDone:
				return m, tea.Quit
			}
		}
	}

	if m.currentStep == StepProjectName {
		m.ProjectName, cmd = m.ProjectName.Update(msg)
	}

	return m, cmd
}

func (m Model) listLen() int {
	switch m.currentStep {
	case StepFramework:
		return len(m.Frameworks)
	case StepDatabase:
		return len(m.Databases)
	case StepSQL:
		return len(m.SqlDrivers)
	default:
		return 0
	}
}
