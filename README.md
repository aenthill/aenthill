# Goals

1. Helping setting up a project on Docker.
2. Easy to use.

# Commands

`anthill init`

1. If the `anthill.yml` does not exist, creates it.
2. Asks for project name.
3. Asks for project description.
4. Asks for project author(s).
5. Asks for project license.

Example of `anthill.yml` file:

```yaml
project: project name
description: project description
authors:
    - Julien Neuhart <j.neuhart@thecodingmachine.com>
license: MIT
```

---

`anthill add ant_name`

* `ant_name` = Docker image with a manifest.

Example of *Dockerfile* with a manifest:

```
FROM alpine:3.7

MAINTAINER Julien Neuhart <j.neuhart@thecodingmachine.com>

LABEL ANT_MANIFEST_PATH="ant-manifest.yml"
```

