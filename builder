#!/bin/sh

set -e

CUR_DIR=`pwd`
APP_VERSION=$(cat VERSION)
REPO_URL=mageddo/dns-proxy-server

create_release(){

	PAYLOAD=`echo '{
			"tag_name": "VERSION",
			"target_commitish": "TARGET",
			"name": "VERSION",
			"body": "",
			"draft": false,
			"prerelease": true
		}' | sed -e "s/VERSION/$APP_VERSION/" | sed -e "s/TARGET/$TRAVIS_BRANCH/"` && \
	TAG_ID=`curl -i -s -f -X POST "https://api.github.com/repos/$REPO_URL/releases?access_token=$REPO_TOKEN" \
--data "$PAYLOAD" | grep -o -E 'id": [0-9]+'| awk '{print $2}' | head -n 1`
}

upload_file(){
	curl --data-binary "@$SOURCE_FILE" -i -w '\n' -f -s -X POST -H 'Content-Type: application/octet-stream' \
"https://uploads.github.com/repos/$REPO_URL/releases/$TAG_ID/assets?name=$TARGET_FILE&access_token=$REPO_TOKEN"
}

case $1 in

	setup-repository )
		git remote remove origin  && git remote add origin https://${REPO_TOKEN}@github.com/$REPO_URL.git
		git checkout -b build_branch ${TRAVIS_BRANCH}
		echo "> Repository added, travisBranch=${TRAVIS_BRANCH}"

	;;

	upload-release )

		if [ "$REPO_TOKEN" = "" ] ; then echo "REPO_TOKEN cannot be empty"; exit 1; fi

		if [ "`git config user.email || echo ''`" = "" ]; then
			echo '> custom config'
			git config user.name `git config user.name || echo 'CI BOT'`
			git config user.email `git config user.email || echo 'ci-bot@mageddo.com'`
		fi
		echo '> config'
		git config -l
		echo ''

		REMOTE="https://${REPO_TOKEN}@github.com/${REPO_URL}.git"

		git checkout -b build_branch ${CURRENT_BRANCH}
		echo "> Repository added, currentBranch=${CURRENT_BRANCH}"

		git commit -am "Releasing ${APP_VERSION}" # if there is nothing to commit the program will exits
		git tag ${APP_VERSION}
		git push "$REMOTE" "build_branch:${CURRENT_BRANCH}"
		git status
		echo "> Branch pushed - Branch $CURRENT_BRANCH"

		create_release
		echo "> Release created with id $TAG_ID"

		SOURCE_FILE="build/dns-proxy-server-$APP_VERSION.tgz"
		TARGET_FILE=dns-proxy-server-$APP_VERSION.tgz
		echo "> Source file hash"
		md5sum $SOURCE_FILE && ls -lha $SOURCE_FILE

		upload_file

	;;

	apply-version )

		# updating files version
		sed -i -E "s/(dns-proxy-server.*)[0-9]+\.[0-9]+\.[0-9]+/\1$APP_VERSION/" docker-compose.yml
		sed -i -E "s/[0-9]+\.[0-9]+\.[0-9]+/$APP_VERSION/g" Dockerfile.hub

	;;

	build )

		echo "> Starting build"

		rm -rf build/
		mkdir -p build/
		go test -cover=false -ldflags "-X github.com/mageddo/dns-proxy-server/flags.version=test" ./.../
		go build -v -o build/dns-proxy-server -ldflags "-X github.com/mageddo/dns-proxy-server/flags.version=$APP_VERSION"
		cp -r static build/
		cd build/
		tar -czvf dns-proxy-server-${APP_VERSION}.tgz *

		echo "> Build success"

	;;


	validate-release )

		if git rev-parse "$APP_VERSION^{}" >/dev/null 2>&1; then
			echo "> Version already exists $APP_VERSION"
			exit 3
		fi

	;;

esac
