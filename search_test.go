package weakand

import (
	"flag"
	"path"
	"strings"
	"testing"

	"github.com/huichen/sego"
	"github.com/stretchr/testify/assert"
)

var (
	sgmt *sego.Segmenter

	pretty       bool
	indexDumpDir string
)

func init() {
	flag.BoolVar(&pretty, "pretty", false, "Pretty print index and frontier when calling Search")
	flag.StringVar(&indexDumpDir, "indexDir", "/tmp", "Directory containing index dumps")
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
		[]string{"incomplete", "ontologies"},
		"aaai14titlesindex.csv")
}

func TestSearchWithZhWikiNews(t *testing.T) {
	testWithBigData(t,
		"github.com/wangkuiyi/weakand/testdata/zhwikinews.txt",
		[]string{"中", "药商"},
		"zhwikinewsindex.csv")
}

func testWithBigData(t *testing.T, corpusFile string, query []string, indexDumpFile string) {
	guaranteeSegmenter(&sgmt)
	idx := NewIndexFromFile(
		path.Join(gosrc(), corpusFile),
		sgmt,
		path.Join(indexDumpDir, indexDumpFile))

	q := NewQuery(query, idx.Vocab)
	for _, r := range Search(q, 10, idx, pretty) {
		doc := idx.Fwd[r.p.DocId].Pretty(idx.Vocab)

		contain := false
		for _, qterm := range query {
			contain = contain || strings.Contains(doc, qterm)
		}
		assert.True(t, contain)
	}
}

func guaranteeSegmenter(sgmt **sego.Segmenter) error {
	if *sgmt == nil {
		s := new(sego.Segmenter)
		if e := s.LoadDictionary(path.Join(gosrc(),
			"github.com/huichen/sego/data/dictionary.txt")); e != nil {
			return e
		}
		*sgmt = s
	}
	return nil
}
