#!/bin/sh


export GOARCH=${GOARCH:-$(go env GOARCH)}   # amd64 | 386 | arm64 | arm | ...
export GOOS=${GOOS:-$(go env GOOS)}         # linux | darwin | windows | freebsd | netbsd | openbsd | ...
export GOROOT=${GOROOT:-$(go env GOROOT)}
export GOPATH=${GOPATH:-$(go env GOPATH)}

PKG="./..."


# repository root
cd $(dirname $0)


go vet   $PKG
go build $PKG   # whether can pass build
