package weakand

import (
	"fmt"
	"testing"
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
	v, ivtIdx, fwdIdx := testBuildIndex()

	var docToId []DocId
	for i := 0; i < len(testingCorpus); i++ {
		docToId = append(docToId, documentHash(testingCorpus[i]))
	}

	query := NewDocument(v.Terms, v)// query includes all terms.
	fr := newFrontier(query, ivtIdx, fwdIdx) 

	for i := 0; i < len(v.Terms); i++ {
	fmt.Println(fr.docId(0))
	fmt.Println(fr.docId(1))
}
