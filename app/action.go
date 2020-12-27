package app

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func Action(c *cli.Context) error {
	args := c.Args().Slice()
	if len(args) == 0 {
		cli.ShowSubcommandHelpAndExit(c, 0)
	}

	until, ok := FlagUntil()
	if !ok {
		cli.ShowSubcommandHelp(c)
		return errors.New(`Required flag "until" not set`)
	}

	if err := InitZap(FlagLogLevel()); err != nil {
		return err
	}

	pending := FlagSignal()
	killAfter := FlagKillAfter()

	now := time.Now()
	kil := until.Next(now)
	zap.L().Info(
		"init",
		zap.String("until", kil.In(time.Local).Format(time.RFC3339Nano)),
		zap.String("duration", kil.Sub(now).String()),
		zap.String("pending", pending.String()),
		zap.String("kill-after", killAfter.String()),
	)

	pcs := Start(args[0], args[1:]...)

	zap.L().Info(
		"run",
		zap.Strings("command", args),
		zap.Int("pid", pcs.Pid),
	)

	err := func() error {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigint)

		if ok, err := deadline(sigint, pcs, kil, pending); ok {
			return err
		}
		return kill(sigint, pcs, killAfter)
	}()

	zap.L().Info(
		"exit",
		zap.Int("pid", pcs.Pid),
		zap.Error(err),
	)

	if err, ok := err.(*exec.ExitError); ok {
		os.Exit(err.ExitCode())
	}

	return err
}

func deadline(sigint <-chan os.Signal, r *Process, d time.Time, sig os.Signal) (bool, error) {
	kill := func(sig os.Signal) {
		err := r.Kill(sig)
		zap.L().Info(
			"kill",
			zap.String("signal", sig.String()),
			zap.Int("pid", r.Pid),
			zap.Error(err),
		)
	}

	next := func(t time.Time) time.Duration {
		a := d.Sub(t)
		switch {
		case a > time.Hour:
			return time.Hour
		case a >= 100*time.Millisecond:
			return a / 2
		}
		return a
	}

	to := next(time.Now())
	zap.L().Debug("next", zap.String("timeout", to.String()))
	t := time.NewTimer(to)
	defer t.Stop()
	for {
		select {
		case v := <-t.C:
			if v.Before(d) {
				to := next(time.Now())
				zap.L().Debug("next", zap.String("timeout", to.String()))
				t.Reset(to)
				continue
			}
			zap.L().Debug("kill")
			kill(sig)
			return false, nil

		case sig := <-sigint:
			kill(sig)

		case err := <-r.C:
			return true, err
		}
	}
}

func kill(sigint <-chan os.Signal, r *Process, d time.Duration) error {
	t := time.NewTimer(d)
	if d <= 0 {
		zap.L().Info(
			"wait forever for the end of the process",
			zap.Int("pid", r.Pid),
		)
		if !t.Stop() {
			<-t.C
		}
	} else {
		zap.L().Info(
			"wait to send SIGKILL",
			zap.String("wait", d.String()),
			zap.Int("pid", r.Pid),
		)
		defer t.Stop()
	}

	kill := func(sig os.Signal) {
		err := r.Kill(sig)
		zap.L().Info(
			"kill",
			zap.String("signal", sig.String()),
			zap.Int("pid", r.Pid),
			zap.Error(err),
		)
	}

	for {
		select {
		case <-t.C:
			kill(syscall.SIGKILL)

		case sig := <-sigint:
			kill(sig)

		case err := <-r.C:
			return err
		}
	}
}
