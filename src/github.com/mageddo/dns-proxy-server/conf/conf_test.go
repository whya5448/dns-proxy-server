package conf

import (
	"testing"
	"fmt"
)

func TestGetCurrentIpAddress(t *testing.T){

	ip, err := getCurrentIpAddress()
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(">> " + ip)

}