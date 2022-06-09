package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	accounts     []string
	choice       int
	isChosen     bool
	cursor       int
	selected     map[int]struct{}
	isTextFormat bool
	textInput    textinput.Model
	status       int
	err          error
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	anaconda.SetConsumerKey(os.Getenv("API_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("API_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("ACCESS_KEY"), os.Getenv("ACCESS_SECRET"))

	fmt.Println("Please input twitter user")
	var userNameArg string
	fmt.Scan(&userNameArg)
	values := url.Values{}
	values.Set("screen_name", userNameArg)
	showTimeLine(api, values)
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		accounts:     []string{"golangch", "GolangTrends", "golang_news"},
		isChosen:     false,
		isTextFormat: false,
		textInput:    ti,
		selected:     make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.accounts)-1 {
				m.cursor++
			}
		case "f":
			m.isTextFormat = true
			return m, nil
		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
		m.textInput, cmd = m.textInput.Update(msg)
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	return m, cmd
}

func (m model) View() string {
	if m.isTextFormat {
		return textInputView(m)
	}

	// The header
	s := "Choose a twitter account you want check.\n\n"

	if m.isChosen {
		return choicesView(m)
	}

	// Iterate over our choices
	for i, choice := range m.accounts {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nOr you could input a twitter account with free text when press f.\n"
	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func choicesView(m model) string {
	c := m.choice
	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		checkbox("Plant carrots", c == 0),
		checkbox("Go to the market", c == 1),
		checkbox("Read something", c == 2),
		checkbox("See friends", c == 3),
	)
	footer := "\nOr you could input a twitter account with free text when press f.\n"
	return fmt.Sprintf("%s\n%s\n", choices, footer)
}

func checkbox(label string, checked bool) string {
	if checked {
		return fmt.Sprint("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func textInputView(m model) string {
	return fmt.Sprintf(
		"What’s your favorite Pokémon?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func showTimeLine(api *anaconda.TwitterApi, v url.Values) {
	tweets, err := api.GetUserTimeline(v)
	if err != nil {
		panic(err)
	}
	for _, tweet := range tweets {
		fmt.Println("tweet: ", tweet.Text)
	}
}
