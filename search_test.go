package weakand

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/wangkuiyi/sego"
)

var (
	sgmt   sego.Segmenter
	pretty bool
)

func init() {
	if e := sgmt.LoadDictionary(path.Join(gosrc(),
		"github.com/huichen/sego/data/dictionary.txt")); e != nil {
		log.Panic(e)
	}

	flag.BoolVar(&pretty, "pretty", false, "Pretty print index and frontier when calling Search")
}

func TestSearch(t *testing.T) {
	idx := testBuildIndex()
	q := NewQuery(idx.Vocab.Terms, idx.Vocab)    // query includes all terms.
	rs := Search(q, 10, idx, pretty)             // Pretty print intermediate steps.
	assert.Equal(t, len(testingCorpus), len(rs)) // All documents should be retrieved.
	for _, r := range rs {
		assert.Equal(t, 0.5, r.s) // Jaccard coeffcient of all documents should be 1/2.
	}
}

func TestSearchWithAAAI14Titles(t *testing.T) {
	testWithBigData(t,
		"github.com/wangkuiyi/weakand/testdata/aaai14papers.txt",
		[]string{"incomplete", "ontologies"})
}

func TestSearchWithZhWikiNews(t *testing.T) {
	testWithBigData(t,
		"github.com/wangkuiyi/weakand/testdata/zhwikinews.txt",
		[]string{"中药", "商行"})
}

func testWithBigData(t *testing.T, corpusFile string, query []string) {
	ch := make(chan []string)
	go func() {
		withFile(path.Join(gosrc(), corpusFile),
			func(f *os.File) {
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

	q := NewQuery(query, idx.Vocab)
	for _, r := range Search(q, 10, idx, pretty) { // No pretty print intermediate steps.
		doc := idx.Fwd[r.p.DocId].Pretty(idx.Vocab)
		fmt.Println(doc) //debug
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
