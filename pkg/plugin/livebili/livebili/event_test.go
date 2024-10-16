package livebili

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDo(t *testing.T) {
	b := &biliPlugin{liveState: make(map[int64]bool)}
	b.liveState[3546597585587017] = false

	assert.NoError(t, b.doCheckLive())

}
