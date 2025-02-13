package shell

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

type Logger interface {
	Error(ctx context.Context, msg string, additionalFields ...map[string]any)
	Warn(ctx context.Context, msg string, additionalFields ...map[string]any)
	Info(ctx context.Context, msg string, additionalFields ...map[string]any)
	Debug(ctx context.Context, msg string, additionalFields ...map[string]any)
	Trace(ctx context.Context, msg string, additionalFields ...map[string]any)
}

// RunCommand runs a script in a given working directory.
func RunCommand(ctx context.Context, interpreter []string, env map[string]string, dir, command string, logger Logger) error {
	cmd := exec.CommandContext(ctx, interpreter[0], append(interpreter[1:], command)...)
	cmd.Dir = dir

	setEnv(cmd, env, true)

	if logger == nil {
		cmd.Stdout = nil
		cmd.Stderr = nil

		return cmd.Run()
	}

	return runCommandLogOutput(ctx, cmd, logger)
}

// setEnv sets the environment variables for a command.
func setEnv(cmd *exec.Cmd, env map[string]string, addOS bool) {
	envList := make([]string, 0, len(env))
	for k, v := range env {
		envList = append(envList, fmt.Sprintf("%s=%s", k, v))
	}

	if addOS {
		cmd.Env = append(os.Environ(), envList...)
	} else {
		cmd.Env = envList
	}
}

// runCommandLogOutput runs a command and logs the output prefixed with [<LEVEL>].
func runCommandLogOutput(ctx context.Context, cmd *exec.Cmd, logger Logger) error {
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout

	scanner := bufio.NewScanner(pipe)
	err = cmd.Start()
	if err != nil {
		return err
	}

	regex := regexp.MustCompile(`^\[(ERROR|WARN|INFO|DEBUG|TRACE)\]\s*(.+)`)
	for scanner.Scan() {
		matches := regex.FindStringSubmatch(scanner.Text())
		if matches != nil {
			switch matches[1] {
			case "ERROR":
				logger.Error(ctx, matches[2])
			case "WARN":
				logger.Warn(ctx, matches[2])
			case "INFO":
				logger.Info(ctx, matches[2])
			case "DEBUG":
				logger.Debug(ctx, matches[2])
			case "TRACE":
				logger.Trace(ctx, matches[2])
			default:
			}
		}
	}

	if scanner.Err() != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return scanner.Err()
	}

	return cmd.Wait()
}
