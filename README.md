# mona

[![GoDoc](https://godoc.org/github.com/davidsbond/mona?status.svg)](http://godoc.org/github.com/davidsbond/mona)
[![CircleCI](https://circleci.com/gh/davidsbond/mona/tree/master.svg?style=shield)](https://circleci.com/gh/davidsbond/mona/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/davidsbond/mona)](https://goreportcard.com/report/github.com/davidsbond/mona)
[![Docker Pulls](https://img.shields.io/docker/pulls/davidsbond/mona.svg)](https://hub.docker.com/r/davidsbond/mona)

Mona is a command-line tool for managing monorepos and is intended for use in CI pipelines. Each independent part of your repository is considered an app than can be built, tested or linted. Each app has a respective `app.yml` file with information on commands to run, artifacts to store and files to exclude when generating hashes. The `mona.yml` file is used at the root of the project. App changes are stored in a `mona.lock` file that should be shared across your builds to ensure they're treated incrementally.

A typical project structure may look like this:

```bash
.
├── src
│   ├── api
│   │   ├── main.go
│   │   └── app.yml
│   └── ui
│       ├── main.go
│       └── app.yml
├── mona.lock
└── mona.yml

```

## Getting Started

Below are steps on how to install mona and set up your first mona project.

### Installation

You can download the latest version of mona from the [releases](https://github.com/davidsbond/mona/releases) and place it anywhere
in your `$PATH`.

Alternatively, you can download and compile the source code using `go get` and `make`

```bash
$ go get -u github.com/davidsbond/mona
$ cd $GOPATH/src/github.com/davidsbond/mona
$ make build
>
```

This will generate a `dist` directory in the source folder containing the compiled binary with the `README.md` and `LICENSE` files. These can be placed anywhere in your `$PATH`

### Creating a project

Navigate to your monorepo and run the `mona init` command. You will see a `mona.yml` and `mona.lock` file have been generated for you.

The `mona.yml` file defines your project and is used to signify the root of the monorepo. It is also used to define global file matchers for hash generation:

```yaml
# mona.yml
name: project # The name of the project
exclude:      # File patterns to ignore across all apps
  - "*.exe"
```

### Adding apps

Each application/library in your monorepo is represented as a mona app. You can add these to your project by using the `mona add-app` command like so:

```bash
$ mona add-app apps/ui
• Created new app ui at apps/ui
```

Within your new app, you'll find a `app.yml` file with several keys:

```yaml
# app.yml
name: app
commands:
  build:
    run: "make build"                      # Command to run to build the app
    image: docker.io/library/golang:alpine # Image to execute the build command in
  test:
    run: "make build"                      # Command to run to test the app
    image: docker.io/library/golang:alpine # Image to execute the test command in
  lint:
    run: "make build"                      # Command to run to lint the app
    image: docker.io/library/golang:alpine # Image to execute the lint command in
exclude:                                   # File patterns to ignore
  - "*.txt"
```

In here you can specify how to build/test your app. Note that this must be a single line command, so it is recommended to pass responsibility for building/testing to a `Makefile` or other script runner.

You can also provide file patterns to ignore from hash generation using the `exclude` key.

### Checking for changes

Mona generates hashes to determine if code within respective apps has changed. You can list what has changed using the `mona diff` command:

```bash
$ mona diff
• 1 app(s) to be built  
• 1 app(s) to be tested
• 3 app(s) to be linted
```

### Building/Testing/Linting apps

You can build, test & lint your modified apps using the `mona build`, `mona test` & `mona lint` commands. Individual hashes are stored separately so mona will know if a app has been built, but not linted or tested etc.

If you provide an image for one of these commands, the app will be mounted to a docker container of your choice and execute your command from that volume's location. If the `image` key is left blank, the commands are ran on the host machine.

The volume is created with a [bind mount](https://docs.docker.com/storage/bind-mounts/) so any build artefacts will be available on the host machine

### Roadmap

Currently, mona is **not production ready**. The project is missing quite a few tests and has not been battle tested. That being said, attempts to use this tool are welcome and all feedback is appreciated. You can see the [contributing guide](CONTRIBUTING.md) for more information.
