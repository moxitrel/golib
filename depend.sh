#!/bin/sh


export GOPATH=${GOPATH:-$(go env GOPATH)}



# repository root
#cd $(dirname $0)


### dep (deprecated)
#
# vendor目录: go编译时，先从repo根下的vendor目录查找代码，若无，再去 GOPATH 查找
#

# install dep if not exist
# if ! which -s dep; then
#     go get -u github.com/golang/dep/cmd/dep
# fi

# 切换至 REPO 目录
#cd $GOPATH/src/$REPO

# support private repository
#
# adds private key to the authentication agent
# default, ~/.ssh/id_*  ,  ~/.ssh/identity
# ssh-add

# set up a new project at [root]
# dep init [root]   #default, current directory
#   -gopath         #search in GOPATH for dependencies
#   -skip-tools     #skip importing other configurations (glide, godep, ...).
#   -v              #enable verbose logging
# dep init -gopath

# install the project's dependencies
# dep ensure
#   -add <pkg>  #add a dependency
#   -update     #update Gopkg.lock according to Gopkg.toml
#
#   -examples   #print detailed usage examples
#   -v          #enable verbose logging
#   -dry-run    #only report the changes
# dep ensure -update


### go get
#
# go get <repo>/<pkg>
# go get <repo>/...     #all pkgs under repo
#   -u                  #update
#

# link github.com/golang to golang.org/x
mkdir -p $GOPATH/src/github.com/golang
mkdir -p $GOPATH/src/golang.org/
ln -s -f $GOPATH/src/github.com/golang $GOPATH/src/golang.org/x

# golib
go get github.com/emirpasic/gods/...
# logrus
#go get github.com/golang/sys/unix
#go get github.com/golang/crypto/ssh/terminal
#go get github.com/Sirupsen/logrus
## gorm
#go get github.com/lib/pq        # postgres
#go get github.com/jinzhu/gorm
