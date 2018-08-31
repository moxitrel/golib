#!/bin/sh


export GOARCH=${GOARCH:-$(go env GOARCH)}   # amd64 | 386 | arm64 | arm | ...
export GOOS=${GOOS:-$(go env GOOS)}         # linux | darwin | windows | freebsd | netbsd | openbsd | ...
export GOROOT=${GOROOT:-$(go env GOROOT)}
export GOPATH=${GOPATH:-$(go env GOPATH)}

PKG="./..."


# repository root
cd $(dirname $0)


# go vet [-n] [-x] [build-flags] [vet-flags] [pkgs]
#   -n      # prints commands that would be executed
#   -x      # prints commands as they are executed
#
#
vet=" go vet "
#vet=" $vet -x "              # -x, prints commands as they are executed
vet=" $vet $PKG "


build=" go build "
#build=" $build -v "              # -v, print the names of packages as they are compiled
#build=" $build -x "              # -x, print the commands
#build=" $build -work "           # -work, print tmp work directory and don't delete when exiting.
build=" $build $PKG "


$vet && $build