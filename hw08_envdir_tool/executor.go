package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.
	var exitCodeErr *exec.ExitError

	if err := setEnv(&env); err != nil {
		fmt.Println(err.Error())
		return -1
	}
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	err := command.Run()
	if err != nil {
		ok := errors.As(err, &exitCodeErr)
		if ok {
			fmt.Println(err.Error())
			return exitCodeErr.ExitCode()
		}
		fmt.Println(err.Error())
		return -1
	}

	return 0
}

func setEnv(env *Environment) error {
	for k, v := range *env {
		if v.NeedRemove {
			if err := os.Unsetenv(k); err != nil {
				return err
			}
		} else {
			if err := os.Setenv(k, v.Value); err != nil {
				return err
			}
		}
	}
	return nil
}
