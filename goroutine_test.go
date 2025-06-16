package async_test

import (
	"testing"

	"github.com/reugn/async"

	"github.com/reugn/async/internal/assert"
)

func TestGoroutineID(t *testing.T) {
	gid, err := async.GoroutineID()

	assert.IsNil(t, err)
	t.Log(gid)
}

func BenchmarkGetGroutineID3(b *testing.B) {
	for range b.N {
		_, err := async.GoroutineID()
		if err != nil {
			b.Error("failed to get gid")
		}
	}
}
