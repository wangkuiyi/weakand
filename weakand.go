package weakand

import "sort"

type DocId uint64 // MD5 hash of document content.
type TermId int   // depends on Vocab.

type InvertedIndex map[TermId]PostingList
type PostingList []Posting
type Posting struct {
	DocId
	TF int // The term frequency in Doc.
}

type ForwardIndex map[DocId]Document
type Document struct {
	Terms map[TermId]int // map makes it fast to compute Σt∈q∩d U_t.
}

// In InvertedIndex, posting lists are sorted by asceding order DocId.
func (p PostingList) Len() int           { return len(p) }
func (p PostingList) Less(i, j int) bool { return p[i].DocId < p[j].DocId }
func (p PostingList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func BuildIndex(corpus chan []string, vocab *Vocab) (InvertedIndex, ForwardIndex) {
	ivtIdx := make(map[TermId]PostingList)
	fwdIdx := make(map[DocId]Document)

	for doc := range corpus {
		d := Document{Terms: make(map[TermId]int)}
		for _, term := range doc {
			d.Terms[vocab.Id(term)]++
		}

		did := DocumentId(d)
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

type Frontier struct {
	terms      []TermId
	postings   []int // indexing posting list of InvertedIndex[terms[i]].
	currentDoc DocId
	ivtIdx     InvertedIndex
	fwdIdx     ForwardIndex
}

func newFrontier(query Document, ivtIdx InvertedIndex, fwdIdx ForwardIndex) *Frontier {
	f := &Frontier{
		terms:      make([]TermId, 0, len(query.Terms)),
		postings:   make([]int, 0, len(query.Terms)),
		currentDoc: 0,
		ivtIdx:     ivtIdx,
		fwdIdx:     fwdIdx}

	for term, tf := range query.Terms {
		if _, ok := ivtIdx[term]; ok {
			// NOTE: Initialziing Frontier.postings to 0 implies postings lists has minimal length 1.
			f.postings = append(f.postings, 0)
			f.terms = append(f.terms, term)
		}
	}
	sort.Sort(f)
	return f
}

// sort.Sort(f) sorts f.terms and f.postrings.
func (f *Frontier) Len() int { return len(f.postings) }
func (f *Frontier) Less(i, j int) bool {
	di := f.ivtIdx[f.terms[i]][f.postings[i]].DocId
	dj := f.ivtIdx[f.terms[j]][f.postings[j]].DocId
	return di < dj
}
func (f *Frontier) Swap(i, j int) {
	f.terms[i], f.terms[j] = f.terms[j], f.terms[i]
	f.postings[i], f.postings[j] = f.postings[j], f.postings[i]
}

func scan(f *Frontier, threshold func() float64, emit chan Posting) {
	for {
		pivotTermIdx := f.findPivotTerm(threshold())
		if pivotTermIdx < 0 {
			return // No more docs
		}

		pivotDocIdx := f.postings[pivotTermIdx]
		if pivotDocIdx >= len(f.ivtIdx[f.terms[pivotTermIdx]]) {
			return // No more docs
		}

		pivot := f.ivtIdx[f.terms[pivotTermIdx]][pivotDocIdx].DocId
		if pivot < f.currentDoc {
			// pivot has been considerred, advance one of the preceeding terms.
			f.postings[f.pickTerm(pivotTermIdx)]++
		} else {
			if p := f.ivtIdx[f.terms[0]][f.postings[0]]; p.DocId == pivot {
				// Success, all terms preceeding pTerm belong to the pivot.
				f.currentDoc = pivot
				emit <- p
			} else {
				// Not enough mass yet on pivot, advance one of the preceeding terms.
				f.postings[f.pickTerm(pivotTermIdx)]++
			}
		}

		sort.Sort(f)
	}
}

type ResultHeap []Result
type Result struct {
	Posting
	Score float64
}

func Retrieve(query Document, cap int, ivtIdx InvertedIndex, fwdIdx ForwardIndex) []Result {
	var results ResultHeap
	threshold := func() float64 {
		if len(results) <= 0 {
			return 0.0
		}
		return results[0].Score // TODO(y): Introduce factor F.
	}

	f := newFrontier(query, ivtIdx, fwdIdx)
	candidates := make(chan Posting)
	go func() {
		scan(f, threshold, candidates)
		close(candidates)
	}()

	for post := range candidates {
		results.Grow(Result{Posting: post, Score: Score(query, post)}, cap)
	}

	sort.Sort(&results)
	return results
}

func (r ResultHeap) Len() int            { return len(r) }
func (r ResultHeap) Less(i, j int) bool  { return r[i].Score < r[j].Score } // TODO(y): Make sure it is a MIN-heap.
func (r ResultHeap) Swap(i, j int)       { r[i], r[j] = r[j], r[i] }
func (r *ResultHeap) Push(x interface{}) { *r = append(*r, x.(Result)) }
func (r *ResultHeap) Pop() interface{}   { l := (*r)[len(*r)-1]; *r = (*r)[:len(*r)-1]; return l }

func (r *ResultHeap) Grow(d Result, cap int) {
	if len(*r) < cap {
		(*r) = append(ResultHeap{d}, (*r)...) // TODO(y): Unit test it.
	}
}
