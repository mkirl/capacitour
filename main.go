package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	loading  bool
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	m := model{
		choices:  []string{"Loading..."},
		selected: make(map[int]struct{}),
		loading:  true,
	}
	go fetchData(&m)
	return m
}

func fetchData(m *model) {
	url := os.Getenv("SERVER_URL")
	if url == "" {
		fmt.Println("SERVER_URL environment variable is not set")
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
		return
	}

	// Update the model with the fetched data
	m.choices = []string{"Choice 1", "Choice 2", "Choice 3", string(body)}
	m.loading = false
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

func (m model) Init() tea.Cmd {
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
