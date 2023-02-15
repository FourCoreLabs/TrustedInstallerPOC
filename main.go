package main

import (
	"fmt"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"golang.org/x/sys/windows"
)

const (
	seDebugPrivilege = "SeDebugPrivilege"
	tiServiceName    = "TrustedInstaller"
	tiExecutableName = "trustedinstaller.exe"
)

func RunAsTrustedInstaller(path string, args []string) error {
	if !checkIfAdmin() {
		if err := elevate(); err != nil {
			return fmt.Errorf("cannot elevate Privs: %v", err)
		}
	}

	if err := enableSeDebugPrivilege(); err != nil {
		return fmt.Errorf("cannot enable %v: %v", seDebugPrivilege, err)
	}

	svcMgr, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("cannot connect to svc manager: %v", err)
	}

	s, err := openService(svcMgr.Handle, tiServiceName)
	if err != nil {
		return fmt.Errorf("cannot open ti service: %v", err)
	}

	status, err := s.Query()
	if err != nil {
		return fmt.Errorf("cannot query ti service: %v", err)
	}

	if status.State != svc.Running {
		if err := s.Start(); err != nil {
			return fmt.Errorf("cannot start ti service: %v", err)
		} else {
			defer s.Control(svc.Stop)
		}
	}

	tiPid, err := getTrustedInstallerPid()
	if err != nil {
		return err
	}

	hand, err := windows.OpenProcess(windows.PROCESS_CREATE_PROCESS|windows.PROCESS_DUP_HANDLE|windows.PROCESS_SET_INFORMATION, true, tiPid)
	if err != nil {
		return fmt.Errorf("cannot open ti process: %v", err)
	}

	cmd := exec.Command(path, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_CONSOLE,
		ParentProcess: syscall.Handle(hand),
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("cannot start new process: %v", err)
	}

	fmt.Println("Started process with PID", cmd.Process.Pid)
	return nil
}

func main() {
	if err := RunAsTrustedInstaller("cmd.exe", []string{"/c", "start", "cmd.exe"}); err != nil {
		panic(err)
	}
}
