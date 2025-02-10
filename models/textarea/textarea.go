package textarea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textarea textarea.Model
	callBack func(string) error
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":

			_ = m.callBack(m.textarea.Value())

			return m, tea.Quit
		case "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n(Press 'Ctrl+S' to submit, 'Esc' to exit)",
		m.textarea.View(),
	)
}

func NewTextAreaModel(callBack func(string) error) model {
	ta := textarea.New()
	ta.Placeholder = "Enter Task Description..."
	ta.Focus()
	ta.SetHeight(5)
	ta.SetWidth(40)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.CharLimit = 200
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{textarea: ta, callBack: callBack}
}
