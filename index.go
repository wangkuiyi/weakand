package weakand

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"sort"
)

type DocId uint64 // MD5 hash of document content.
type TermId int   // depends on Vocab.

type InvertedIndex map[TermId]PostingList
type PostingList []Posting
type Posting struct {
	DocId
	TF int // The term frequency in Doc.
}

type ForwardIndex map[DocId]*Document
type Document struct {
	Terms map[TermId]int // map makes it fast to compute Σt∈q∩d U_t.
	Len   int            // sum over Terms.
}

// If word exists in content but not in vocab, add it into vocab.
func NewDocument(content []string, vocab *Vocab) *Document {
	d := &Document{Terms: make(map[TermId]int)}

	for _, term := range content {
		d.Terms[vocab.Id(term)]++
		d.Len++
	}
	return d
}

// In InvertedIndex, posting lists are sorted by asceding order DocId.
func (p PostingList) Len() int           { return len(p) }
func (p PostingList) Less(i, j int) bool { return p[i].DocId < p[j].DocId }
func (p PostingList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func BuildIndex(corpus chan []string, vocab *Vocab) (InvertedIndex, ForwardIndex) {
	ivtIdx := make(map[TermId]PostingList)
	fwdIdx := make(map[DocId]*Document)

	for doc := range corpus {
		did := documentHash(doc)
		d := NewDocument(doc, vocab)
		fwdIdx[did] = d
		for term, tf := range d.Terms {
			ivtIdx[term] = append(ivtIdx[term], Posting{DocId: did, TF: tf})
		}
	}

	for _, ps := range ivtIdx {
		sort.Sort(ps)
	}

	return ivtIdx, fwdIdx
}

func documentHash(terms []string) DocId {
	var buf bytes.Buffer
	for _, t := range terms {
		fmt.Fprintf(&buf, "%s\t", t)
	}
	md5Bytes := md5.Sum(buf.Bytes())
	return DocId(binary.BigEndian.Uint64(md5Bytes[:]))
}
