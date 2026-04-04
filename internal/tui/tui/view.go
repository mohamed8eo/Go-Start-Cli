package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	primary   = lipgloss.Color("#7C3AED")
	secondary = lipgloss.Color("#A78BFA")
	muted     = lipgloss.Color("#6B7280")
	success   = lipgloss.Color("#10B981")
	white     = lipgloss.Color("#F9FAFB")

	appStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primary)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(primary).
			Padding(0, 2).
			MarginBottom(1)

	stepActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primary).
			PaddingLeft(1).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(primary)

	stepInactiveStyle = lipgloss.NewStyle().
				Foreground(muted).
				PaddingLeft(2)

	stepDoneStyle = lipgloss.NewStyle().
			Foreground(success).
			PaddingLeft(2)

	cursorStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(success).
			Bold(true)

	unselectedStyle = lipgloss.NewStyle().
			Foreground(white)

	descStyle = lipgloss.NewStyle().
			Foreground(muted).
			Italic(true)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(secondary).
			Bold(true).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(muted)

	doneStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(success)

	summaryKeyStyle = lipgloss.NewStyle().
			Foreground(secondary).
			Bold(true).
			Width(12)

	summaryValStyle = lipgloss.NewStyle().
			Foreground(white)
)

func (m Model) View() tea.View {
	header := titleStyle.Render("  gostart — Go project scaffolder  ")
	footer := helpStyle.Render(renderHelp(m))

	// Layout sizes
	sidebarWidth := 20
	contextWidth := m.width / 3
	contentWidth := m.width - sidebarWidth - contextWidth - 6
	if contentWidth < 0 {
		contentWidth = 0
	}

	// Panels
	sidebar := renderSidebar(m)

	content := lipgloss.NewStyle().
		Width(contentWidth).
		PaddingLeft(2).
		Render(renderContent(m))

	context := lipgloss.NewStyle().
		Width(contextWidth).
		BorderLeft(true).
		BorderForeground(muted).
		PaddingLeft(2).
		Render(renderContext(m))

	// Body
	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		content,
		context,
	)

	// Vertical sizing
	headerH := lipgloss.Height(header)
	footerH := lipgloss.Height(footer)
	bodyH := m.height - headerH - footerH - 2
	if bodyH < 0 {
		bodyH = 0
	}

	body = lipgloss.NewStyle().
		Height(bodyH).
		Render(body)

	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
		footer,
	)

	full := appStyle.
		Width(m.width).
		Height(m.height).
		Render(layout)

	return tea.NewView(full)
}

func renderContext(m Model) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(" Project Structure "),
		descStyle.Render(buildTree(m)),
	)
}

func renderSidebar(m Model) string {
	steps := []struct {
		label string
		s     step
	}{
		{"Project Name", StepProjectName},
		{"Framework", StepFramework},
		{"Database", StepDatabase},
		{"SQL Driver", StepSQL},
	}

	var sb strings.Builder
	for _, item := range steps {
		var line string
		switch {
		case item.s == m.currentStep:
			line = stepActiveStyle.Render("● " + item.label)
		case item.s < m.currentStep:
			line = stepDoneStyle.Render("✓ " + item.label)
		default:
			line = stepInactiveStyle.Render("○ " + item.label)
		}
		sb.WriteString(line + "\n")
	}

	return lipgloss.NewStyle().
		Width(18).
		BorderRight(true).
		BorderForeground(muted).
		PaddingRight(2).
		Render(sb.String())
}

func renderContent(m Model) string {
	switch m.currentStep {
	case StepProjectName:
		return renderProjectName(m)
	case StepFramework:
		return renderList(m, "Choose a Framework", frameworkItems(m))
	case StepDatabase:
		return renderList(m, "Choose a Database", databaseItems(m))
	case StepSQL:
		return renderList(m, "Choose a SQL Driver", sqlItems(m))
	case StepDone:
		return doneView(m)
	}
	return ""
}

func renderProjectName(m Model) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		inputLabelStyle.Render("Project Name"),
		m.ProjectName.View(),
		descStyle.Render("This will be used as the module name"),
	)
}

type listItem struct {
	name        string
	description string
	selected    bool
}

func renderList(m Model, title string, items []listItem) string {
	var sb strings.Builder
	sb.WriteString(inputLabelStyle.Render(title) + "\n\n")

	for i, item := range items {
		cursor := "  "
		if m.cursor == i {
			cursor = cursorStyle.Render("▸ ")
		}

		var nameRendered string
		if item.selected {
			nameRendered = selectedStyle.Render("◉ " + item.name)
		} else if m.cursor == i {
			nameRendered = unselectedStyle.Render("○ " + item.name)
		} else {
			nameRendered = lipgloss.NewStyle().Foreground(muted).Render("○ " + item.name)
		}

		desc := descStyle.Render("  " + item.description)

		sb.WriteString(fmt.Sprintf("%s%s\n%s\n", cursor, nameRendered, desc))
	}

	return sb.String()
}

func renderHelp(m Model) string {
	if m.currentStep == StepProjectName {
		return "type project name • enter to continue • ctrl+q to quit"
	}
	return "j/k move • space select • enter confirm • ↑/↓ switch • ctrl+q quit"
}

// Done screen
func doneView(m Model) string {
	var sb strings.Builder

	sb.WriteString(doneStyle.Render("✓ Project ready to scaffold!") + "\n\n")

	rows := []struct{ k, v string }{
		{"Name", m.ProjectName.Value()},
		{"Framework", m.Frameworks[m.SelectedFramework].Name},
		{"Database", m.Databases[m.SelectedDatabase].Name},
		{"SQL Driver", m.SqlDrivers[m.SelectedSQL].Name},
	}

	for _, r := range rows {
		sb.WriteString(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				summaryKeyStyle.Render(r.k+":"),
				summaryValStyle.Render(r.v),
			) + "\n",
		)
	}

	return sb.String()
}

func frameworkItems(m Model) []listItem {
	items := make([]listItem, len(m.Frameworks))
	for i, f := range m.Frameworks {
		items[i] = listItem{
			name:        f.Name,
			description: f.description,
			selected:    m.SelectedFramework == i,
		}
	}
	return items
}

func databaseItems(m Model) []listItem {
	items := make([]listItem, len(m.Databases))
	for i, d := range m.Databases {
		items[i] = listItem{
			name:        d.Name,
			description: d.description,
			selected:    m.SelectedDatabase == i,
		}
	}
	return items
}

func sqlItems(m Model) []listItem {
	items := make([]listItem, len(m.SqlDrivers))
	for i, s := range m.SqlDrivers {
		items[i] = listItem{
			name:        s.Name,
			description: s.description,
			selected:    m.SelectedSQL == i,
		}
	}
	return items
}

func buildTree(m Model) string {
	name := m.ProjectName.Value()
	if name == "" {
		name = "myproject"
	}

	var sb strings.Builder

	sb.WriteString(name + "/\n")
	sb.WriteString("├── cmd/\n")
	sb.WriteString("│   ├── api/\n")
	sb.WriteString("│   │   └── main.go\n")
	sb.WriteString("│   └── worker/\n")
	sb.WriteString("│       └── main.go\n")
	sb.WriteString("├── internal/\n")
	sb.WriteString("│   ├── auth/\n")
	sb.WriteString("│   ├── storage/\n")
	sb.WriteString("│   └── transport/\n")
	sb.WriteString("├── pkg/\n")
	sb.WriteString("│   ├── logger/\n")
	sb.WriteString("│   └── crypto/\n")
	sb.WriteString("├── api/\n")
	sb.WriteString("│   └── openapi.yaml\n")
	sb.WriteString("├── config/\n")
	sb.WriteString("│   └── config.yaml\n")
	sb.WriteString("├── scripts/\n")
	sb.WriteString("│   └── deploy.sh\n")

	// Dynamic framework folder
	if m.currentStep >= StepFramework {
		fw := strings.ToLower(m.Frameworks[m.SelectedFramework].Name)
		switch fw {
		case "gin":
			sb.WriteString("│   └── server/\n")
		case "echo":
			sb.WriteString("│   └── handlers/\n")
		case "fiber":
			sb.WriteString("│   └── fiber/\n")
		case "chi":
			sb.WriteString("│   └── routes/\n")
		}
	}

	// Dynamic database files
	if m.currentStep >= StepDatabase {
		db := strings.ToLower(m.Databases[m.SelectedDatabase].Name)
		switch db {
		case "postgresql":
			sb.WriteString("├── db/\n")
			sb.WriteString("│   └── postgres.go\n")
		case "mysql":
			sb.WriteString("├── db/\n")
			sb.WriteString("│   └── mysql.go\n")
		case "sqlite":
			sb.WriteString("├── db/\n")
			sb.WriteString("│   └── sqlite.go\n")
		case "mongodb":
			sb.WriteString("├── db/\n")
			sb.WriteString("│   └── mongodb.go\n")
		}
	}

	// Dynamic SQL driver folder
	if m.currentStep >= StepSQL {
		driver := strings.ToLower(m.SqlDrivers[m.SelectedSQL].Name)
		if driver != "none" {
			sb.WriteString("├── migrations/\n")
		}
	}

	// Root files
	sb.WriteString("├── go.mod\n")
	sb.WriteString("├── go.sum\n")
	sb.WriteString("└── README.md\n")

	return sb.String()
}
