---
title: "Introduction"
weight: 1
kind: "use"
menu:
  main:
    identifier: "use"
    name: "Use"
    url: "/use"
    weight: 2
---

# Aenthill

Aenthill is a command-line tool that helps bootstraping your [Docker](https://www.docker.com/) projects easily.

Using Aenthill, in a few minutes, you can have:

- Your containers in a `docker-compose.yml` file ready
- ... along with the [Traefik](https://traefik.io/) reverse-proxy to access your web containers
- `Dockerfiles` to build your production containers
- ... and [Kubernetes](https://kubernetes.io/) deployment files to deploy your project
- CI/CD integration to build your containers and deploy them
- ... in test environments or in production!

## How does it work?

Most project scaffolders rely on a list of template files that are used to 
generate a project. But this approach is fundamentally limited.

Aenthill **does not** work like that. Instead, Aenthill relies on a set of
programs hosted in separate containers. We call those programs *aents*.

Working together, these aents will build your project infrastructure.
This architecture is very flexible. It means anyone can write its own aent
and extend the system.

## How does it compare to Helm?

Helm is a Kubernetes only tool that can be used to make Kubernetes deployments
reusable. You typically will write Helm deployments if you are a product owner 
helping clients deploy your application. Writing a Helm deployment requires a 
good understanding of Kubernetes.

Aenthill on the other end is targetted at web developers that are not expert 
DevOps. It helps them starting a web-application and encompass the whole toolset
needed (from docker-compose to image building, CI/CD, etc...).

# TheAentMachine

[TheAentMachine](https://github.com/theaentmachine/) is a *colony* of aents which are able to work which each other.

We provide three kinds of aents:

- Orchestrators (e.g. Docker, Kubernetes)
- CI/CD (GitLab)
- Services (PHP, MySQL, Node.js etc.)

TODO explain briefly how they work together