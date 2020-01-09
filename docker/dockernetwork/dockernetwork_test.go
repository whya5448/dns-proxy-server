package dockernetwork

import (
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/network"
	"github.com/mageddo/dns-proxy-server/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetGatewayIp_MustFindNetworkFromDps(t *testing.T) {
	// arrange
	ctx := reference.Context()
	mockApiClient := &MockApiClient{}
	SetCli(mockApiClient)

	mockApiClient.On("NetworkList", ctx, mock.Anything).
		Return([]types.NetworkResource{{
			Driver: "bridge",
			IPAM: network.IPAM{
				Config: []network.IPAMConfig{{
					Gateway: "192.168.0.1",
				}},
			},
		}}, nil)

	// act
	ip, err := GetGatewayIp(ctx)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, ip, "192.168.0.1")
}
