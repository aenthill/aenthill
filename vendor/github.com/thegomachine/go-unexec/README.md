<p align="center">
    <img src="https://user-images.githubusercontent.com/8983173/41920404-f4baf4e2-7960-11e8-8880-6b54bcef12e2.png" alt="Logo" width="200" height="200" />
</p>
<h3 align="center">go-unexec</h3>
<p align="center">A simple library easing the use of <a href="https://golang.org/pkg/os/exec/">os/exec</a> package for running cross-platform external commands</p>
<p align="center">
    <a href="https://travis-ci.org/thegomachine/go-unexec">
        <img src="https://travis-ci.org/thegomachine/go-unexec.svg?branch=master" alt="Travis CI">
    </a>
    <a href="https://godoc.org/github.com/thegomachine/go-unexec">
        <img src="https://godoc.org/github.com/thegomachine/go-unexec?status.svg" alt="GoDoc">
    </a>
    <a href="https://goreportcard.com/report/thegomachine/go-unexec">
        <img src="https://goreportcard.com/badge/github.com/thegomachine/go-unexec" alt="Go Report Card">
    </a>
    <a href="https://codecov.io/gh/thegomachine/go-unexec/branch/master">
        <img src="https://codecov.io/gh/thegomachine/go-unexec/branch/master/graph/badge.svg" alt="Codecov">
    </a>
</p>

---

While using the `os/exec` package, you may have encountered some consistency issues:
a command which was working fine on your command line interpreter fails miserably while calling it
with the said package.

To address this common problem, the go-unexec library tries to detect your default command line
interpreter by looking for the `SHELL` environment variable on UNIX systems or `COMSPEC` environment variable
on Windows.

## Installation

```bash
$ go get github.com/thegomachine/go-unexec
```

## Usage

So previously your code might have looked like this:

```golang
import "os/exec"

func main() {
    cmd := exec.Command("echo", "Hello world")
    // will run "echo Hello world".
 }
```

With this package:

```golang
import unexec "github.com/thegomachine/go-unexec"

func main() {
    cmd, err := unexec.Command("echo", "Hello world")
    // will run "/bin/sh -c echo Hello world" (or "/bin/zsh -c echo Hello world" etc.)
    // on UNIX systems or "cmd.exe /c echo Hello world" on Windows.
}
```

See [GoDoc](https://godoc.org/github.com/thegomachine/go-unexec) for full documentation.

## FAQ

**When should I use this library?**

If your external commands are not platform specific.

**When should I NOT use this library?**

For every others cases! 😄