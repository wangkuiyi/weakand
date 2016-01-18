package weakand

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVocab(t *testing.T) {
	vocab := `
 apple
 pie
 iphone
 jailbreak
`
	v := NewVocab(strings.NewReader(vocab))

	assert := assert.New(t)
	assert.Equal(4, len(v.Terms))
	assert.Equal(4, len(v.TermIndex))
	for i, t := range v.Terms {
		assert.Equal(i, v.TermIndex[t])
	}
}
