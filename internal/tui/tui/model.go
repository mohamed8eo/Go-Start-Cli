package tui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type step int

const (
	StepProjectName step = iota
	StepFramework
	StepDatabase
	StepSQL
	StepDone
)

type framework struct {
	Name        string
	description string
}

type database struct {
	Name        string
	description string
}

type sqlDriver struct {
	Name        string
	description string
}

type Model struct {
	// navigation
	currentStep step
	cursor      int
	width       int
	height      int

	ProjectName textinput.Model

	Frameworks        []framework
	SelectedFramework int

	Databases        []database
	SelectedDatabase int

	SqlDrivers  []sqlDriver
	SelectedSQL int
}

var defaultFrameworks = []framework{
	{Name: "Gin", description: "Fast HTTP web framework"},
	{Name: "Echo", description: "High performance, minimalist"},
	{Name: "Fiber", description: "Express-inspired, built on Fasthttp"},
	{Name: "Chi", description: "Lightweight, idiomatic router"},
	{Name: "None", description: "No framework, stdlib only"},
}

var defaultDatabases = []database{
	{Name: "PostgreSQL", description: "Advanced open source RDBMS"},
	{Name: "MySQL", description: "Popular open source RDBMS"},
	{Name: "SQLite", description: "Lightweight file-based DB"},
	{Name: "MongoDB", description: "NoSQL document database"},
	{Name: "None", description: "No database"},
}

var defaultSQLDrivers = []sqlDriver{
	{Name: "GORM", description: "Full-featured ORM for Go"},
	{Name: "sqlx", description: "Extensions on top of database/sql"},
	{Name: "sqlc", description: "Generate type-safe code from SQL"},
	{Name: "pgx", description: "PostgreSQL driver & toolkit"},
	{Name: "None", description: "Raw database/sql"},
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "my-project"
	ti.Focus()
	ti.CharLimit = 70
	ti.SetWidth(30)
	ti.SetVirtualCursor(false)

	return Model{
		currentStep:       StepProjectName,
		cursor:            0,
		ProjectName:       ti,
		Frameworks:        defaultFrameworks,
		SelectedFramework: 0,
		Databases:         defaultDatabases,
		SelectedDatabase:  0,
		SqlDrivers:        defaultSQLDrivers,
		SelectedSQL:       0,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
