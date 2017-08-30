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