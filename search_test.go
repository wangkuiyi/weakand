package weakand

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// func TestRetrieve(t *testing.T) {
// 	ch := make(chan []string)
// 	go func() {
// 		for _, d := range testingCorpus {
// 			ch <- d
// 		}
// 		close(ch)
// 	}()

// 	v := NewVocab(nil)
// 	ivtIdx, fwdIdx := BuildIndex(ch, v)

// 	rs := Retrieve(NewDocument([]string{"apple"}, v), 10, ivtIdx, fwdIdx)
// 	spew.Dump(rs)
// }

func TestNewFrontier(t *testing.T) {
	v, ivt, fwd := testBuildIndex()
	q := NewDocument(v.Terms, v)
	spew.Dump(Retrieve(q, 10, ivt, fwd, v))
}
