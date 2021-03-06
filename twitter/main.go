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
	tweets       []anaconda.Tweet
	status       int
	err          error
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
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

	case tea.KeyMsg:
		// The actual key pressed
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.accounts)-1 {
				m.cursor++
			}
		case "!":
			m.isTextFormat = true
			return m, nil
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
				m.isChosen = false
			} else {
				m.selected[m.cursor] = struct{}{}
				m.isChosen = true
			}
		case "enter":
			anaconda.SetConsumerKey(os.Getenv("API_KEY"))
			anaconda.SetConsumerSecret(os.Getenv("API_SECRET"))
			api := anaconda.NewTwitterApi(os.Getenv("ACCESS_KEY"), os.Getenv("ACCESS_SECRET"))
			values := url.Values{}

			if m.isChosen {
				fmt.Println(m.accounts[m.cursor])
				values.Set("screen_name", m.accounts[m.cursor])
				m.fetchTweetsByAccount(api, values)
			}
			if m.isTextFormat {
				values.Set("screen_name", m.textInput.Value())
				m.fetchTweetsByAccount(api, values)
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
	return choicesView(m)
}

func choicesView(m model) string {
	// The header
	s := "Choose a twitter account you want check.\n\n"
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
	// The footer
	s += "\nOr you could input a twitter account with free text when press `!`.\n"
	s += "\nPress q to quit.\n"
	// Send the UI for rendering
	return s
}

func textInputView(m model) string {
	return fmt.Sprintf(
		"What???s the twitter account you want know?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func (m model) fetchTweetsByAccount(api *anaconda.TwitterApi, v url.Values) tea.Msg {
	tweets, err := api.GetUserTimeline(v)
	if err != nil {
		panic(err)
	}
	for _, tweet := range tweets {
		fmt.Println("tweet: ", tweet.Text)
	}
	return model{tweets: tweets}
}
