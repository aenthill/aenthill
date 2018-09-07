---
title: "Getting started"
weight: 2
kind: "use"
---

# Installation

Aenthill is distributed in a binary form and the only requirement is to have Docker installed and running on your host (always use the latest version).

## Linux/MacOS

The quickest way to get Aenthill is to run the following command:

```bash
$ curl -sf https://raw.githubusercontent.com/aenthill/aenthill/master/install.sh | BINDIR=/usr/local/bin sh
```

It will installs Aenthill in the `/usr/local/bin` directory. You may change the installation path by updating the `BINDIR` value.

You may also install a specific version with:

```bash
$ curl -sf https://raw.githubusercontent.com/aenthill/aenthill/master/install.sh | BINDIR=/usr/local/bin sh -s version
```

All available versions are listed in the [releases pages](https://github.com/aenthill/aenthill/releases)

## Windows

# Commands