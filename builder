#!/bin/sh

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

		rm -rf build/ && \
		mkdir -p build/ && \
		git submodule init && \
		git submodule update && \
		cd src && \
		go test -cover=false ./github.com/mageddo/dns-proxy-server/.../ && \
		go build -v -o -ldflags "-X github.com/mageddo/dns-proxy-server/flags.version=`cat VERSION`" ../build/dns-proxy-server && \
		cp -r ../static ../build/ && \
		cp ../docker-compose.yml ../build/docker-compose.yml && \
		cp ../dns-proxy-service ../build/dns-proxy-service && \
		cd ../build/ && \
		tar -cvf dns-proxy-server-2.0.19.tgz * && \
		cd ../

		echo "build success"

	;;

esac