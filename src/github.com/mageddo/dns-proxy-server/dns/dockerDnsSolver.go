package dns

import (
	"github.com/miekg/dns"
	"fmt"
	"encoding/json"
	"io"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"golang.org/x/net/context"
	"github.com/docker/engine-api/types/filters"
	"github.com/mageddo/log"
)

type DockerDnsSolver struct {

}

func (*DockerDnsSolver) Solve(name string) *dns.Msg {
	// adaptar a api do docker aqui
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, nil)
	//cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// more about list containers https://docs.docker.com/engine/reference/commandline/ps/
	options := types.ContainerListOptions{}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		fmt.Printf("%s - %s\n", c.Names[0], c.ID)
	}

	// more about events here: http://docs-stage.docker.com/v1.10/engine/reference/commandline/events/
	var eventFilter = filters.NewArgs()

	eventFilter.Add("event", "start")

	eventFilter.Add("event", "die")
	eventFilter.Add("event", "stop")

	body, err := cli.Events(context.Background(), types.EventsOptions{ Filters: eventFilter })
	if err != nil {
		log.Logger.Fatal(err)
	}

	dec := json.NewDecoder(body)
	for {
		var event events.Message
		err := dec.Decode(&event)
		if err != nil && err == io.EOF {
			break
		}

		log.Logger.Info(event)
	}
	return nil
}
