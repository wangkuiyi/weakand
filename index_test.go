package weakand

import (
	"bufio"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testingCorpus = [][]string{
		{"apple", "pie"},
		{"apple", "iphone"},
		{"iphone", "jailbreak"}}
)

func TestBuildIndex(t *testing.T) {
	ch := make(chan []string)
	go func() {
		for _, d := range testingCorpus {
			ch <- d
		}
		close(ch)
	}()

	v := NewVocab(nil)
	ivtIdx, fwdIdx := BuildIndex(ch, v)

	assert := assert.New(t)
	assert.Equal(4, len(v.Terms))
	assert.Equal(4, len(v.TermIndex))
	assert.Equal(len(testingCorpus), len(fwdIdx))
	assert.Equal(4, len(ivtIdx))
	for i := range ivtIdx {
		assert.True(sort.IsSorted(ivtIdx[i]))
	}
}

func TestDocumentHashCollision(t *testing.T) {
	withFile(path.Join(gosrc(), "github.com/wangkuiyi/weakand/testdata/internet-zh.num"),
		func(f *os.File) {
			dict := make(map[DocId][][]string)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				fs := strings.Fields(scanner.Text())
				if len(fs) == 2 {
					content := fs[1:]
					did := documentHash(content)
					if _, ok := dict[did]; ok {
						t.Errorf("Collision between %v and %v", content, dict[did])
					}
					dict[did] = append(dict[did], content)
				}
			}
			if e := scanner.Err(); e != nil {
				t.Errorf("Reading %s error: %v", f.Name(), e)
			}
		})
}

func gosrc() string {
	return path.Join(os.Getenv("GOPATH"), "src")
}

func openOrDie(file string) *os.File {
	f, e := os.Open(file)
	if e != nil {
		log.Panic(e)
	}
	return f
}

func withFile(file string, fn func(f *os.File)) {
	f := openOrDie(file)
	defer f.Close()
	fn(f)
}
