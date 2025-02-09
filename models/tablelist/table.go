package tablelist

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	selectedRow  = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(true)
	defaultStyle = lipgloss.NewStyle()
)

type SelectFunc func(int) error

type model struct {
	table      table.Model
	selectFunc SelectFunc
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "up":
			m.table.MoveUp(1)
		case "down":
			m.table.MoveDown(1)
		case "enter":
			selectedRow := m.table.Cursor()

			err := m.selectFunc(selectedRow)

			if err != nil {
				return m, tea.Quit
			}

			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return m.table.View() + "\nPress 'q' to quit."
}

// Function to create the table model
func NewTableModel(rows []table.Row, height int, selectFunc SelectFunc) model {
	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Description", Width: 30},
		{Title: "Status", Width: 10},
		{Title: "Creation Date", Width: 15},
		{Title: "Update Date", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true), // Enable keyboard navigation
		table.WithHeight(height+1),
	)
	t.SetStyles(table.Styles{
		Header:   headerStyle,
		Selected: selectedRow,
		Cell:     defaultStyle,
	})

	return model{table: t, selectFunc: selectFunc}
}
