# Commands

`anthill init`

1. If the `anthill.yml` file does not exist, creates it.
2. Asks for project name.
3. Asks for project description.
4. Asks for project author(s).
5. Asks for project license.

Example of `anthill.yml` file:

```yaml
project: "project name"
description: "project description"
authors:
    - "Julien Neuhart <j.neuhart@thecodingmachine.com>"
license: "MIT"
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

## Example: ant-docker-compose-builder

1. Asks for the builder name with a default value (ex: ant-docker-composer-builder-1).
2. Once done, updates the `anthill.yml` file:

```yaml
project: "project name"
description: "project description"
authors:
    - "Julien Neuhart <j.neuhart@thecodingmachine.com>"
license: "MIT"
ants:
    - local: "ant-docker-compose-builder-1"
      ant: "anthill-docker/ant-docker-compose-builder"
      image: "anthill-docker/ant-docker-compose-builder:1.0"
      version: "1.0.0"
```

3. Updates the working tree:

```
.
+-- ants
|   +-- ant-docker-compose-builder-1
+-- anthill.yml
```

* proposal: `ant-manifest.yml`

```yaml
name: "anthill-docker/ant-docker-compose-builder"
version: "1.0.0"

events_handle:
  - name: "ADD_SERVICE"
    script:
      - ant-docker-compose-builder add-service
  - name: "REMOVE_SERVICE"
    script:
      - ant-docker-compose-builder remove-service --service-name={{ .Event.SERVICE_NAME }}
  - name: "ADD_IMAGE"
    script:
      - ant-docker-compose-builder add-image --service-name={{ .Event.SERVICE_NAME }}
  - name: "REMOVE_IMAGE"
    script:
      - ant-docker-compose-builder remove-image --service-name={{ .Event.SERVICE_NAME }}
  - name: "REMOVE_IMAGE"
    script:
      - ant-docker-compose-builder remove-image --service-name={{ .Event.SERVICE_NAME }}
  - name: "ADD_NETWORK"
    script:
      - ant-docker-compose-builder add-network --default-network-name={{ .Event.DEFAULT_NETWORK_NAME }}
  - name: "REMOVE_NETWORK"
    script:
      - ant-docker-compose-builder remove-network --network-name={{ .Event.NETWORK_NAME }}
  - name: "ADD_ENVIRONMENT"
    script:
      - ant-docker-compose-builder add-environment --service-name={{ .Event.SERVICE_NAME }} --environment-key-name={{ .Event.ENVIRONMENT_KEY_NAME }} --value-pattern={{ .Event.VALUE_PATTERN }} --default-value={{ .Event.DEFAULT_VALUE }} --value={{ .Event.VALUE }}
  - name: "REMOVE_ENVIRONMENT"
    script:
      - ant-docker-compose-builder remove-environment --service-name={{ .Event.SERVICE_NAME }} --environment-key-name={{ .Event.ENVIRONMENT_KEY_NAME }}
  - name: "ADD_LABEL"
    script:
      - ant-docker-compose-builder add-label --service-name={{ .Event.SERVICE_NAME }} --label={{ .Event.LABEL }}
  - name: "REMOVE_LABEL"
    script:
      - ant-docker-compose-builder remove-label --service-name={{ .Event.SERVICE_NAME }} --label={{ .Event.LABEL }}
  - name: "ADD_PORT"
    script:
      - ant-docker-compose-builder add-port --service-name={{ .Event.SERVICE_NAME }} --container-port={{ .Event.CONTAINER_PORT }}
  - name: "REMOVE_PORT"
    script:
      - ant-docker-compose-builder remove-port --service-name={{ .Event.SERVICE_NAME }} --container-port={{ .Event.CONTAINER_PORT }}
  - name: "ADD_HTTP_PORT"
    script:
      - ant-docker-compose-builder add-port --service-name={{ .Event.SERVICE_NAME }} --container-port={{ .Event.CONTAINER_PORT }}
  - name: "REMOVE_HTTP_PORT"
    script:
      - ant-docker-compose-builder remove-port --service-name={{ .Event.SERVICE_NAME }} --container-port={{ .Event.CONTAINER_PORT }}
  - name: "ADD_VOLUME"
    script:
      - ant-docker-compose-builder add-volume --service-name={{ .Event.SERVICE_NAME }} --container-path={{ .Event.CONTAINER_PATH }}
  - name: "REMOVE_VOLUME"
    script:
      - ant-docker-compose-builder remove-volume --service-name={{ .Event.SERVICE_NAME }} --container-path={{ .Event.CONTAINER_PATH }}
  - name: "ADD_MAPPED_VOLUME":
    script:
      - ant-docker-compose-builder add-mapped-volume --service-name={{ .Event.SERVICE_NAME }} --container-path={{ .Event.CONTAINER_PATH }} --host-path={{ .Anthill.HOST_ROOT_PATH }}
  - name: "REMOVE_MAPPED_VOLUME":
    script:
      - ant-docker-compose-builder remove-mapped-volume --service-name={{ .Event.SERVICE_NAME }} --container-path={{ .Event.CONTAINER_PATH }} --host-path={{ .Anthill.HOST_ROOT_PATH }}
```

## Example: ant-traefik

1. Asks for the ant name with a default value (ex: ant-traefik-1).
2. Once done, updates the `anthill.yml` file:

```yaml
project: "project name"
description: "project description"
authors:
    - "Julien Neuhart <j.neuhart@thecodingmachine.com>"
license: "MIT"
ants:
    - local: "ant-docker-compose-builder-1"
      ant: "anthill-docker/ant-docker-compose-builder"
      image: "anthill-docker/ant-docker-compose-builder:1.0"
      version: "1.0.0"
    - local: "ant-traefik-1"
      ant: "anthill-docker/ant-traefik"
      image: "anthill-docker/ant-traefik:1.0"
      version: "1.0.0"
```

3. Updates the working tree:

```
.
+-- ants
|   +-- ant-docker-compose-builder-1
|   +-- ant-traefik-1
+-- anthill.yml
```

* proposal: `ant-manifest.yml`

```yaml
name: "anthill-docker/ant-traefik"
version: "1.0.0"

events_send:
  - for: "anthill-docker/ant-docker-compose-builder"
    version: "~1.0.0"
    notify:
      - name: "ADD_SERVICE"
        then:
          - name: "ADD_IMAGE"
            variables:
              - SERVICE_NAME={{ .Response.SERVICE_NAME }}
              - IMAGE_NAME="traefik:1.5"
          - name: "ADD_NETWORK"
            variables:
              - SERVICE_NAME={{ .Response.SERVICE_NAME }}
              - DEFAULT_NETWORK_NAME="frontend"
    
events_handle:
  - name: "NEW_HTTP_PORT"
    script:
      - add-virtual-host.sh
    then:
      - notify:
        - name: "ADD_LABEL"
          variables:
            - SERVICE_NAME={{ .Sender.SERVICE_NAME }}
            - LABEL="traefik.frontend.rule=HOST:{{ .Container.Env.VIRTUAL_HOST }}"
      - notify:
        - name: "REMOVE_HTTP_PORT"
          variables:
            - SERVICE_NAME=
                    
```

---

`anthill remove ant_name`

---

`anthill rollback`