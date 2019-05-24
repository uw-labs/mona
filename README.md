# mona

Mona is a command-line tool for managing monorepos and is intended for use in CI pipelines. Each application in your monorepo is declared in a `mona.yml` file as a module. In turn, each application has a `module.yml` file that specifies commands for testing/building the application. Changes to these modules are tracked in a `mona.lock` file which should be cached and shared amonst your builds.

## Getting Started

Below are steps on how to install mona and set up your first mona project

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
  build: ""     # Command to run to build the module
  test: ""      # Command to run to test the module
excludes:       # File patterns to ignore
  - "*.txt"
```

In here you can specify how to build/test your module. Note that this must be a single line command, so it is recommended to pass responsibility for building/testing to a `Makefile` or other script runner.

You can also provide file patterns to ignore from hash generation using the `excludes` key.

### Checking for changes

Mona generates hashes to determine if code within respective modules has changed. You can list what has changed using the `mona diff` command:

```bash
mona diff

# Output:
# Modules to be built:
# module

# Modules to be tested:
# module
```

### Building/Testing modules

You can build/test your modified modules using the `mona build` and `mona test` commands. Build and test hashes are stored seperately so mona will know if a built module has also been tested and vice-versa. Any subsequent output to `stdin` or `stderr` will be written to the console and mona will stop executing if your test/build commands return an error exit code.