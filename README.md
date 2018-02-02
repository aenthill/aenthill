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

`anthill add ant_image_name`

* `ant_image_name` = Docker image with a manifest.

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
4. Once done, updates the `anthill.yml` file:

```yaml
project: project name
description: project description
authors:
    - Julien Neuhart <j.neuhart@thecodingmachine.com>
license: MIT
ant-builders:
  - dev:
    - local_name: ant-docker-compose-builder-1
      ant_name: anthill-docker/ant-docker-compose-builder
      image: anthill-docker/ant-docker-compose-builder:1.0
      version: 1.0
```

5. Updates the working tree:

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
type: ant-builder
name: anthill-docker/ant-docker-compose-builder
version: 1
workingdir: build # once all events have been handled, `anthill` will put the generated files/folders from this directory to `anthill/ant-builders/dev/ant-docker-compose-builder-1`.
handle:
    - on: NEW_SERVICE
      do: 
        - ant-docker-compose-builder add-service # will ask user for service name and on which networks he wants to attach it
        - anthill notify NEW_CONTAINER --optional --service-name={{ notification.SERVICE_NAME }}
    - on: NEW_IMAGE:
      do:
        - ant-docker-compose-builder add-image --service-name={{ notification.SERVICE_NAME }}
    - on: NEW_CONTAINER
      do:
        - ant-docker-compose-builder add-container-name --service-name={{ notification.SERVICE_NAME }} # will ask user for container name
    - on: NEW_NETWORK
      do:
        - ant-docker-compose-builder add-network [--default-network-name=networkName] # will ask user for network name and which services he want to attach it
    - on: NEW_ENV
      do:
        - ant-docker-compose-builder add-env --service-name={{ notification.SERVICE_NAME }} --env-key-name={{ notification.ENV_KEY_NAME }} [--value-pattern={{ notification.VALUE_PATTERN }} --default-value={{ notification.DEFAULT_VALUE }} --value={{ notification.VALUE }}] # will ask user for env value if --value has not been specified
    - on: NEW_LABEL
      do:
        - ant-docker-compose-builder add-env --service-name={{ notification.SERVICE_NAME }} --label-key-name={{ notification.LABEL_KEY_NAME }} [--value-pattern={{ notification.VALUE_PATTERN }} --default-value={{ notification.DEFAULT_VALUE }} --value={{ notification.VALUE }}] # will ask user for label value if --value has not been specified
    - on: NEW_HTTP_PORT
      do:
        - ant-docker-compose-builder add-http-port --service-name={{ notification.SERVICE_NAME }} --container-http-port:80 [--default-host-http-port=80] # will ask user for port from host to map
    - on: NEW_PORT
      do:
        - ant-docker-compose-builder add-port --service-name={{ notification.SERVICE_NAME }} --container-port:3306 [--default-host-port=3306] # will ask user for port from host to map
    - on: NEW_VOLUME
      do:
        - ant-docker-compose-builder add-volume --service-name={{ notification.SERVICE_NAME }} # will ask user for volume name
    - on: NEW_MAPPED_VOLUME
      do:
        - ant-docker-compose-builder add-mapped-volume --service-name={{ notification.SERVICE_NAME }} --container-path={{ notification.CONTAINER_PATH }} [--host-path=hostPath] # will ask user for host path to map --host-path has not been specified
    - ...
```

## ant

1. When added, let the user select an environment.
2. If the environment does not exist, asks for a environment name.
3. Asks for the ant name with a default value (ex: ant-traefik-1).
4. Once done, updates the `anthill.yml` file:

```yaml
project: project name
description: project description
authors:
    - Julien Neuhart <j.neuhart@thecodingmachine.com>
license: MIT
ant-builders:
  - dev:
    - local_name: ant-docker-compose-builder-1
      ant_name: anthill-docker/ant-docker-compose-builder
      image: anthill-docker/ant-docker-compose-builder:1.0
      version: 1.0
ants:
  - dev:
    - local_name: ant-traefik-1
      ant_name: anthill-docker/ant-traefik
      image: anthill-docker/ant-traefik:1.5
      version: 1.0
```

5. Updates the working tree:

```
.
+-- anthill
|   +-- ant-builders
|   |   +-- dev
|   |   |   +-- ant-docker-compose-builder-1
|   +-- ants
|   |   +-- dev
|   |   |   +-- ant-traefik-1
+-- anthill.yml
```

### Example: ant-traefik

* `ant-manifest.yml`:

```yaml
type: ant
name: anthill-docker/ant-traefik
version: 1
notify:
    - anthill notify NEW_SERVICE
    - anthill notify NEW_IMAGE --image-name=anthill-docker/ant-traefik:1.15
    - anthill notify NEW_NETWORK --default-network-name=frontend
    - custom command # should returns to stdout a list of notifications or exit 0 
handle:
    - on: NEW_HTTP_PORT
      do:
        - add-virtual-host.sh {{ notification.SERVICE_NAME }} # should returns to stdout a list of notifications like NEW_LABEL with mapped port + virtualhost
        - anthill notify TRAEFIK_ACTIVATE_LOCAL_HTTPS --service-name={{ notification.SERVICE_NAME }} --virtual-host={{ CONTAINER_ENV.VIRTUAL_HOST }}
        - anthill notify REMOVE_HTTP_PORT --service-name={{ notification.SERVICE_NAME }}
    - on: TRAEFIK_ACTIVATE_LOCAL_HTTPS
      do:
      # ask the user if he wants to generate a self-signed-certificate. 
      # if true, will add it to the workingdir.
      # will also returns to stdout "anthill notify NEW_MAPPED_VOLUME --serviceName={{ notification.SERVICE_NAME }} --container-path=build/certs --host-path={{ ant.HOST_PATH }}/certs"
      # if fale, returns to stdout "anthill notify TRAEFIK_ACTIVATE_ACME"
      - add-local-https.sh --virtual-host={{ CONTAINER_ENV.VIRTUAL_HOST }} 
    - on: TRAEFIK_ACTIVATE_ACME
      do:
      - ...
```


---

`anthill remove ant_name`

---

`anthill rollback`