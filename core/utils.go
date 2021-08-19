package core

import (
	"fmt"
	"github.com/xjasonlyu/tun2socks/log"
	"os"
	"os/exec"
	"strings"
)

func Must(action string, err error) {
	if err != nil {
		log.Fatalf("failed to %s: %v", action, err)
	}
}

func Mustf(action string, err error, args ...interface{}) {
	if err != nil {
		Must(fmt.Sprintf(action, args), err)
	}
}

func Maybe(action string, err error) {
	if err != nil {
		log.Errorf("failed to %s: %v", action, err)
	}
}

func Maybef(action string, err error, args ...interface{}) {
	if err != nil {
		Maybe(fmt.Sprintf(action, args), err)
	}
}

func ExecShell(cmdIn string, errIn error, name string, arg ...string) (cmd string, err error) {
	if errIn != nil {
		return cmdIn, errIn
	}
	cmd = strings.Join([]string{name, strings.Join(arg, " ")}, " ")
	shell := exec.Command(name, arg...)
	log.Debugf(">> %s", cmd)
	shell.Stderr = os.Stderr
	err = shell.Start()
	if err == nil {
		err = shell.Wait()
	}
	return
}

func ExecProc(name string, arg []string) error {
	shell := exec.Command(name, arg...)
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	err := shell.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	status, err := shell.Process.Wait()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(status.ExitCode())
	return nil
}
