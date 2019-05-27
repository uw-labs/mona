# mona

[![GoDoc](https://godoc.org/github.com/davidsbond/mona?status.svg)](http://godoc.org/github.com/davidsbond/mona)
[![CircleCI](https://circleci.com/gh/davidsbond/mona/tree/master.svg?style=shield)](https://circleci.com/gh/davidsbond/mona/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/davidsbond/mona)](https://goreportcard.com/report/github.com/davidsbond/mona)

Mona is a command-line tool for managing monorepos and is intended for use in CI pipelines. Each application/library in your monorepo is tracked in a `mona.yml` file as a module. In turn, each module has a `module.yml` file that specifies commands for testing/building. Changes to these modules are tracked in a `mona.lock` file which should be cached and shared amongst your builds.

A typical project structure may look like this:

```bash
.
├── src
│   ├── api
│   │   ├── main.go
│   │   └── module.yml
│   └── ui
│       ├── main.go
│       └── module.yml
├── mona.lock
└── mona.yml

```

## Getting Started

Below are steps on how to install mona and set up your first mona project.

### Installation

You can download the latest version of mona from the [releases](https://github.com/davidsbond/mona/releases) and place it anywhere
in your `$PATH`.

Alternatively, you can download and compile the source code using `go get`

```bash
go get github.com/davidsbond/mona
go install github.com/davidsbond/mona
```

### Creating a project

Navigate to your monorepo and run the `mona init` command. You will see a `mona.yml` and `mona.lock` file have been generated for you.

The `mona.yml` file defines your project and locally keeps track of the modules you've added. The `mona.lock` file stores hashes of your module content so it knows when things need building/testing.

### Adding modules

Each application/library in your monorepo is represented as a mona module. You can add these to your project by using the `mona add-module` command like so:

```bash
mona add-module path/to/module
```

Within your new module, you'll find a `module.yml` file with several keys:

```yaml
name: module
commands:
  build: "make build"     # Command to run to build the module
  test: "make test"       # Command to run to test the module
  lint: "make lint"       # Command to run to lint the module
exclude:                  # File patterns to ignore
  - "*.txt"
artefacts:                # Build artefacts to move
  - "dist/mybinary"
```

In here you can specify how to build/test your module. Note that this must be a single line command, so it is recommended to pass responsibility for building/testing to a `Makefile` or other script runner.

You can also provide file patterns to ignore from hash generation using the `exclude` key.

### Checking for changes

Mona generates hashes to determine if code within respective modules has changed. You can list what has changed using the `mona diff` command:

```bash
mona diff

> Modules to be built:
> module
>
> Modules to be tested:
> module
```

### Building/Testing/Linting modules

You can build, test & lint your modified modules using the `mona build`, `mona test` & `mona lint` commands. Individual hashes are stored seperately so mona will know if a module has been built, but not linted or testec etc. Any subsequent output to `stdin` or `stderr` will be written to the console and mona will stop executing if your commands return an error exit code.

### Storing build artefacts

In the `mona.yml` file, you can use the `artefacts` key to specify a directory that all artefacts for each module built will be stored. In each `module.yml` file, you can specify file matchers that will cause those files to be moved to the directory specified in `mona.yml`

```yaml
# mona.yml
name: example
version: v1
modules:
- app-01 apps/app-01
- app-02 apps/app-02
artefacts: "dist"
```

```yaml
# module.yml
name: app-01
commands:
  build: "make build"
  test: "make test"
  lint: "make lint"
artefacts:
  - "app-01"
```