#!/usr/bin/env bash

set -euo pipefail

if [[ -d $PWD/go-module-cache && ! -d ${GOPATH}/pkg/mod ]]; then
  mkdir -p ${GOPATH}/pkg
  ln -s $PWD/go-module-cache ${GOPATH}/pkg/mod
fi

commit() {
  git commit -a -m "Dependency Upgrade: $1 $2" || true
}

cd "$(dirname "${BASH_SOURCE[0]}")/.."
git config --local user.name "$GIT_USER_NAME"
git config --local user.email ${GIT_USER_EMAIL}

go build -ldflags='-s -w' -o bin/dependency github.com/cloudfoundry/libcfbuildpack/dependency

bin/dependency mariadb-jdbc "[\d]+\.[\d]+\.[\d]+" $(cat ../mariadb-jdbc/version) $(cat ../mariadb-jdbc/uri)  $(cat ../mariadb-jdbc/sha256)
commit mariadb-jdbc $(cat ../mariadb-jdbc/version)

bin/dependency postgresql-jdbc "[\d]+\.[\d]+\.[\d]+" $(cat ../postgresql-jdbc/version) $(cat ../postgresql-jdbc/uri)  $(cat ../postgresql-jdbc/sha256)
commit postgresql-jdbc $(cat ../postgresql-jdbc/version)
