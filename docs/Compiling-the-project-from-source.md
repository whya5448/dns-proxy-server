# Building from source with docker

Docker-Dns-proxy uses docker to simplify the compile process


Generate the binaries

	$ docker-compose up prod-build-dns-proxy-server
	Starting docker-dns-server-compiler
	Attaching to docker-dns-server-compiler
	docker-dns-server-compiler           | ok  	github.com/mageddo/dns-proxy-server/conf	0.008s
	docker-dns-server-compiler           | ?   	github.com/mageddo/dns-proxy-server/controller	[no test files]
	...
	docker-dns-server-compiler           | github.com/mageddo/dns-proxy-server/flags
	...
	docker-dns-server-compiler           | _/app/src
	docker-dns-server-compiler exited with code 0

Create the docker image

    $ docker-compose build prod-build-docker-dns-proy-server

# Used tecnologies 

* Docker
* Docker Compose
* Git
* Golang 1.7


