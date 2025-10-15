# Pacis

[See the documentation](https://ui.canpacis.com)

## Introduction

Server-side rendering is back in style. It's fast, it's good for SEO, and users see content right away instead of waiting for JavaScript to load. But most SSR tools lock you into JavaScript and come with a ton of build tools you have to learn and maintain. This library lets you build server-rendered web apps with Go on the backend and modern frontend tooling on the frontend. Your routes, your business logic, your server rendering, all Go. Your frontend assets get bundled with Bun and Vite, so you get the best of both worlds without fighting two separate ecosystems. No Node runtime required to serve your app. Just one binary that includes everything.

Building web UIs in Go has always been clunky. You were stuck with basic templating or reinventing the wheel with string concatenation. This library changes that. You get composable components that feel natural to use, type-safe templates that catch mistakes before they hit production, and a development experience that doesn't waste your time. Hot reload means you see changes instantly. Streaming means your pages load faster. You can add caching where it actually matters instead of opt-out strategies from other frameworks.

Deployment is where this approach really shines. You bundle your frontend assets during build, then compile everything into one binary with your Go server. No runtime dependencies, no separate frontend and backend deployments, no configuration drift between dev and production. It deploys anywhere, cloud VM, Docker, Kubernetes, edge platforms. You get the safety of compiled code and the simplicity of one artifact to ship.

## Installation

To install, and get started, your best bet is creating a new repository from Pacis' [Github Template](https://github.com/canpacis/pacis-template). Go ahead create yours and clone it to your local environment.

```shell
git clone <your-repository>
```
You can follow the simple readme guide in the template repository or follow this one.

### Install Dependencies

Install your go and bun dependencies. 

> Pacis uses bun and vite for bundling your frontend assets. You can replace it with another package manager but the template heavily relies on it.

```shell
bun install
```

```shell
go mod download
```

Pacis also uses [air](https://github.com/air-verse/air) and [taskfile](https://taskfile.dev/) to streamline development and build processes. If you don't have them already, follow their installation guides:

- [Air](https://github.com/air-verse/air?tab=readme-ov-file#installation)
- [Taskfile](https://taskfile.dev/docs/installation)

### Development

You are ready to run your development server. There is a task called `dev` already defined for you in the `Taskfile.yml` file. 

```shell
task dev
```

Running this command will spin up a server on localhost port 8080. It will also run a vite dev server on port 5173.

> The dev server actually runs on port 8081 but it is proxied to port 8080 with hot reloading using air. Check the `.air.toml` file for its configuration.