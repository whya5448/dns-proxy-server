#!/bin/sh

set -e

CUR_DIR=$PWD

case $1 in

	pull-all )
		git pull
		for i in `git submodule | awk '{print $2}'`; do
			MATCH=`echo $i | grep -o "mageddo"`
			MATCH2=`echo $i | grep -o "ElvisDeFreitas"`

				echo "pulling $i"
				cd $i
				git pull
				cd $CUR_DIR

		done;
	;;

	build )

		echo "starting build"
		VERSION=`cat VERSION`

		rm -rf build/ && \
		mkdir -p build/ && \
		git submodule init && \
		git submodule update && \
		cd src && \
		go test -cover=false \
			-ldflags "-X github.com/mageddo/dns-proxy-server/flags.version=test" \
			./github.com/mageddo/dns-proxy-server/.../ && \
		go build -v -o ../build/dns-proxy-server \
			-ldflags "-X github.com/mageddo/dns-proxy-server/flags.version=$VERSION" && \
		cp -r ../static ../build/ && \
		cd ../build/ && \
		tar -cvf dns-proxy-server-$VERSION.tgz * && \
		cd ../

		echo "build success"

	;;

esac