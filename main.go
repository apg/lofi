package main

import (
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/apg/flenv"
	"github.com/apg/go-trello"
)

type config struct {
	AppKey string `env:"TRELLO_APP_KEY" flag:"-k,--app-key" help:"Trello app key"`
	Token  string `env:"TRELLO_TOKEN" flag:"-t,--token" help:"Trello token"`
	Board  string `env:"TRELLO_BOARD" flag:"-b,--board" help:"Trello board id with list"`
	List   string `env:"TRELLO_LIST" flag:"-l,--list" help:"Trello list to export"`
}

func findBoard(c *trello.Client, id string) (board trello.Board, err error) {
	b, err := c.Board(id)
	if err != nil {
		return board, fmt.Errorf("Error while retrieving board %q: %q", id, err)
	}

	return *b, nil
}

func findList(board trello.Board, name string) (list trello.List, err error) {
	lists, err := board.Lists()
	if err != nil {
		return list, fmt.Errorf("Error while retrieving lists for board %q: %q", board.Name, err)
	}

	for _, l := range lists {
		if l.Name == name {
			return l, nil
		}
	}

	return list, fmt.Errorf("No list found named %q, in board %q", name, board.Name)
}

func findCards(list trello.List) (cards []trello.Card, err error) {
	cards, err = list.Cards()
	if err != nil {
		err = fmt.Errorf("Error while retrieving cards for list %q: %q", list.Name, err)
	}
	return
}

func main() {
	var conf config
	flagSet, err := flenv.DecodeArgs(&conf)
	if err != nil {
		flagSet.Usage()
		os.Exit(1)
	}

	if conf.AppKey == "" || conf.Token == "" || conf.Board == "" || conf.List == "" {
		fmt.Fprintf(os.Stderr, "TRELLO_APP_KEY, TRELLO_TOKEN, TRELLO_BOARD, and TRELLO_LIST are all required, or use the flags.")
		fmt.Fprintf(os.Stderr, "%+v\n", conf)
		os.Exit(1)
	}

	client, err := trello.NewAuthClient(conf.AppKey, &conf.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		os.Exit(1)
	}

	board, err := findBoard(client, conf.Board)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		os.Exit(1)
	}

	list, err := findList(board, conf.List)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		os.Exit(1)
	}

	cards, err := findCards(list)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		os.Exit(1)
	}

	tmpl.Execute(os.Stdout, struct {
		List  trello.List
		Cards []trello.Card
		Now   string
	}{list, cards, time.Now().Format("2006-02-01")})
}

var tmplRaw = `
% {{.List.Name}}: WITTY_TITLE
%
% {{.Now}}

# {{.List.Name}}: WITTY_TITLE

SOME_INTRO_TEXT

# Hackity Hacks (in pseudo random order)

{{range .Cards}}
## {{.Name}}

{{.Desc}}

{{end}}

# SOME_OTHER_STUFF??

# SOME_WITTY_CLOSING

Well, that was a lot of hacks! I guess this is what happens when you
skip a week!

I'm always on the lookout for interesting hacks, projects, pieces of
art, neat papers, etc. If you have something I should see, don't
hesitate to reply to this email!

Oh! And, are you or someone you know working on a big project with an
interesting story (or even a boring story, but with interesting
problems)? Let me know! I'd love to ch chat with you/them.

Until next time!

Happy hacking,

Andrew
`

var tmpl = template.Must(template.New("hifi").Parse(tmplRaw))
