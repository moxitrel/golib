#!/bin/sh
#
#
### Install compiles and installs the packages named by the import paths
#
# go install [-i] [build flags] [packages]
#    -i          # install the dependent packages
#
#
#
### Compile the packages, along with their dependencies
#
# go build [-o output] [-i] [build flags] [pkgs]
#    -o output   # save as, only allowed when compiling a single package
#    -i          # install the dependent packages
#
#    -a      # rebuild packages even up-to-date
#    -p <n>  # parallel build, default to CPU number
#
#    -n      # print the commands but do not run them
#    -v      # print the names of packages as they are compiled
#    -x      # print the commands
#    -work   # print the name of the tmp work directory and do not delete it when exiting.
#
#    -race   # enable data race detection.
#            # only for linux/amd64, freebsd/amd64, darwin/amd64 and windows/amd64.
#    -msan   # enable inter-operation with memory sanitizer.
#            # only for linux/amd64 and linux/arm64, with Clang/LLVM as the host C compiler.



export GOARCH=${GOARCH:-$(go env GOARCH)}   # amd64 | 386 | arm64 | arm | ...
export GOOS=${GOOS:-$(go env GOOS)}         # linux | darwin | windows | freebsd | netbsd | openbsd | ...
export GOROOT=${GOROOT:-$(go env GOROOT)}
export GOPATH=${GOPATH:-$(go env GOPATH)}

PKG=""



# repository root
cd $(dirname $0)

cmd=" go build "
#cmd=" $cmd -i "              # -i, install the dependent packages
#cmd=" $cmd -v "              # -v, print the names of packages as they are compiled
#cmd=" $cmd -x "              # -x, print the commands
#cmd=" $cmd -work "           # -work, print tmp work directory and don't delete when exiting.
#cmd=" $cmd -race "           # enable data race detection
#cmd=" $cmd -msan "           # enable memory sanitizer, only linux/amd64 with Clang/LLVM


$cmd $PKG
