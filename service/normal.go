package service

import (
	"github.com/mageddo/dns-proxy-server/utils"
)

func NewNormalScript() (*Script) {
	return &Script{utils.GetPath("/dns-proxy-server")}
}