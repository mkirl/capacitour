package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mkirl/capacitour/api"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load .env file
	config, err := api.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(config))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
} // <-- Added closing brace

type Icon struct {
	Type string `json:"type"`
	Val  string `json:"val"`
}

type Space struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Icon  Icon   `json:"icon"`
}

type SpacesResponse struct {
	Spaces []Space `json:"spaces"`
}

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	loading  bool
	config   *api.Config
	spaces   []Space
}

func initialModel(config *api.Config) tea.Model {
	m := model{}
	m.choices = []string{"Loading..."}
	m.cursor = 0
	m.selected = make(map[int]struct{})
	m.loading = true
	m.config = config
	m.spaces = []Space{}
	return &m
}

func (m *model) Init() tea.Cmd {
	spacesData, err := api.FetchSpacesData(m.config)
	if err != nil {
		fmt.Printf("Error fetching spaces data: %v", err)
		return nil
	}

	var spacesResponse SpacesResponse
	err = json.Unmarshal(spacesData, &spacesResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling spaces data: %v", err)
		return nil
	}

	m.spaces = append(m.spaces, spacesResponse.Spaces...)
	fmt.Println(m.spaces[0].Title)

	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m *model) View() string {
	s := "Welcome to Capacitour\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}
