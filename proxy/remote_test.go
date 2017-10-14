package proxy

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/miekg/dns"
	"github.com/mageddo/go-logging"
	"github.com/mageddo/dns-proxy-server/events/local"
	"context"
	"github.com/mageddo/dns-proxy-server/log"
)

type MockedRemoteDnsSolver struct {
	mock.Mock
	RemoteDnsSolver
}

func (m MockedRemoteDnsSolver) loadConfiguration(ctx context.Context) (*local.LocalConfiguration, error) {
	log.LOGGER.Infof("status=mocked-config")
	return local.LoadConfiguration(ctx)
}

func TestRemoteDnsSolver_SolveCacheSuccess(t *testing.T) {

	ctx := logging.NewContext()

	remoteSolver := new(MockedRemoteDnsSolver)
	remoteSolver.confloader = func(ctx context.Context) (*local.LocalConfiguration, error) {
		remoteSolver.MethodCalled("confloader", ctx)
		return local.LoadConfiguration(ctx)
	}

	remoteSolver.On("confloader", ctx).Once()

	question := new(dns.Question)
	remoteSolver.Solve(ctx, *question)
	remoteSolver.Solve(ctx, *question)

	remoteSolver.AssertExpectations(t)

}

