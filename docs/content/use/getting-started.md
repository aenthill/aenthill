---
title: "Getting started"
weight: 2
kind: "use"
---

# Installation

Aenthill is distributed in a binary form and the only requirement is to have Docker installed and running on your host (always use the latest version).

We currently provide pre-built binaries for the following:

- Linux
- macOS
- Windows

You may find the latest release on the [releases pages](https://github.com/aenthill/aenthill/releases).

## Install script (Unix)

The quickest way to get Aenthill is to run the following command:

```bash
$ curl -sf https://raw.githubusercontent.com/aenthill/aenthill/master/install.sh | BINDIR=/usr/local/bin sh
```

It will installs Aenthill in the `/usr/local/bin` directory. You may change the installation path by updating the `BINDIR` value.

## Homebrew (macOS)

If you are on macOs and use [Homebrew](https://brew.sh/) for package management
you may install Aenthill with the following command:

```bash
$ brew install aenthill/tap/aenthill
```

## Scoop (Windows)

If you are on Windows and use [Scoop](https://scoop.sh/) for package management
you may install Aenthill with the following commands:

```bash
$ scoop bucket add aenthill https://github.com/aenthill/scoop-bucket.git
$ scoop install aenthill
```


# Commands