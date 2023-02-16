# TrustedInstaller
A simple Proof of Concept in Golang to start a new shell as TrustedInstaller. This code accompanies FourCore's blog about TrustedInstaller. It is important to note that you need to run this as a user which has SeDebugPrivileges. Upon running, it will automatically ask for UAC in case you are not running as an Administrator.

Use the `RunAsTrustedInstaller` function to pass any executable to be run with TrustedInstaller privileges.

To run 
1. git clone the repository
2. ensure you have go compiler installed
3. You can either build a binary using `go build ti` or run it directly using `go run ti`

It will spawn a new cmd shell as TrustedInstaller which you can check by running `whoami /all`
![demo](https://user-images.githubusercontent.com/26490648/219342533-79d0cf34-0bf2-4f63-b805-34fca5aff012.gif)
