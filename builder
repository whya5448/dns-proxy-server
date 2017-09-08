#!/bin/sh

set -e

CUR_DIR=$PWD

case $1 in

	upload-release )

	APP_VERSION=$(cat VERSION)
	REPO_URL=mageddo/dns-proxy-server
	SOURCE_FILE="build/dns-proxy-server-$APP_VERSION.tgz"
	TARGET_FILE=dns-proxy-server-$APP_VERSION.tgz

	echo "> Source file hash"
	md5sum $SOURCE_FILE && ls -lha $SOURCE_FILE

	git remote remove origin  && git remote add origin https://${REPO_TOKEN}@github.com/$REPO_URL.git
	echo "> Repository added"

	git push origin $TRAVIS_BRANCH
	echo "> Branch pushed"

	PAYLOAD=`echo '{
			"tag_name": "VERSION",
			"target_commitish": "TARGET",
			"name": "VERSION",
			"body": "",
			"draft": false,
			"prerelease": false
		}' | sed -e "s/VERSION/$APP_VERSION/" | sed -e "s/TARGET/$TRAVIS_BRANCH/"` && \
	TAG_ID=`curl -i -s -f -X POST "https://api.github.com/repos/$REPO_URL/releases?access_token=$REPO_TOKEN" \
--data "$PAYLOAD" | grep -o -E 'id": [0-9]+'| awk '{print $2}' | head -n 1`
	echo "> Release created with id $TAG_ID"

	curl --data-binary "@$SOURCE_FILE" -i -w '\n' -f -s -X POST -H 'Content-Type: application/octet-stream' \
"https://uploads.github.com/repos/$REPO_URL/releases/$TAG_ID/assets?name=$TARGET_FILE&access_token=$REPO_TOKEN"

	;;

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

		echo "> Starting build"
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
		tar -czvf dns-proxy-server-$VERSION.tgz * && \
		cd ../

		echo "> Build success"

	;;

esac