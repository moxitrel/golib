#!/bin/sh


PKG="./..."


# repository root
cd $(dirname $0)

go fmt $PKG
go fix $PKG
