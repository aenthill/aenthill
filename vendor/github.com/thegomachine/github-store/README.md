<p align="center">
    <img src="https://user-images.githubusercontent.com/8983173/42808440-9fb27048-89b3-11e8-9be1-22b1031dec9a.png" alt="Logo" width="200" height="200" />
</p>
<h3 align="center">github-store</h3>
<p align="center">A simple package to be used as a GitHub store for <a href="https://github.com/tj/go-update">go-update</a> which does not prefix release tag with "v"</p>
<p align="center">
    <a href="https://godoc.org/github.com/thegomachine/github-store">
        <img src="https://godoc.org/github.com/thegomachine/github-store?status.svg" alt="GoDoc">
    </a>
</p>

---

[go-update](https://github.com/tj/go-update) is a go package for auto-updating system-specific binaries via GitHub releases.

The current store used by this package automatically prefix the release tag with "v". 

Package `thegomachine/github-store` provides the same API but without prefixing release tag.

## Installation

```bash
$ go get github.com/thegomachine/github-store
```

## Usage

So previously your code might have looked like this:

```golang
import (
    "github.com/tj/go-update/stores/github"
    update "github.com/tj/go-update"
)

func main() {
    project := &update.Manager{
		Command: "foo",
		Store: &github.Store{
			Owner:   "foo",
			Repo:    "bar",
			Version: version,
		},
    }
    release, err := project.Store.GetRelease("1.0.0")
    // will try to find the GitHub release v1.0.0.
 }
```

With this package:

```golang
import (
    "github.com/thegomachine/github-store"
    update "github.com/tj/go-update"
)

func main() {
    project := &update.Manager{
		Command: "foo",
		Store: &github.Store{
			Owner:   "foo",
			Repo:    "bar",
			Version: version,
		},
    }
    release, err := project.Store.GetRelease("1.0.0")
    // will try to find the GitHub release 1.0.0.
 }
```