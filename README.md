[# TrustedInstaller

A simple Proof of Concept in Go to spawn a new shell as TrustedInstaller. Read more about how this PoC works on this [blog about TrustedInstaller](https://fourcore.io/blogs/no-more-access-denied-i-am-trustedinstaller). It is important to note that this should be executed as a user which has SeDebugPrivileges. Upon execution, it will automatically ask for UAC in case it is not executed as as an Administrator.

## POC

1. Clone the repository

```
$ git clone https://github.com/FourCoreLabs/TrustedInstallerPOC.git
```

2. Ensure you have Go installed. This POC has been tested on Go 1.19.
3. Either build the binary and execute it

```
$ go build ti
$ ./ti.exe
```

4. Or run it directly

```
$ go run ti
```


This will spawn a new cmd shell with TrustedInstaller privileges which can be confirmed by running the command `whoami /all`

![demo](https://user-images.githubusercontent.com/26490648/219342533-79d0cf34-0bf2-4f63-b805-34fca5aff012.gif)

## API

- RunAsTrustedInstaller
  - Use the `RunAsTrustedInstaller` function to pass any executable to be run with TrustedInstaller privileges.
