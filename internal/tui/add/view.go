package add

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	primary = lipgloss.Color("#7C3AED")
	success = lipgloss.Color("#10B981")
	muted   = lipgloss.Color("#6B7280")
	white   = lipgloss.Color("#F9FAFB")
	yellow  = lipgloss.Color("#F59E0B")
	red     = lipgloss.Color("#EF4444")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(primary).
			Padding(0, 2)

	descStyle = lipgloss.NewStyle().
			Foreground(muted)

	queuedBadge = lipgloss.NewStyle().
			Foreground(success).
			Bold(true)

	versionStyle = lipgloss.NewStyle().
			Foreground(yellow)

	helpStyle = lipgloss.NewStyle().
			Foreground(muted)
)

const (
	visibleResults = 7  // max results shown at once
	leftPanelW     = 28 // fixed left panel width
)

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func (m Model) View() tea.View {
	str := "loading..."
	if m.width == 0 {
		return tea.NewView(str)
	}

	rightW := m.width - leftPanelW - 3
	innerH := m.height - 4

	// ── left panel ────────────────────────────────────────────────────────
	var left strings.Builder

	left.WriteString(lipgloss.NewStyle().
		Foreground(white).Bold(true).Render("Search") + "\n\n")

	left.WriteString(m.input.View() + "\n\n")

	if m.loading {
		left.WriteString(descStyle.Render("  searching...") + "\n")
	}

	if m.err != nil {
		left.WriteString(lipgloss.NewStyle().Foreground(red).
			Render("  "+m.err.Error()) + "\n")
	}

	if len(m.queued) > 0 {
		left.WriteString("\n" + queuedBadge.Render(
			fmt.Sprintf("Queued (%d)", len(m.queued))) + "\n")
		for i, q := range m.queued {
			left.WriteString(descStyle.Render(
				fmt.Sprintf("  %d. %s", i+1, truncate(q.Path, leftPanelW-5))) + "\n")
		}
	}

	leftPanel := lipgloss.NewStyle().
		Width(leftPanelW).
		Height(innerH).
		PaddingRight(2).
		Render(left.String())

	// ── right panel with scroll ───────────────────────────────────────────
	var right strings.Builder

	right.WriteString(lipgloss.NewStyle().
		Foreground(white).Bold(true).Render("Results") + "\n\n")

	if len(m.results) == 0 && !m.loading && m.input.Value() != "" {
		right.WriteString(descStyle.Render("no results found") + "\n")
	}

	// calculate scroll offset so cursor stays visible
	scrollOffset := 0
	if m.cursor >= visibleResults {
		scrollOffset = m.cursor - visibleResults + 1
	}

	total := len(m.results)
	for i := scrollOffset; i < total && i < scrollOffset+visibleResults; i++ {
		pkg := m.results[i]
		isSelected := m.cursor == i
		isQueued := alreadyQueued(m.queued, pkg.Path)

		queued := ""
		if isQueued {
			queued = queuedBadge.Render(" ✓")
		}

		if isSelected {
			row := lipgloss.NewStyle().
				Foreground(white).
				Background(primary).
				Bold(true).
				Width(rightW-2).
				Padding(0, 1).
				Render(fmt.Sprintf("▸ %s %s%s",
					truncate(pkg.Path, rightW-25),
					pkg.Version,
					queued,
				))
			desc := lipgloss.NewStyle().
				Foreground(white).
				Background(lipgloss.Color("#4C1D95")).
				Width(rightW-2).
				Padding(0, 1).
				Render(truncate(pkg.Description, rightW-4))
			right.WriteString(row + "\n" + desc + "\n")
		} else {
			name := lipgloss.NewStyle().Foreground(white).Bold(true).
				Render(truncate(pkg.Path, rightW-25))
			version := versionStyle.Render(" " + pkg.Version)
			desc := descStyle.Render("  " + truncate(pkg.Description, rightW-4))
			right.WriteString(fmt.Sprintf("  %s%s%s\n%s\n",
				name, version, queued, desc))
		}
	}

	// scroll indicator
	if total > visibleResults {
		shown := fmt.Sprintf("  %d-%d of %d",
			scrollOffset+1,
			min(scrollOffset+visibleResults, total),
			total,
		)
		right.WriteString("\n" + descStyle.Render(shown))
	}

	rightPanel := lipgloss.NewStyle().
		Width(rightW).
		Height(innerH).
		PaddingLeft(2).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(muted).
		Render(right.String())

	// ── assemble ──────────────────────────────────────────────────────────
	title := titleStyle.Render("  gostart add — package search  ")
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	help := helpStyle.Render(
		"j/k navigate  •  tab queue  •  enter install  •  ctrl+d install queued  •  ctrl+r unqueue  •  esc quit",
	)

	str = lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			title,
			body,
			help,
		))

	return tea.NewView(str)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
