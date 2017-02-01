package utils

type LocalConfiguration struct {
	RemoteDnsServers [][4]byte
	Envs []EnvVo
}

type EnvVo struct {
	Name string
	Hostnames []HostnameVo
}

type HostnameVo struct {
	ip [4]byte
	ttl int
}
