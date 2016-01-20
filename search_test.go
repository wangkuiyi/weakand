package weakand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrieve(t *testing.T) {
	v, ivt, fwd := testBuildIndex()
	q := NewDocument(v.Terms, v)                 // query includes all terms.
	rs := Retrieve(q, 10, ivt, fwd, v)           // NOTE: change v to nil to disable debug output.
	assert.Equal(t, len(testingCorpus), len(rs)) // All documents should be retrieved.
	for _, r := range rs {
		assert.Equal(t, 0.5, r.s) // Jaccard coeffcient of all documents should be 1/2.
	}
}
