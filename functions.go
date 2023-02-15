package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
)

func openService(mgrHandle windows.Handle, name string) (*mgr.Service, error) {
	n, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	h, err := windows.OpenService(mgrHandle, n, windows.SERVICE_QUERY_STATUS|windows.SERVICE_START|windows.SERVICE_STOP|windows.SERVICE_USER_DEFINED_CONTROL)
	if err != nil {
		return nil, err
	}
	return &mgr.Service{Name: name, Handle: h}, nil
}

func enableSeDebugPrivilege() error {
	var t windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_ALL_ACCESS, &t); err != nil {
		return err
	}

	var luid windows.LUID

	if err := windows.LookupPrivilegeValue(nil, windows.StringToUTF16Ptr(seDebugPrivilege), &luid); err != nil {
		return fmt.Errorf("LookupPrivilegeValueW failed, error: %v", err)
	}

	ap := windows.Tokenprivileges{
		PrivilegeCount: 1,
	}

	ap.Privileges[0].Luid = luid
	ap.Privileges[0].Attributes = windows.SE_PRIVILEGE_ENABLED

	if err := windows.AdjustTokenPrivileges(t, false, &ap, 0, nil, nil); err != nil {
		return fmt.Errorf("AdjustTokenPrivileges failed, error: %v", err)
	}

	return nil
}

func parseProcessName(exeFile [windows.MAX_PATH]uint16) string {
	for i, v := range exeFile {
		if v <= 0 {
			return string(utf16.Decode(exeFile[:i]))
		}
	}
	return ""
}

func getTrustedInstallerPid() (uint32, error) {

	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(snapshot)

	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	if err := windows.Process32First(snapshot, &procEntry); err != nil {
		return 0, err
	}

	for {
		if strings.EqualFold(parseProcessName(procEntry.ExeFile), tiExecutableName) {
			return procEntry.ProcessID, nil
		} else {
			if err = windows.Process32Next(snapshot, &procEntry); err != nil {
				if err == windows.ERROR_NO_MORE_FILES {
					break
				}
				return 0, err
			}
		}
	}
	return 0, fmt.Errorf("cannot find %v in running process list", tiExecutableName)
}

func checkIfAdmin() bool {
	f, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func elevate() error {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	if err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd); err != nil {
		return err
	}

	os.Exit(0)
	return nil
}
