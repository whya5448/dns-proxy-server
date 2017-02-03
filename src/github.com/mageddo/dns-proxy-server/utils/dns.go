package utils

type LocalConfiguration struct {
	RemoteDnsServers [][4]byte
	Envs []EnvVo
	ActiveEnv string
}

type EnvVo struct {
	Name string
	Hostnames []HostnameVo
}

type HostnameVo struct {
	Hostname string
	Ip [4]byte
	Ttl int
}
func (lc *LocalConfiguration) GetActiveEnv() *EnvVo {

	for _, env := range lc.Envs {
		if env.Name == lc.ActiveEnv {
			return &env
		}
	}
	return nil
}

func(env *EnvVo) GetHostname(hostname string) *HostnameVo {
	for _, host := range env.Hostnames {
		if host.Hostname == hostname {
			return &host
		}
	}
	return nil
}