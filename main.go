package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mkirl/capacitour/api"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	// Load configuration
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

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	loading  bool
}

func initialModel(config *api.Config) tea.Model {
	m := model{}
	m.choices = []string{"Loading..."}
	m.cursor = 0
	m.selected = make(map[int]struct{})
	m.loading = true
	go api.FetchData(config)
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func sendToServer() tea.Msg {
	url := os.Getenv("SERVER_URL")
	if url == "" {
		fmt.Println("SERVER_URL environment variable is not set")
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error sending request: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
		return nil
	}

	fmt.Printf("Response from server: %s", body)

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				if m.choices[m.cursor] == "Send to server" {
					sendToServer()
				}
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Welcome to Capacitour?\n\n"

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
