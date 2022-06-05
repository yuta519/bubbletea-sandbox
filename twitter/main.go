package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

type model struct {
	status int
	err    error
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

// func (m model) Init() tea.Cmd
// func (m model) Update(tea.Model, tea.Cmd)

func main() {
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
