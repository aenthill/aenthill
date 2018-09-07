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

Aenthill provides a set of commands allowing you to launch aents.

To view details for a command at any time use `aenthill help` or `aenthill help command`.

```
NAME:
   aenthill - May the swarm be with you!

USAGE:
   aenthill [global options] command [command options] [arguments...]

COMMANDS:
     start, s    Starts an aent
     add, a      Adds an aent in the manifest
     upgrade, u  Upgrades Aenthill
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, -v  print verbose output to the console
   --debug, -d    print debug output to the console
   --help, -h     show help
   --version      print the version
```

## Add

The command `add` is used to *install* an aent of type [orchestrator](/use/orchestrators) in the manifest `aenthill.json`.

For instance, if you want to add Docker Compose orchestrator, you may run:

```bash
$ aenthill add theaentmachine/aent-docker-compose
```

## Start

The command `start` is used to *launch* an aent of type [service](/use/services).

For instance, if you want to add a PHP service, you may run:

```bash
$ aenthill start theaentmachine/aent-php
```

## Upgrade

You may automatically upgrade your version of Aenthill with the latest release thanks to the `upgrade` command.

```bash
$ aenthill upgrade
```

You may also provide a specific version thanks to the `-t --target` flag.

```bash
$ aenthill upgrade -t version
```