package flags

import (
	"flag"
	"os"
)

var (
	Cpuprofile      = flag.String("cpuprofile", "", "write cpu profile to file")
	Compress        = flag.Bool("compress", false, "compress replies")
	Tsig            = flag.String("tsig", "", "use MD5 hmac tsig: keyname:base64")
	WebServerPort   = flag.Int("web-server-port", 8080, "The web server port")
	DnsServerPort   = flag.Int("server-port", 53, "The DNS server to start into")
	SetupResolvconf = flag.Bool("default-dns", true, "This DNS server will be the default server for this machine")
	ConfPath        = flag.String("conf-path", "conf/config.json", "The config file path ")
	Help            = flag.Bool("help", false, "This message")
)

func init(){

	flag.Parse()
	if *Help {
		flag.PrintDefaults()
		os.Exit(1)
	}

}