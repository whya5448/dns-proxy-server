Obs: The version 2 is in development and this documentation is being built

# Introduction
dns-proxy-server is a end user(developers, Server Administrators) DNS server tool with some extra features like:
* Solve names from docker containers using docker **hostname** option or **HOSTNAMES** env
* Solve names from local configuration database
* Solve names from a list of configured DNS servers(as a proxy) if no answer of two above
* Graphic interface to manage it

This tool is in the version 2 from nodejs version, improving:
* Performance - this version uses much less RAM and is much faster
* Bug fixes
* Binary distribution - now you can simply download a linux executable and use it, without need to install anything
* Code design quality
* And more

# Installing

	git submodule init && git submodule update



# Test hostnames

	nslookup -port=8980 bookmarks-node.mageddo.in 127.0.0.1
