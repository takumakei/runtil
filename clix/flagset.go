package clix

import "github.com/urfave/cli/v2"

type FlagSet map[string]struct{}

func (fs FlagSet) Init(c *cli.Context) error {
	fs.set(c)
	return nil
}

func (fs FlagSet) Chain(before cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		fs.set(c)
		return before(c)
	}
}

func (fs FlagSet) set(c *cli.Context) {
	for _, name := range c.LocalFlagNames() {
		fs[name] = struct{}{}
	}
}

func (fs FlagSet) IsSet(flag cli.Flag) bool {
	if flag.IsSet() {
		return true
	}
	for _, name := range flag.Names() {
		if _, ok := fs[name]; ok {
			return true
		}
	}
	return false
}
