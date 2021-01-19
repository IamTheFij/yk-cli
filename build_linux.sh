#! /bin/sh
set -e

# Use first arg as docker image to build inside of (ubuntu or golang)
if [ $# -gt 0 ]; then
    docker run --rm -it --workdir /src -v "$(pwd):/src" "$1" ./build_linux.sh
fi

set -x

apt-get update

# For ubuntu
# apt-get install wget
if ! command -v go; then
    apt-get install -y --no-install-recommends make golang git ca-certificates
    mkdir /go
    export GOPATH=/go
    export PATH=$GOPATH/bin:$PATH
fi

# For ubuntu and golang
apt-get install -y --no-install-recommends libpcsclite-dev # gnome-keyring

make dist/yk-linux-amd64
