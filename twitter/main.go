package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	status   int
	err      error
	choices  []string
	cursor   int
	selected map[int]struct{}
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

func initialModel() model {
	return model{
		// Our shopping list is a grocery list
		choices: []string{"golangch", "GolangTrends", "golang_news"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

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
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

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

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
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

	var userNameArg string
	fmt.Println("Please input twitter user")
	fmt.Scan(&userNameArg)
	values := url.Values{}
	values.Set("screen_name", userNameArg)
	showTimeLine(api, values)
}
