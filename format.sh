#!/bin/sh


PKG="./..."


# repository root
cd $(dirname $0)

#go fix $PKG
go fmt $PKG
