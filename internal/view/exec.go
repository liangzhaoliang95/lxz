// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package view

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/liangzhaoliang95/lxz/internal/render"
	"github.com/liangzhaoliang95/lxz/internal/slogs"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

const (
	shellCheck   = `command -v bash >/dev/null && exec bash || exec sh`
	bannerFmt    = "<<LXZ-Shell>> Container: %s \n"
	outputPrefix = "[output]"
)

var editorEnvVars = []string{"K9S_EDITOR", "KUBE_EDITOR", "EDITOR"}

type shellOpts struct {
	clear, background bool
	pipes             []string
	binary            string
	banner            string
	args              []string
}

func (s shellOpts) String() string {
	return fmt.Sprintf("%s %s", s.binary, strings.Join(s.args, " "))
}

func runDockerExec(a *App, opts *shellOpts) error {
	bin, err := exec.LookPath("docker")
	if errors.Is(err, exec.ErrDot) {
		return fmt.Errorf("docker command must not be in the current working directory: %w", err)
	}
	if err != nil {
		return fmt.Errorf("docker command is not in your path: %w", err)
	}
	opts.binary = bin

	suspended, errChan, stChan := run(a, opts)
	if !suspended {
		return fmt.Errorf("unable to run command")
	}
	for v := range stChan {
		slog.Debug("stdout", slogs.Line, v)
	}
	var errs error
	for e := range errChan {
		errs = errors.Join(errs, e)
	}

	return errs
}

func runK9sExec(a *App, opts *shellOpts) error {
	bin, err := exec.LookPath("k9s")
	if errors.Is(err, exec.ErrDot) {
		return fmt.Errorf("k9s command must not be in the current working directory: %w", err)
	}
	if err != nil {
		return fmt.Errorf("k9s command is not in your path: %w", err)
	}
	slog.Debug("K9s exec command found", slogs.Command, bin)
	opts.binary = bin
	slog.Info("options", slogs.Options, opts)
	suspended, errChan, stChan := run(a, opts)
	if !suspended {
		return fmt.Errorf("unable to run command")
	}
	for v := range stChan {
		slog.Debug("stdout", slogs.Line, v)
	}
	var errs error
	for e := range errChan {
		errs = errors.Join(errs, e)
	}

	return errs
}

func run(a *App, opts *shellOpts) (ok bool, errC chan error, outC chan string) {
	errChan := make(chan error, 1)
	statusChan := make(chan string, 1)

	if opts.background {
		if err := execute(opts, statusChan); err != nil {
			errChan <- err
			a.UI.Flash().Errf("Exec failed %q: %s", opts, err)
		}
		close(errChan)
		return true, errChan, statusChan
	}

	a.Halt()
	defer a.Resume()

	return a.UI.Suspend(func() {
		if err := execute(opts, statusChan); err != nil {
			errChan <- err
			a.UI.Flash().Errf("Exec failed %q: %s", opts, err)
		}
		close(errChan)
	}), errChan, statusChan
}

func execute(opts *shellOpts, statusChan chan<- string) error {
	if opts.clear {
		clearScreen()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if !opts.background {
			cancel()
			clearScreen()
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func(cancel context.CancelFunc) {
		defer slog.Debug("Got signal canceled")
		select {
		case sig := <-sigChan:
			slog.Debug("Command canceled with signal", slogs.Sig, sig)
			cancel()
		case <-ctx.Done():
			slog.Debug("Signal context canceled!")
		}
	}(cancel)

	cmds := make([]*exec.Cmd, 0, 1)
	cmd := exec.CommandContext(ctx, opts.binary, opts.args...)
	slog.Debug("Exec command", slogs.Command, opts)

	if env := os.Getenv("K9S_EDITOR"); env != "" {
		// There may be situations where the user sets the editor as the binary
		// followed by some arguments (e.g. "code -w" to make it work with vscode)
		//
		// In such cases, the actual binary is only the first token
		binTokens := strings.Split(env, " ")

		if bin, err := exec.LookPath(binTokens[0]); err == nil {
			binTokens[0] = bin
			cmd.Env = append(os.Environ(), fmt.Sprintf("KUBE_EDITOR=%s", strings.Join(binTokens, " ")))
		}
	}

	cmds = append(cmds, cmd)

	for _, p := range opts.pipes {
		tokens := strings.Split(p, " ")
		if len(tokens) < 2 {
			continue
		}
		cmd := exec.CommandContext(ctx, tokens[0], tokens[1:]...)
		slog.Debug("Exec command", slogs.Command, cmd)
		cmds = append(cmds, cmd)
	}

	var o, e bytes.Buffer
	err := pipe(ctx, opts, statusChan, &o, &e, cmds...)
	if err != nil {
		slog.Error("Exec failed",
			slogs.Error, err,
			slogs.Command, cmds,
		)
		return errors.Join(err, fmt.Errorf("%s", e.String()))
	}

	return nil
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func pipe(_ context.Context, opts *shellOpts, statusChan chan<- string, w, e *bytes.Buffer, cmds ...*exec.Cmd) error {
	if len(cmds) == 0 {
		return nil
	}

	if len(cmds) == 1 {
		cmd := cmds[0]
		if opts.background {
			go func() {
				cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, w, e
				if err := cmd.Run(); err != nil {
					slog.Error("Command exec failed", slogs.Error, err)
				} else {
					for _, l := range strings.Split(w.String(), "\n") {
						if l != "" {
							statusChan <- fmt.Sprintf("%s %s", outputPrefix, l)
						}
					}
					statusChan <- fmt.Sprintf("Command completed successfully: %q", render.Truncate(cmd.String(), 20))
					slog.Info("Command ran successfully", slogs.Command, cmd.String())
				}
				close(statusChan)
			}()
			return nil
		}
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		_, _ = cmd.Stdout.Write([]byte(opts.banner))

		slog.Debug("Exec started")
		err := cmd.Run()
		slog.Debug("Running exec done", slogs.Error, err)
		if err == nil {
			statusChan <- fmt.Sprintf("Command completed successfully: %q", cmd.String())
		}
		close(statusChan)

		return err
	}

	last := len(cmds) - 1
	for i := range cmds {
		cmds[i].Stderr = os.Stderr
		if i+1 < len(cmds) {
			r, w := io.Pipe()
			cmds[i].Stdout, cmds[i+1].Stdin = w, r
		}
	}
	cmds[last].Stdout = os.Stdout

	for _, cmd := range cmds {
		slog.Debug("Starting command", slogs.Command, cmd)
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	return cmds[len(cmds)-1].Wait()
}
