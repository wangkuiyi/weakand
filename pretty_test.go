package weakand

import (
	"os"
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	idx := testBuildIndex()
	PrettyPrint(NewPlotTable(os.Stdout), idx, nil, nil, 0)

	query := NewDocument(idx.Vocab.Terms, idx.Vocab) // query includes all terms.
	fr := newFrontier(query, idx)
	PrettyPrint(NewPlotTable(os.Stdout), idx, fr.terms, fr.postings, fr.cur)
}
