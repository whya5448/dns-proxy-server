### 2.2.2
* Cache Rest API v1 is exposed

### 2.2.1
* Preventing nil pointer when container inspection fails

### 2.2.0
* Increased code coverage
* Implementing cache at local hostnames and remote server resolution
* Considering TTL to invalidate hostname cache for local resolution

### 2.1.7
* All build and release process is made inside docker (no travis dependency)

### 2.1.6
* Refactor project structure to save dependencies in vendor folder

### 2.1.5
* Automating build with Travis

### 2.1.1
* Fix - `Error response from daemon: No such container...` message. see #29  
* Fix - hostname don't get removed when the container has killed. see #26  

### 2.1.0
* Turn publish port optional when running as service using docker mode

### 2.0.21
* BugFix - Service stopped of work in normal mode

### 2.0.20 
* Support for --version option that shows the current version
* Docker Compose is not required anymore to run DNS Proxy Server as a docker service

### 2.0.19
* Ability to customize remote server
* Fixing DNS solution order from (local, docker, remote) to (docker, local, remote)
* Now, at least docker 1.9 API v1.21 is necessary

### 2.0.18
* Making it compatible with docker 1.8 api v1.20
