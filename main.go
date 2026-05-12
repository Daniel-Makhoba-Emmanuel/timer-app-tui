package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// Unique datatype to identify our timer's tick messages when they arrive in Update()
type tickMsg time.Time

// Defines the application state, which includes a text input for the user to enter
// the number of seconds, a counter for the remaining seconds, and a boolean to track
// whether the timer is active or not.
type model struct {
	input       textinput.Model
	seconds     int
	timerActive bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "60"
	ti.SetWidth(20)
	ti.Focus()

	return model{
		input:       ti,
		seconds:     0,
		timerActive: false, //start in input mode, waiting for user to enter time
	}
}

// tick() creates a Command that tells the Runtime to wait 1 second, then send a tickMsg to Update()
// You can find out more about tea.Tick() here: https://pkg.go.dev/charm.land/bubbletea/v2#Tick
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Starts the first tick when the program starts
func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if !m.timerActive {
			switch {
			// Verifies that the user enters a valid number
			case msg.String() >= "0" && msg.String() <= "9":
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				return m, cmd

			// Handles time input by the user and starts the timer
			case msg.String() == "enter":
				val := m.input.Value()
				secs, _ := strconv.Atoi(val)
				m.seconds = secs

				// Sets timerActive to true to indicate that the timer is now running
				m.timerActive = true
				return m, tick()

			// Handles backspace key to allow user to correct their input
			case msg.String() == "backspace":
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
			return m, nil
		}

	case tickMsg:
		if m.seconds > 0 {
			m.seconds--
			// This tells the Runtime to schedule another 1-second wait.
			return m, tick()
		}

		//The tier has reached 0, so we set timerActive to false to indicate that the timer is no longer running
		m.timerActive = false
		m.input.SetValue("") // Clear the input field for the next timer
		return m, nil
	}

	return m, nil
}

func (m model) View() tea.View {
	var s string

	if !m.timerActive {
		// Room A: The Input Screen ⌨️
		s = "Set timer (seconds):\n\n"
		s += m.input.View() // This shows the actual text box
		s += "\n\n(press enter to start)"
	} else {
		// Room B: The Timer Screen 🕰️
		s = fmt.Sprintf("Time remaining: %d\n\n", m.seconds)
		if m.seconds == 0 {
			s = "Time's up!\n\n"
		}
		s += "(press q to quit)"
	}
	return tea.NewView(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
