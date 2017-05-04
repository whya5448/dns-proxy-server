**Obs**: The version 2 is in development(the main features are working but in beta) and this documentation is being built

# Features
dns-proxy-server is a end user(developers, Server Administrators) DNS server tool with some extra features like:

* Solve names from local configuration database
* Solve names from docker containers using docker **hostname** option or **HOSTNAMES** env
* Solve names from a list of configured DNS servers(as a proxy) if no answer of two above
* [Graphic interface to manage it](http:/127.0.0.1:5380/static/)
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

# Version 2 Improvements
This tool comes from from nodejs version(1.0), improving:
* Performance - this version uses much less RAM and is much faster
* Bug fixes
* Binary distribution - now you can simply download a linux executable and use it, without need to install anything
* Code design quality
* And more

# Installing from binary

>~~not documented yet~~

# Running it

### Default start server
 
	$ ./dns-proxy-server

### Help view

	$ ./dns-proxy-server --help


# Testing

	$ go test -cover=false ./src/github.com/mageddo/dns-proxy-server/.../

# Building from source

clone the repository

	git submodule init && git submodule update

compile the packages

	$ docker rm -f docker-dns-server-compiler && \
	docker-compose up prod-build-dns-proxy-server

Installing

>~~not documented yet~~

# Configuration and data

```javascript
{
  "remoteDnsServers": [], // not used
  "envs": [ // there areall possible environments 
    {
      "name": "", // empty string is the default
      "hostnames": [ // there are all local hostnames entries
        {
          "id": 1,
          "hostname": "github.com",
          "ip": [192, 168, 0, 1],
          "ttl": 255
        }
      ]
    }
  ],
  "activeEnv": "", // what is default env name 
  "lastId": 1, // hostnames sequence
  "webServerPort": 0, // web admin port, when 0 the default value is used
  "dnsServerPort": 8980 // dns server port, when 0 the default value is used
}
```

# Commandline help Documentation

```
-compress
  compress replies
-conf-path string
  The config file path  (default "conf/config.json")
-cpuprofile string
  write cpu profile to file
-default-dns
  This DNS server will be the default server for this machine (default true)
-help
  This message
-server-port int
  The DNS server to start into (default 53)
-tsig string
  use MD5 hmac tsig: keyname:base64
-web-server-port int
  The web server port (default 5380)
```

# Test hostnames

	nslookup -port=8980 bookmarks-node.mageddo.in 127.0.0.1