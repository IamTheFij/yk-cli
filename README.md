# yk-cli

Yubikey cli for retrieving TOTP codes

## Installation

Currently, there is no binary release published, but it can be built with `make build` and then you can copy `./build/yk` to somewhere in your path.

If you're a [`fish`](https://fishshell.com) user, you can also add `./scripts/yk.fish` to your `conf.d/` directory to get completions.

## Building

Executing `make build` will compile to `./build/yk`. Additionally, distribution builds should be possible with `make all` or by building a particular target. Eg `make ./dist/yk-darwin-amd64`. There is also an alias present and the `./dist/` prefix can be left off.

### Dependencies

This tool depends on libpcsclite to talk to the smart card (scard) interface on your computer. Different distributions will require different steps to install the proper dependencies.

#### Ubuntu

    apt-get install libpcsclite-dev

You may also need to install `gnome-keyring` as well for credential storage.

#### CachyOS (Arch)

    sudo packman -S pcsclite ccid

You may need to activate the service and socket as well. You can check their status using:

    systemctl status pcscd.service pcscd.socket

If they need to be enabled, you can try:

    systemctl activate pcscd.service
    systemctl start pcscd.service

If you still have issues, make sure the socket is enabled and started too.

### Note on distribution builds

Currently cross compiling is not working correctly.

### Building for linux distros

This is a work in progress, but it can be done by running `./build_linux.sh [golang|ubuntu]`.
