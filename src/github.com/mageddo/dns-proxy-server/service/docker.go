package service

import (
	"strings"
	"fmt"
	"github.com/mageddo/dns-proxy-server/conf"
	"github.com/mageddo/dns-proxy-server/flags"
)

func NewDockerScript() *Script {

	script := `'PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin ; ` +
		`docker rm -f dns-proxy-server &> /dev/null ;` +
		`docker run --hostname dns.mageddo --name dns-proxy-server -p %d:%d ` +
		`-v /opt/dns-proxy-server/conf:/app/conf ` +
		`-v /var/run/docker.sock:/var/run/docker.sock ` +
		`-v /etc/resolv.conf:/etc/resolv.conf ` +
		`defreitas/dns-proxy-server:%s'`
	script = strings.Replace(script, "/", "\\/", -1)
	script = strings.Replace(script, "&", "\\&", -1)
	script = fmt.Sprintf(script, conf.WebServerPort(), conf.WebServerPort(), flags.GetRawCurrentVersion())
	return &Script{script}
}