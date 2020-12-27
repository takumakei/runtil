package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/takumakei/runtil/clix"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap/zapcore"
)

var FlagSet = make(clix.FlagSet)

var Flags = []cli.Flag{
	UntilFlag,
	SignalFlag,
	KillAfterFlag,
	LogLevelFlag,
	DebugFlag,
}

var (
	UntilFlag = &cli.GenericFlag{
		Name:     "until",
		Aliases:  []string{"u"},
		Usage:    "process running until `HH:MM:SS`",
		EnvVars:  []string{"RUNTIL_UNTIL"},
		Required: true,
		Value:    UntilFlagValue,
	}

	UntilFlagValue = new(HHMMSSFlag)

	SignalFlag = &cli.GenericFlag{
		Name:    "signal",
		Aliases: []string{"s"},
		Usage:   "`signal` to send at HH:MM:SS",
		EnvVars: []string{"RUNTIL_SIGNAL"},
		Value:   SignalFlagValue,
	}

	SignalFlagValue = MustParseSignal("SIGINT")

	KillAfterFlag = &cli.DurationFlag{
		Name:        "kill-after",
		Aliases:     []string{"k"},
		Usage:       "waiting `duration` before sending KILL after the signal sent",
		EnvVars:     []string{"RUNTIL_KILL_AFTER"},
		Value:       10 * time.Second,
		Destination: new(time.Duration),
	}

	LogLevelFlag = &cli.GenericFlag{
		Name:    "log-level",
		Aliases: []string{"l"},
		Usage:   "`level` [debug|info|warn|error|dpanic|panic|fatal]",
		EnvVars: []string{"RUNTIL_LOG_LEVEL"},
		Value:   LogLevelFlagValue,
	}

	LogLevelFlagValue = new(zapcore.Level)

	DebugFlag = &cli.BoolFlag{
		Name:        "debug",
		Aliases:     []string{"d"},
		Usage:       "use debug config of zap logger",
		EnvVars:     []string{"RUNTIL_DEBUG"},
		Destination: new(bool),
	}
)

func FlagUntil() (HHMMSS, bool) {
	return UntilFlagValue.HHMMSS, FlagSet.IsSet(UntilFlag)
}

func FlagSignal() os.Signal {
	return SignalFlagValue.Signal()
}

func FlagKillAfter() time.Duration {
	return *KillAfterFlag.Destination
}

func FlagLogLevel() zapcore.Level {
	return *LogLevelFlagValue
}

func FlagDebug() bool {
	return *DebugFlag.Destination
}

type HHMMSSFlag struct {
	HHMMSS
	value string
}

func (hmsf *HHMMSSFlag) Set(value string) error {
	hmsf.value = value
	return hmsf.HHMMSS.Parse(value)
}

func (hmsf *HHMMSSFlag) String() string {
	return hmsf.value
}

type Signal struct {
	sig   os.Signal
	value string
}

func MustParseSignal(s string) *Signal {
	sig := new(Signal)
	if err := sig.Set(s); err != nil {
		panic(err)
	}
	return sig
}

func (sig *Signal) Set(value string) error {
	sig.value = value
	if i, err := strconv.Atoi(value); err == nil {
		sig.sig = syscall.Signal(i)
		return nil
	}
	value = strings.ToUpper(value)
	if strings.HasPrefix(value, "SIG") {
		value = value[3:]
	}
	if v, ok := signals[value]; ok {
		sig.sig = v
		return nil
	}
	return fmt.Errorf("cannot parse %q into signal", sig.value)
}

func (sig *Signal) String() string {
	return sig.value
}

func (sig *Signal) Signal() os.Signal {
	return sig.sig
}
