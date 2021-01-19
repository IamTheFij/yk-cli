# yk-cli

Yubikey cli for retrieving TOTP codes

## Installation

Currently, there is no binary release published, but it can be built with `make build` and then you can copy `./build/yk` to somewhere in your path.

If you're a [`fish`](https://fishshell.com) user, you can also add `./scripts/yk.fish` to your `conf.d/` directory to get completions.

## Building

Executing `make build` will compile to `./build/yk`. Additionally, distribution builds should be possible with `make all` or by building a particular target. Eg `make ./dist/yk-darwin-amd64`. There is also an alias present and the `./dist/` prefix can be left off.

### Note on distribution builds

Currently cross compiling is not working correctly.

### Building for linux distros

This is a work in progress, but it can be done by running `./build_linux.sh [golang|ubuntu]`.
