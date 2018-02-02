# Goals

1. Helping setting up a project on Docker.
2. Easy to use.

# Commands

`anthill init`

1. If the `anthill.yml` fiel does not exist, creates it.
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

There are two types of manifest:

1. `ant-builder`: container which will handle events for generating files/folders
2. `ant`: container which will send events to builders and will also handle event from others ant.

## ant-builder

1. When added, asks for an environment
2. Once the environment is specified by the user, updates the `anthill.yml` file:

```yaml
project: project name
description: project description
authors:
    - Julien Neuhart <j.neuhart@thecodingmachine.com>
license: MIT
ant-builders:
  - dev:
    - ant-builder-1
```

3. Updates the working tree:

```
.
+-- anthill
|   +-- ant-builders
|   |   +-- dev
|   |   |   +-- ant-builder-1
+-- anthill.yml
```

### Example: ant-docker-compose-builder

* `ant-manifest.yml`:

```yaml
type: ant-builder
handle:
    - for: NEW_SERVICE
      do: 
        - ant-docker-compose-builder add-service --image-name=imageName # will ask user for service name and on which networks he wants to attach it
        - send NEW_CONTAINER --optional
    - for: NEW_CONTAINER
      do:
        - ant-docker-compose-builder add-container # will ask user for container name
    - for: NEW_NETWORK
      do:
        - ant-docker-compose-builder add-network [--default-network-name=networkName] # will ask user for network name and which services he want to attach it
    - for: NEW_ENV
      do:
        - ant-docker-compose-builder add-env --env-key-name=envKeyName [--value-pattern=[0-9]* --default-value=defaultValue] # will ask user for env value
    - for: NEW_LABEL
      do:
        - ant-docker-compose-builder add-env --label-key-name=labelKeyName [--value-pattern=[0-9]* --default-value=defaultValue] # will ask user for label value
    - for: NEW_HTTP_PORT
      do:
        - ant-docker-compose-builder add-http-port --container-http-port:80 [--default-host-http-port=80] # will ask user for port from host to map
    - for: NEW_PORT
      do:
        - ant-docker-compose-builder add-port --container-port:3306 [--default-host-port=3306] # will ask user for port from host to map
    - for: NEW_VOLUME
      do:
        - ant-docker-compose-builder add-volume # will ask user for volume name
    - for: NEW_MAPPED_VOLUME
      do:
        - ant-docker-compose-builder add-mapped-volume --container-path=containerPath # will ask user for host path to map 
```

## ant

---

`anthill remove ant_name`