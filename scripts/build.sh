#!/usr/bin/env bash

# usage:
#   ./build.sh          # builds jj for the current system architecture
#   ./build.sh package  # builds jj for windows/mac/linux/freebsd

set -e

VERSION="1.9.2"

cd $(dirname "${BASH_SOURCE[0]}")/..

package(){
	echo Packaging $1 Binary
	bdir=jj-${VERSION}-$2-$3
	rm -rf packages/$bdir && mkdir -p packages/$bdir
	GOOS=$2 GOARCH=$3 scripts/build.sh
	if [ "$2" == "windows" ]; then
		mv jj packages/$bdir/jj.exe
	else
		mv jj packages/$bdir
	fi
	cp README.md packages/$bdir
	cd packages
	if [ "$2" == "linux" ]; then
		tar -zcf $bdir.tar.gz $bdir
	else
		zip -r -q $bdir.zip $bdir
	fi
	rm -rf $bdir
	cd ..
}

if [ "$1" == "package" ]; then
	rm -rf packages/
	package "Windows" "windows" "amd64"
	package "Mac" "darwin" "amd64"
	package "Linux" "linux" "amd64"
	package "FreeBSD" "freebsd" "amd64"
	package "OpenBSD" "openbsd" "amd64"
	exit
fi


# build and store objects into original directory.
go build -ldflags "-X main.version=$VERSION" -o jj cmd/jj/*.go

