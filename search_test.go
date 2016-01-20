package weakand

import (
	"bufio"
	"os"
	"path"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/wangkuiyi/sego"
)

func TestSearch(t *testing.T) {
	idx := testBuildIndex()
	q := NewDocument(idx.Vocab.Terms, idx.Vocab) // query includes all terms.
	rs := Search(q, 10, idx, true)               // Pretty print intermediate steps.
	assert.Equal(t, len(testingCorpus), len(rs)) // All documents should be retrieved.
	for _, r := range rs {
		assert.Equal(t, 0.5, r.s) // Jaccard coeffcient of all documents should be 1/2.
	}
}

func TestSearchAAAI14Data(t *testing.T) {
	ch := make(chan []string)
	go func() {
		withFile(path.Join(gosrc(), "github.com/wangkuiyi/weakand/testdata/aaai14papers.txt"),
			func(f *os.File) {
				var sgmt sego.Segmenter
				assert.Nil(t, sgmt.LoadDictionary(path.Join(gosrc(), "github.com/huichen/sego/data/dictionary.txt")))

				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					var terms []string
					for _, seg := range sgmt.Segment([]byte(scanner.Text())) {
						if term := seg.Token().Text(); !allPunctOrSpace(term) {
							terms = append(terms, term)
						}
					}
					ch <- terms
				}
				if e := scanner.Err(); e != nil {
					t.Skipf("Scanning corpus error:", e)
				}
			})
		close(ch)
	}()

	idx := NewIndex(NewVocab(nil)).BatchAdd(ch)

	// Note: to print the forward&inverted index into a CSV that
	// can be loaded into Apple Numbers or Microsoft Excel, just
	// uncomment the following line:
	//
	// PrettyPrint(NewCSVTable(os.Stdout), fwd, ivt, vocab, nil, nil, 0)

	q := NewDocument([]string{"incomplete", "ontologies"}, idx.Vocab)
	for _, r := range Search(q, 10, idx, false) { // No pretty print intermediate steps.
		doc := idx.Fwd[r.p.DocId].Pretty(idx.Vocab)
		assert.True(t, strings.Contains(doc, "incomplete") || strings.Contains(doc, "ontologies"))
	}
}

func allPunctOrSpace(s string) bool {
	for _, u := range s {
		if !unicode.IsPunct(u) && !unicode.IsSpace(u) {
			return false
		}
	}
	return true
}
