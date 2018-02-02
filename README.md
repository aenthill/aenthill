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

1. `ant-builder`: container which will handle events for generating files/folders.
2. `ant`: container which will send events to builders and will also handle event from others ant.

## ant-builder

1. When added, let the user select an environment.
2. If the environment does not exist, asks for a environment name.
3. Asks for the builder name with a default value (ex: ant-docker-composer-builder-1).
2. Once done, updates the `anthill.yml` file:

```yaml
project: project name
description: project description
authors:
    - Julien Neuhart <j.neuhart@thecodingmachine.com>
license: MIT
ant-builders:
  - dev:
    - ant-docker-compose-builder-1: anthill-docker/ant-docker-compose-builder:1.0
```

3. Updates the working tree:

```
.
+-- anthill
|   +-- ant-builders
|   |   +-- dev
|   |   |   +-- ant-docker-compose-builder-1
+-- anthill.yml
```

### Example: ant-docker-compose-builder

* `ant-manifest.yml`:

```yaml
version: 1
type: ant-builder
workingdir: build # all generated files/folders will be added to the host directory 
handle:
    - on: NEW_SERVICE
      do: 
        - ant-docker-compose-builder add-service --image-name=imageName # will ask user for service name and on which networks he wants to attach it
        - anthill notify NEW_CONTAINER --optional
    - on: NEW_CONTAINER
      do:
        - ant-docker-compose-builder --service-name=serviceName add-container-name # will ask user for container name
    - on: NEW_NETWORK
      do:
        - ant-docker-compose-builder add-network [--default-network-name=networkName] # will ask user for network name and which services he want to attach it
    - on: NEW_ENV
      do:
        - ant-docker-compose-builder add-env --service-name=serviceName --env-key-name=envKeyName [--value-pattern=[0-9]* --default-value=defaultValue --value=value] # will ask user for env value if --value has not been specified
    - on: NEW_LABEL
      do:
        - ant-docker-compose-builder add-env --service-name=serviceName --label-key-name=labelKeyName [--value-pattern=[0-9]* --default-value=defaultValue --value=value] # will ask user for label value if --value has not been specified
    - on: NEW_HTTP_PORT
      do:
        - ant-docker-compose-builder add-http-port --service-name=serviceName --container-http-port:80 [--default-host-http-port=80] # will ask user for port from host to map
    - on: NEW_PORT
      do:
        - ant-docker-compose-builder add-port --service-name=serviceName --container-port:3306 [--default-host-port=3306] # will ask user for port from host to map
    - on: NEW_VOLUME
      do:
        - ant-docker-compose-builder add-volume --service-name=serviceName # will ask user for volume name
    - on: NEW_MAPPED_VOLUME
      do:
        - ant-docker-compose-builder add-mapped-volume --service-name=serviceName --container-path=containerPath # will ask user for host path to map 
```

Once all events have been handled, `anthill` will put the generated files/folders from `workingdir` inside the container to `anthill/ant-builders/dev/ant-docker-compose-builder-1`.

## ant

---

`anthill remove ant_name`