package darkcli

import (
	"fmt"

	"github.com/darklab8/go-utils/utils/ptr"
)

type Parser struct {
	actions_by_nick map[string]Action
	ParserOpts
}

type ParserOpts struct {
	ParentArgs    []string // needed for help info to see what was previous called
	DefaultAction *string  // by default will be "help"
}

func NewParser(
	actions []Action,
	opts ParserOpts,
) *Parser {
	p := &Parser{
		actions_by_nick: map[string]Action{},
		ParserOpts:      opts,
	}
	for _, action := range actions {
		p.actions_by_nick[action.Nickname] = action
	}
	help_cmd := &Action{
		Nickname:    "help",
		Description: "Print help info about sub group of commands for darkmap",
		Func: func(info ActionInfo) error {
			p.PrintHelp()
			return nil
		},
	}
	p.actions_by_nick["help"] = *help_cmd
	if p.DefaultAction == nil {
		p.DefaultAction = ptr.Ptr("help")
	}

	return p
}

type ActionInfo struct {
	CmdArgs []string
}

func (p *Parser) Run(args []string) error {
	var action_nickname string
	if len(args) >= 1 {
		action_nickname = args[0]
	}
	fmt.Println("act:", action_nickname)

	info := ActionInfo{CmdArgs: args}

	if action, ok := p.actions_by_nick[action_nickname]; ok {
		return action.Func(info)
	}

	return p.actions_by_nick[*p.DefaultAction].Func(info)
}

func (p *Parser) PrintHelp() {

	fmt.Println("you called commands:", p.ParentArgs)
	fmt.Println(`
but missed one more cli argument or called help
Input extra argument
Possible commands:`)
	for _, command := range p.actions_by_nick {
		fmt.Printf("  %s - %s\n", command.Nickname, command.Description)
	}
}

type Action struct {
	Nickname    string
	Description string
	Func        func(info ActionInfo) error
}

func (a Action) GetNickname() string { return string(a.Nickname) }
