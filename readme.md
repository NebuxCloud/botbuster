<p align="center">
    <a href="https://nebux.cloud">
        <picture>
            <source media="(prefers-color-scheme: dark)" srcset="https://nebux.cloud/assets/brand/imagotype_light.svg">
            <img alt="Nebux logo" src="https://nebux.cloud/assets/brand/imagotype_dark.svg" height="60px">
        </picture>
    </a>
</p>

# Nebux Botbuster

A high-performance, self-hosted, proof-of-work captcha server fully compatible with [Altcha](https://altcha.org)'s [front-end widget integrations](https://altcha.org/docs/v2/widget-integration/) for web (native, React, Vue, Svelte, Solid, Lit, Angular) and mobile (Flutter, React Native) applications.

Botbuster gives you full control over your bot protection stack: no third‑party tracking, no usage caps, no external dependencies you don’t own, and no black‑box risk scoring.

Designed as a great alternative to paid, proprietary services such as reCAPTCHA, hCaptcha, Cloudflare Turnstile, FriendlyCaptcha or Altcha Sentinel.

## Features

- 💸 **Proof-of-work.** Instead of frustrating puzzles, it uses invisible challenges that are negligible for humans but costly at scale for all kinds of bots, including AI‑driven ones.
- 🔒 **Security.** With HMAC‑signed challenges, built‑in protection against replay/DDoS attacks, mandatory challenge expiration, and strict CORS enforcement.
- ⚡️ **High performance.** Creating new challenges doesn't require storing data, so a single process can serve **more than 1 million challenges per minute** without breaking a sweat.
- 🔒 **Privacy-first.** Designed to be compliant with GDPR (European Union), PIPEDA/CPPA (Canada), HIPAA (USA), CCPA (California), LGPD (Brazil), DPDPA (India), and PIPL (China).
- 🧑‍🦯‍➡️ **Accessibility.** Fully adhering to WCAG 2.2 AA-level guidelines to leave no one behind.
- 🔌 **Plug and play.** Self-hostable, distributed as a single static binary (less than 10 MiB) and a container image (less than 20 MiB) for ARM64 and AMD64.
- 💾 **Stateless.** Without local persistence and with a single external dependency (Valkey/Redis) for caching, so several replicas can be conveniently run in stateless machines (such as Kubernetes clusters).
- ⚖️ **FOSS.** Completely free and open-source under the GPL-3.0 license.
- ☑️ **Compatibility.** Can be used as a drop‑in back‑end replacement for Altcha without any front‑end changes.

## Installation

### Docker

You can run the Docker image from [GitHub Container Registry](https://github.com/NebuxCloud/botbuster/pkgs/container/botbuster) and the command below.

```Shell
docker run --rm ghcr.io/nebuxcloud/botbuster
```

### Binary

Download the binary from [GitHub Releases](https://github.com/NebuxCloud/botbuster/releases) for Linux (ARM64 or AMD64), and place the executable file in your `PATH`.

### Source

Build from source with Go version 1.18 or higher.

The `go install` command will download the source, compile it, and install the binary in your `$GOBIN` directory.

To install, simply run

```Shell
go install github.com/NebuxCloud/botbuster@latest
```

### Helm

This project can easily be deployed to Kubernetes with [Nebux' generic Helm chart](https://github.com/NebuxCloud/helm-charts/tree/main/charts/generic).

See the following example values:

```yaml
workloads:
  default:
    revisionHistoryLimit: 10

    strategy:
      rollingUpdate:
        maxUnavailable: 1
        maxSurge: 1

    containers:
      default:
        image: ghcr.io/nebuxcloud/botbuster:<tag>
        envFrom:
          configMaps:
            - default
          secrets:
            - default
        resources:
          requests:
            cpu: 50m
            memory: 32Mi
          limits:
            cpu: 200m
            memory: 128Mi
        ports:
          - name: http
            containerPort: 8000
        probes:
          readiness:
            httpGet:
              path: /_health
              port: http
            initialDelaySeconds: 5
          liveness:
            httpGet:
              path: /_health
              port: http
            initialDelaySeconds: 5

    autoscaling:
      targetCPUUtilizationPercentage: 90
      replicas:
        min: 2
        max: 6

    disruptionBudget:
      maxUnavailable: 1

    networkPolicy:
      ingress:
        - ports:
            - port: http
          from: []

    service:
      annotations:
        service.kubernetes.io/topology-mode: Auto
      type: ClusterIP
      ports:
        - name: http
          port: 80
          targetPort: http

    securityContext:
      runAsUser: 1000
      runAsGroup: 1000

configMaps:
  default:
    ALLOWED_ORIGINS: https://example.org

secrets:
  default:
    HMAC_KEY: '<secret>'
    VALKEY_URL: 'redis://<user>:<password>@<host>:6379/<database>'

httpRoutes:
  - name: default
    parentRefs:
      - namespace: networking
        name: default
    hostnames:
      - captcha.example.org
    rules:
      - matches:
          - path:
              value: /
        filters:
          - type: ResponseHeaderModifier
            responseHeaderModifier:
              add:
                - name: strict-transport-security
                  value: max-age=31536000; includeSubDomains; preload
                - name: x-robots-tag
                  value: noindex, nofollow
        backendRefs: [{}]
```

## Configuration

The software is configured with environment variables.

| Name              | Description                                                                                                                                                           | Default value      |
|-------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------|
| `DEBUG`           | Whether to print debug logs.                                                                                                                                          | `false`            |
| `LISTEN_PORT`     | Port on which the HTTP server will listen to requests.                                                                                                                | `8000`             |
| `ALLOWED_ORIGINS` | A comma-separated list of origins a cross-domain request can be executed from. See more information [here](https://github.com/rs/cors?tab=readme-ov-file#parameters). | none               |
| `HMAC_KEY`        | **Secret** HMAC key used to create challenges. Generate a new one with  `botbuster generate:key`.                                                                     | none               |
| `EXPIRATION`      | Expiration time for challenges.                                                                                                                                       | `10m` (10 minutes) |
| `VALKEY_URL`      | Valkey or Redis URL for caching used challenges until they expire. It's required to prevent replay attacks. For example: `redis://valkey:6379/42`.                    | none               |
| `VALKEY_PREFIX`   | Prefix to prepend to all Valkey/Redis keys to avoid collisions.                                                                                                       | `botbuster:`       |

## Development

### Directory structure

The project follows the [_de facto_ standard Go project layout](https://github.com/golang-standards/project-layout) with the additions below:

- `Containerfile`, `compose.yml`, `Makefile`, `.dockerignore` and `.env.example` contain the configuration and manifests that define the development and runtime environments with [OCI](https://opencontainers.org) containers and [Compose](https://docs.docker.com/compose).
- `.github` holds the [GitHub Actions](https://github.com/features/actions) CI/CD pipelines.

### Getting started

This project comes with a containerized environment that has everything necessary to work on any platform without having to install dependencies on the developers' machines.

**TL;TR**

```Shell
make
```

#### Requirements

Before starting using the project, make sure that the following dependencies are installed on the machine:

- [Git](https://git-scm.com).
- An [OCI runtime](https://opencontainers.org), like [Docker Desktop](https://www.docker.com/products/docker-desktop/), [Podman Desktop](https://podman.io) or [OrbStack](https://orbstack.dev).
- [Docker Compose](https://docs.docker.com/compose/install/).

It is necessary to install the latest versions before continuing. You may follow the previous links to read the installation instructions.

#### Initializing

First, initialize the project and run the environment.

```Shell
make
```

Then, download third-party dependencies.

```Shell
make deps
```

You may stop the environment by running the following command.

```Shell
make down
```

### Usage

Commands must be run inside the containerized environment by starting a shell in the main container (`make shell`).

#### Running the development server

Run the following command to start the development server:

```Shell
make run
```

> Note that Git is not available in the container, so you should use it from the host machine. It is strongly recommended to use a Git GUI (like [VS Code's](https://code.visualstudio.com/docs/editor/versioncontrol) or [Fork](https://git-fork.com)) instead of the command-line interface.

#### Running tests

To run all automated tests, use the following command.

```Shell
make test
```

#### Debugging

It is possible to debug the software with [Delve](https://github.com/go-delve/delve). To run the application in debug mode, run the command below.

```Shell
make debug
```

For more advanced scenarios, such as debugging tests, you may open a shell in the container and use the Delve CLI directly.

```Shell
make shell
dlv test --listen=:2345 --headless --api-version=2 ./internal/api
```
