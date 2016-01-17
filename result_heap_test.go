package weakand

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func shuffledFloat64Slice(size int) []float64 {
	r := make([]float64, size)
	for i := range r {
		r[i] = float64(i)
	}

	for i := range r {
		j := rand.Intn(i + 1)
		r[i], r[j] = r[j], r[i]
	}

	return r
}

func TestResultHeap(t *testing.T) {
	assert := assert.New(t)

	var mh ResultHeap
	size := 10

	sh := shuffledFloat64Slice(1024 * 1024)
	for _, s := range sh {
		mh.Grow(Result{Score: s}, size)
	}
	assert.Equal(size, len(mh))

	mh.Sort()
	for i := 0; i < len(mh); i++ {
		assert.Equal(float64(1024*1024-i-1), mh[i].Score)
	}
}
