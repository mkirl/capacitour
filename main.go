package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	api "github.com/mkirl/capacitour/api"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	docStyle          = lipgloss.NewStyle().Margin(1, 2)
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Title())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type model struct {
	currentView   string
	err           error
	config        api.Config
	spaces        []api.Space
	spinner       spinner.Model
	list          list.Model
	choice        string
	quitting      bool
	selectedSpace api.Space
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	selectedSpace := api.Space{}

	items := []list.Item{}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Capacitour"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return model{
		currentView:   "loading",
		spinner:       s,
		list:          l,
		selectedSpace: selectedSpace,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			config, err := loadConfig()
			if err != nil {
				return errMsg{err: err}
			}
			m.config = config
			spaces, err := api.FetchAllSpaces(&m.config)
			if err != nil {
				return errMsg{err: err}
			}
			items := make([]list.Item, len(spaces))
			for i, space := range spaces {
				items[i] = item{title: space.Title, desc: "Description for " + space.Title}
			}
			return spacesResponseMsg{spaces: spaces, items: items}
		},
	)
}

func loadConfig() (api.Config, error) {
	apiURL := os.Getenv("CAPACITIES_API_URL")
	if apiURL == "" {
		return api.Config{}, fmt.Errorf("CAPACITIES_API_URL is not set")
	}

	apiToken := os.Getenv("CAPACITIES_API_TOKEN")
	if apiToken == "" {
		return api.Config{}, fmt.Errorf("CAPACITIES_API_TOKEN is not set")
	}

	return api.Config{APIURL: apiURL, APIToken: apiToken}, nil
}

type spacesResponseMsg struct {
	spaces []api.Space
	items  []list.Item
}

type spaceSelectedMsg struct {
	space api.Space
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		fmt.Println("Error:", msg.err)
		return m, nil

	case spacesResponseMsg:
		m.spaces = msg.spaces
		m.list.SetItems(msg.items) // Ensure items are set correctly
		m.currentView = "spaces"
		fmt.Println("Spaces loaded:", len(msg.spaces))
		return m, nil

	case spaceSelectedMsg:
		m.selectedSpace = msg.space
		m.currentView = "space"
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.title
			}
			m.handleSpaceSelection(m.spaces[m.list.Index()])
			return m, func() tea.Msg {
				return spaceSelectedMsg{space: m.spaces[m.list.Index()]}
			}
		case "up", "k":
			m.list.CursorUp()
		case "down", "j":
			m.list.CursorDown()
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		fmt.Println("Window resized:", msg.Width, msg.Height)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
func (m model) handleSpaceSelection(space api.Space) tea.Cmd {
	spaceData, err := api.FetchSpace(&m.config, space)
	if err != nil {
		fmt.Println("Error fetching space:", err)
		panic(err)
	}
	fmt.Println(string(spaceData))
	return nil
}

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That's cool.")
	}
	switch m.currentView {
	case "loading":
		return m.spinner.View()
	case "spaces":
		return docStyle.Render(m.list.View())
	case "space":
		return docStyle.Render(m.selectedSpace.Title)
	}
	return "Unknown view"
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("./debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
	m := initialModel()

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
