Obs: The version 2 is in development and this documentation is being built

# Features
dns-proxy-server is a end user(developers, Server Administrators) DNS server tool with some extra features like:

* Solve names from local configuration database
* Solve names from docker containers using docker **hostname** option or **HOSTNAMES** env
* Solve names from a list of configured DNS servers(as a proxy) if no answer of two above
* Graphic interface to manage it
	* List and edit DNS local entries
	* ~~List docker containers hostnames~~
* ~~Cache for remote DNS increasing internet velocity, and options to enable/disable~~
* ~~List docker containers using [http://dns.mageddo:5380/containers](http://dns.mageddo:5380/containers)~~
* ~~List cached hosts using [http://127.0.0.1:5380/cache](http://127.0.0.1:5380/cache)(without docker) or [http://dns.mageddo:5380/cache](http://dns.mageddo:5380/cache) (with docker)~~

# DNS resolution order
The Dns Proxy Server basically follow the bellow order to solve the names:

* DNS try to solve the hosts from **docker** containers
* then from local database file
* then from 3rd configured remote DNS servers

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