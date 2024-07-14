package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/mkirl/capacitour/api"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func main() {
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
}

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
	list     list.Model
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func initialModel(config *api.Config) tea.Model {
	m := model{}
	m.choices = []string{}
	m.cursor = 0
	m.selected = make(map[int]struct{})
	m.loading = true
	m.config = config
	m.spaces = []Space{}
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
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
	for _, space := range m.spaces {
		m.choices = append(m.choices, space.Title)
	}
	m.loading = false

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
	if m.loading {
		s += "Loading..."
		return s
	} else {
		s += "Select a space to view its contents\n\n"
	}

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

	items := []item{
		{title: "Raspberry Pi's", desc: "I have 'em all over my house"},
	}

	for _, itm := range items {
		s += fmt.Sprintf("%s: %s\n", itm.Title(), itm.Description())
	}

	return docStyle.Render(items.View())
}
