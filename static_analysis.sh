#!/bin/sh


PKG="./..."


# repository root
cd $(dirname $0)


# go vet [-n] [-x] [build-flags] [vet-flags] [pkgs]
#   -n      # prints commands that would be executed
#   -x      # prints commands as they are executed
#
#
cmd=" go vet "
#cmd=" $cmd -x "              # -x, prints commands as they are executed
$cmd $PKG

