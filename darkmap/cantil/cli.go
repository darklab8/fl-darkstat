package cantil

/*
 */

import (
	"fmt"
	"os"

	"github.com/darklab8/fl-darkstat/darkstat/settings"
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

func NewConsoleParser(
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
	fmt.Println("HELP INFO")

	fmt.Println()
	fmt.Println("possible environment variables:")
	for _, enver := range settings.Enverants {
		fmt.Println()
		if enver.Description != "" {
			fmt.Println(enver.Description)
		} else {
			fmt.Println("ENVERANT env block:")
		}

		for _, env_var := range enver.GetParams() {
			fmt.Printf(" %s - [%s] %s (default: %v)\n", env_var.PrefixedKey, env_var.VarType.ToStr(), env_var.Description, env_var.Default)
		}
	}

	fmt.Println()
	fmt.Println("Possible commands:")
	for _, command := range p.actions_by_nick {
		fmt.Printf("  %s - %s\n", command.Nickname, command.Description)
	}

	fmt.Println()
	fmt.Println("your called args", os.Args[1:])
	fmt.Println("command parent args:", p.ParentArgs)
	fmt.Println(`missed one more cli argument or called help. Input extra argument if necessary`)
}

type Action struct {
	Nickname    string
	Description string
	Func        func(info ActionInfo) error
}

func (a Action) GetNickname() string { return string(a.Nickname) }
