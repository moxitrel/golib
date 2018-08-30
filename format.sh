#!/bin/sh


PKG="./..."

# repository root
cd $(dirname $0)


### Format Codes
#
# go fmt [-n] [-x] [packages]   #= gofmt -l -w ...
#   -n      # print commands that would be executed, dry-run
#   -x      # print commands as they are executed
#
# gofmt [opts] [path ...]
#   -r <string>     # rewrite rule (e.g., 'a[b:len(a)] -> a[b:]')
#   -s	            # simplify code
#   -e	            # report all errors
#   -l	            # list files to be formatted
#   -d	            # show diffs
#   -w              # write to (source) file instead of stdout
#
cmd=" go fmt "
#cmd=" $cmd -x "     # -x, print commands as they are executed

$cmd $PKG
