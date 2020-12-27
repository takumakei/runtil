package app

import (
	"github.com/urfave/cli/v2"
)

var Version = "unset"
var Commit = "unset"

var App = &cli.App{
	Usage:   "execute command until sometime in a day",
	Flags:   Flags,
	Before:  FlagSet.Init,
	Action:  Action,
	Version: Version + " (" + Commit + ") ",

	HideHelpCommand: true,
}
