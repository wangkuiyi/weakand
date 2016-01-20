package weakand

import (
	"os"
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	idx := testBuildIndex()
	idx.Pretty(NewPlotTable(os.Stdout), nil, nil, 0)

	query := NewDocument(idx.Vocab.Terms, idx.Vocab) // query includes all terms.
	fr := newFrontier(query, idx)
	idx.Pretty(NewPlotTable(os.Stdout), fr.terms, fr.postings, fr.cur)
}
