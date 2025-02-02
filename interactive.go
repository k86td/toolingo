package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type state struct {
	modified     []Transaction
	list         []Transaction
	lookingIndex int
	typing       bool
}

func (m state) Init() tea.Cmd {
	return nil
}

func NewState(list []Transaction) state {
	var mod []Transaction
	return state{
		modified:     mod,
		list:         list,
		lookingIndex: 0,
		typing:       false,
	}
}

func (m state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.typing {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "i":
				m.typing = true
				return m, nil
			}
		} else {
			tran := &m.list[m.lookingIndex]
			switch msg.String() {
			case "esc":
				m.typing = false
				return m, nil
			case "backspace":
				if tran.Destination != "" {
					tran.Destination = tran.Destination[:len(tran.Destination)-1]
				}
				return m, nil
			default:
				tran.Destination = tran.Destination + msg.String()
				return m, nil
			}
		}
	}

	return m, nil
}

func (m state) View() string {
	builder := strings.Builder{}
	t := m.list[m.lookingIndex]

	if m.typing {
		builder.WriteString("[EDIT]")
	} else {
		builder.WriteString("[VIEW]")
	}

	builder.WriteString(fmt.Sprintf(" %v/%v ('i': enter typing mode, 'esc': exit typing mode, 'q': quit, 'm': mark as modified)\n", m.lookingIndex, len(m.list)))
	builder.WriteString(t.String())

	return builder.String()
}
