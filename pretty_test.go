package weakand

import (
	"os"
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	v, ivt, fwd := testBuildIndex()
	PrettyPrint(os.Stdout, fwd, ivt, v, nil, nil, 0)

	query := NewDocument(v.Terms, v) // query includes all terms.
	fr := newFrontier(query, ivt, fwd)
	PrettyPrint(os.Stdout, fwd, ivt, v, fr.terms, fr.postings, fr.cur)
}
