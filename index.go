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

type SearchIndex struct {
	Fwd   ForwardIndex
	Ivt   InvertedIndex
	Vocab *Vocab
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

func NewIndex(vocab *Vocab) *SearchIndex {
	if vocab == nil {
		vocab = NewVocab(nil)
	}
	return &SearchIndex{
		Ivt:   make(map[TermId]PostingList),
		Fwd:   make(map[DocId]*Document),
		Vocab: vocab}
}

// add Add a document into the index, but not sorting posting lists.
func (idx *SearchIndex) add(doc []string, changed map[TermId]int) {
	did := documentHash(doc)
	if _, ok := idx.Fwd[did]; !ok { // De-duplicate
		d := NewDocument(doc, idx.Vocab)
		idx.Fwd[did] = d
		for term, tf := range d.Terms {
			idx.Ivt[term] = append(idx.Ivt[term], Posting{DocId: did, TF: tf})
			changed[term]++
		}
	}
}

// Add a single document into index.
func (idx *SearchIndex) Add(doc []string) *SearchIndex {
	changed := make(map[TermId]int)
	idx.add(doc, changed)

	for termId := range changed {
		sort.Sort(idx.Ivt[termId])
	}

	return idx
}

// BatchAdd adds all documents read from channel corpus into the index.
func (idx *SearchIndex) BatchAdd(corpus chan []string) *SearchIndex {
	changed := make(map[TermId]int)
	for doc := range corpus {
		idx.add(doc, changed)
	}

	for termId := range changed {
		sort.Sort(idx.Ivt[termId])
	}

	return idx
}

func documentHash(terms []string) DocId {
	var buf bytes.Buffer
	for _, t := range terms {
		fmt.Fprintf(&buf, "%s\t", t)
	}
	md5Bytes := md5.Sum(buf.Bytes())
	return DocId(binary.BigEndian.Uint64(md5Bytes[:]))
}

func (d *Document) Pretty(vocab *Vocab) string {
	var buf bytes.Buffer
	for t, n := range d.Terms {
		if n > 1 {
			fmt.Fprintf(&buf, "%dx", n)
		}
		fmt.Fprintf(&buf, "%s ", vocab.Term(t))
	}
	return buf.String()
}
