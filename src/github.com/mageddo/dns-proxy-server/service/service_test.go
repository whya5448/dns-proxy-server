package service

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"github.com/mageddo/log"
)

func TestSetupFor_NormalModeSuccess(t *testing.T) {
	ctx := log.GetContext()

	const SERVICE_FILE = "/tmp/serviceFile"
	sc := NewService(ctx)
	err := sc.SetupFor(SERVICE_FILE, &Script{"ls"})
	if err != nil {
		t.Error(err)
	}

	bytes, err := ioutil.ReadFile(SERVICE_FILE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	assert.Equal(t, EXPECTED_SERVICE_FILE, string(bytes))

}

const EXPECTED_SERVICE_FILE = `x`
