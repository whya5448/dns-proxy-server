package docker

import (
	"fmt"
	"github.com/docker/engine-api/types/filters"
	"encoding/json"
	"github.com/docker/engine-api/types/events"
	"io"
	"github.com/docker/engine-api/client"
	"golang.org/x/net/context"
	"github.com/mageddo/log"
	"github.com/docker/engine-api/types"
	"strings"
)

var cache = make(map[string]string)

func HandleDockerEvents(){

	logCtx := log.GetContext()
	logger := log.GetLogger(logCtx)

	// adaptar a api do docker aqui
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, nil)
	//cli, err := client.NewEnvClient()
	if err != nil {
		logger.Errorf("m=solve, status=error-to-connect-at-host, solver=docker, err=%v", err)
		return
	}

	// more about list containers https://docs.docker.com/engine/reference/commandline/ps/
	options := types.ContainerListOptions{}
	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, options)
	if err != nil {
		logger.Errorf("m=solve, status=error-to-list-container, solver=docker, err=%v", err)
		return
	}

	for _, c := range containers {
		cInspection, err := cli.ContainerInspect(ctx, c.ID)
		logger.Infof("status=container-from-list-begin, container=%s", cInspection.Name)
		if err != nil {
			logger.Errorf("status=inspect-error-at-list, container=%s, err=%v", c.Names, err)
		}
		hostnames := getHostnames(cInspection)
		putHostnames(logCtx, hostnames, cInspection)
		logger.Infof("status=container-from-list-success, container=%s, hostnames=%s", cInspection.Name, hostnames)
	}

	// more about events here: http://docs-stage.docker.com/v1.10/engine/reference/commandline/events/
	var eventFilter = filters.NewArgs()
	eventFilter.Add("event", "start")

	eventFilter.Add("event", "die")
	eventFilter.Add("event", "stop")

	body, err := cli.Events(ctx, types.EventsOptions{ Filters: eventFilter })
	if err != nil {
		logger.Errorf("m=solve, status=error-to-attach-at-events-handler, solver=docker, err=%v", err)
		return
	}

	dec := json.NewDecoder(body)
	for {

		var event events.Message
		err := dec.Decode(&event)
		if err != nil && err == io.EOF {
			break
		}

		cInspection, err := cli.ContainerInspect(ctx, event.ID)
		if err != nil {
			logger.Errorf("status=inspect-error, container=%s, err=%v", cInspection.Name, err)
		}
		hostnames := getHostnames(cInspection)
		logger.Infof("status=resolved-hosts, action=%s, hostnames=%s", event.Action, hostnames)

		switch event.Action {
		case "start":
			putHostnames(ctx, hostnames, cInspection)
			break

		case "die":
		case "stop":
			for _, host := range hostnames {
				remove(host)
			}
			break

		}
	}

}

func ContainsKey(key string) bool {
	_, ok := cache[key]
	if ok {
		return true
	}
	return false
}

func Get(key string) string {
	return cache[key]
}

func GetCache() map[string]string {
	return cache
}

func remove(key string){
	if ContainsKey(key) {
		delete(cache, key)
	}
}

func getHostnames(inspect types.ContainerJSON) []string {

	const hostnameEnv = "HOSTNAMES="
	hostnames := *new([]string)

	if len(inspect.Config.Hostname) != 0 {
		if len(inspect.Config.Domainname) != 0 {
			hostnames = append(hostnames, fmt.Sprintf("%s.%s", inspect.Config.Hostname, inspect.Config.Domainname))
		}else {
			hostnames = append(hostnames, inspect.Config.Hostname)
		}
	}
	for _, env := range inspect.Config.Env {

		envName := strings.Index(env, hostnameEnv)
		if envName == 0 {
			envValue := env[envName + len(hostnameEnv) : ]
			hostnames = append(hostnames, strings.Split(envValue, ",")...)
			return hostnames
		}
	}
	return hostnames
}

func putHostnames(ctx context.Context, hostnames []string, inspect types.ContainerJSON){
	logger := log.GetLogger(ctx)
	for _, host := range hostnames {

		var ip string
		for _, network := range inspect.NetworkSettings.Networks {
			ip = network.IPAddress
		}
		if len(ip) == 0 {
			panic(fmt.Sprintf("no network found to %s", inspect.Name))
		}
		logger.Debugf("m=putHostnames, host=%s, ip=%s", host, ip)
		cache[host] = ip
	}
}