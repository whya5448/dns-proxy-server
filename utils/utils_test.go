package utils

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestDiffMillis(t *testing.T) {

	before, err := time.Parse("2006-01-02 15:04:05.999", "2017-05-25 23:58:04.555")
	assert.Nil(t, err)

	after, err := time.Parse("2006-01-02 15:04:05.999", "2017-05-26 23:59:04.555")
	assert.Nil(t, err)

	assert.Equal(t, int64(86460000), DiffMillis(before, after))
}
