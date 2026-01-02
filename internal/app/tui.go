package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// State constants
const (
	StateInit = iota
	StateCatchUpPrompt
	StateDayEntry
	StateReview
	StateSaving
	StateDone
)

type Model struct {
	state  int
	cfg    *Config      // Direct reference to app.Config
	client *SheetClient // Direct reference to app.SheetClient

	missingDates []string
	currentIndex int
	answers      map[string]interface{}
	habitIndex   int

	noteInput textinput.Model
	err       error
	statusMsg string
}

func NewModel(cfg *Config, client *SheetClient) Model {
	ti := textinput.New()
	ti.Placeholder = "Optional notes..."
	ti.CharLimit = 200

	return Model{
		state:     StateInit,
		cfg:       cfg,
		client:    client,
		answers:   make(map[string]interface{}),
		noteInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		headers, err := m.client.FetchHeaders()
		if err != nil {
			return errMsg(err)
		}
		_ = headers

		dates, err := m.client.FetchExistingDates()
		if err != nil {
			return errMsg(err)
		}

		missing, err := CalculateMissingDays(dates, m.cfg.Timezone)
		if err != nil {
			return errMsg(err)
		}

		return missingDatesMsg(missing)
	}
}

// Internal message types (unexported is fine)
type errMsg error
type missingDatesMsg []string
type savedMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, tea.Quit

	case missingDatesMsg:
		m.missingDates = msg
		if len(m.missingDates) == 0 {
			m.statusMsg = "All caught up!"
			m.state = StateDone
			return m, tea.Quit
		}

		if len(m.missingDates) > 1 {
			m.state = StateCatchUpPrompt
		} else {
			m.state = StateDayEntry
			m.habitIndex = 0
		}

	case savedMsg:
		m.currentIndex++
		if m.currentIndex >= len(m.missingDates) {
			m.state = StateDone
			return m, tea.Quit
		}
		m.habitIndex = 0
		m.answers = make(map[string]interface{})
		m.noteInput.Reset()
		m.state = StateDayEntry

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		switch m.state {
		case StateCatchUpPrompt:
			if msg.String() == "y" || msg.String() == "Y" || msg.String() == "enter" {
				m.state = StateDayEntry
				m.habitIndex = 0
			} else if msg.String() == "n" || msg.String() == "N" {
				m.currentIndex = len(m.missingDates) - 1
				m.state = StateDayEntry
				m.habitIndex = 0
			}

		case StateDayEntry:
			if m.habitIndex < len(m.cfg.Habits) {
				currentHabit := m.cfg.Habits[m.habitIndex]
				switch msg.String() {
				case "y", "Y":
					m.answers[currentHabit.Column] = true
					m.habitIndex++
				case "n", "N":
					m.answers[currentHabit.Column] = false
					m.habitIndex++
				case "enter":
					m.answers[currentHabit.Column] = currentHabit.Default
					m.habitIndex++
				}
			} else {
				var cmd tea.Cmd
				m.noteInput, cmd = m.noteInput.Update(msg)
				if msg.String() == "enter" {
					m.state = StateReview
				}
				return m, cmd
			}

		case StateReview:
			if msg.String() == "enter" {
				m.state = StateSaving
				return m, m.saveRow
			}
		}
	}

	return m, nil
}

func (m Model) saveRow() tea.Msg {
	currentDate := m.missingDates[m.currentIndex]
	m.answers["date"] = currentDate
	m.answers["timestamp_submitted"] = time.Now().Format(time.RFC3339)
	m.answers["timezone"] = m.cfg.Timezone
	m.answers["notes"] = m.noteInput.Value()

	if err := m.client.AppendRow(m.answers); err != nil {
		return errMsg(err)
	}
	return savedMsg{}
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	switch m.state {
	case StateInit:
		return "Connecting to Google Sheets...\n"
	case StateCatchUpPrompt:
		return fmt.Sprintf("Missing %d days (%s to %s). Start catch-up? (Y/n)\n",
			len(m.missingDates), m.missingDates[0], m.missingDates[len(m.missingDates)-1])
	case StateDayEntry:
		currentDate := m.missingDates[m.currentIndex]
		if m.habitIndex < len(m.cfg.Habits) {
			h := m.cfg.Habits[m.habitIndex]
			def := "n"
			if h.Default {
				def = "y"
			}
			return fmt.Sprintf("Date: %s\n\n%s? [y/n] (default: %s)\n", currentDate, h.Label, def)
		} else {
			return fmt.Sprintf("Date: %s\n\nNotes: %s\n(Enter to finish)", currentDate, m.noteInput.View())
		}
	case StateReview:
		s := fmt.Sprintf("Review for %s:\n", m.missingDates[m.currentIndex])
		for _, h := range m.cfg.Habits {
			s += fmt.Sprintf("- %s: %v\n", h.Label, m.answers[h.Column])
		}
		s += fmt.Sprintf("Notes: %s\n\nSave to Sheet? (Enter to save, Ctrl+C to quit)", m.noteInput.Value())
		return s
	case StateSaving:
		return "Saving..."
	case StateDone:
		return "Done!\n"
	}
	return ""
}
