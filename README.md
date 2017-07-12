**Obs**: This is the version 2, for the old version 1 [see this link](https://github.com/mageddo/dns-proxy-server/tree/v1-nodejs) 

# Features
Dns-proxy-server is a end user(developers, Server Administrators) DNS server tool to develop systems with docker solving docker containers hostnames:

* Solve names from local configuration database
* Solve names from docker containers using docker **hostname** option or **HOSTNAMES** env
* Solve names from a list of configured DNS servers(as a proxy) if no answer of two above
* [Graphic interface to manage it](http:/127.0.0.1:5380/static/)
	* List and edit DNS local entries

![](http://i.imgur.com/Bhe9P36.png)

# Requirements
* Linux
* Docker 1.8.x
* Docker Compose 1.6.0 (if you use as docker service)

# DNS resolution order
The Dns Proxy Server basically follow the bellow order to solve the names:

* DNS try to solve the hosts from **docker** containers
* then from local database file
* then from 3rd configured remote DNS servers

# Version 2 Improvements
This tool comes from nodejs version(1.0), improving:
* Performance - this version uses much less RAM and is much faster
* Bug fixes
* Binary distribution - now you can simply download a linux executable and use it, without need to install anything
* Code design quality
* And more

# Running it

### From docker

	$ docker run --hostname dns.mageddo --name dns-proxy-server -p 5380:5380 \
    -v /opt/dns-proxy-server/conf:/app/conf \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /etc/resolv.conf:/etc/resolv.conf \
    defreitas/dns-proxy-server

### Standalone run

Download the [latest version](https://github.com/mageddo/dns-proxy-server/releases), extract and run

	$ sudo ./dns-proxy-server

# If you need options 

	$ ./dns-proxy-server --help
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

You can also configure the options at the configuration file

./conf/config.json

```javascript
{
  "remoteDnsServers": [], // not used at the current version
  "envs": [ // all existent environments 
    {
      "name": "", // empty string is the default
      "hostnames": [ // all local hostnames entries
        {
          "id": 1,
          "hostname": "github.com",
          "ip": [192, 168, 0, 1],
          "ttl": 255
        }
      ]
    }
  ],
  "activeEnv": "", // the default env keyname 
  "lastId": 1, // hostnames sequence, don't touch here
  "webServerPort": 0, // web admin port, when 0 the default value is used, see --help option
  "dnsServerPort": 8980 // dns server port, when 0 the default value is used
}
```

# Installing it as a service

1. Download the [latest release](https://github.com/mageddo/dns-proxy-server/releases) and extract it
2. Run the service installer

		$ sudo ./dns-proxy-server -service=docker

this way it will start with the **OS**

if you want to stop 

	$ sudo service dns-proxy-server stop
	Stopping service…
	Service stopped

if you don't want this service anymore

	$ sudo service dns-proxy-server uninstall
	Are you really sure you want to uninstall this service? That cannot be undone. [yes|No] 
	yes

# Testing the DNS server

Testing website

	$ nslookup google.com <dns-server-ip>
	Server:   172.17.0.2
	Address:  172.17.0.2#53

	Non-authoritative answer:
	Name: google.com
	Address: 216.58.202.142

Testing container hostname

	$ nslookup dns.mageddo <dns-server-ip>
	Server:   172.17.0.2
	Address:  172.17.0.2#53
	
	Non-authoritative answer:
	Name: dns.mageddo
	Address: 172.17.0.2

Specifying a port

	nslookup -port=8980 bookmarks-node.mageddo.in 127.0.0.1
