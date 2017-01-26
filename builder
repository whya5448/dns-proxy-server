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

esac