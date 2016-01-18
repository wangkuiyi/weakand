package weakand

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testVocabContent = `
 apple
 pie
 iphone
 jailbreak
`
)

func TestVocab(t *testing.T) {
	v := NewVocab(strings.NewReader(testVocabContent))

	assert := assert.New(t)
	assert.Equal(4, len(v.Terms))
	assert.Equal(4, len(v.TermIndex))
	for i, t := range v.Terms {
		assert.Equal(i, v.TermIndex[t])
	}
}
