package core

import (
	"golang.org/x/sys/windows"
	"os"
	"strings"
	"syscall"
)

func ExecSu() error {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err == nil {
		return nil
	}

	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	err = windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, windows.SW_NORMAL)
	if err == nil {
		os.Exit(0)
	}
	return nil
}
