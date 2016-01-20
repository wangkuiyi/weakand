package weakand

import (
	"os"
	"sort"
)

// A special value indicating the end of posting list.
const EndOfPostingList = DocId(^uint64(0))

type Frontier struct {
	*SearchIndex
	terms    []TermId
	postings []int // indexing posting list of InvertedIndex[terms[i]].
	cur      DocId
}

func newFrontier(query *Document, idx *SearchIndex) *Frontier {
	f := &Frontier{
		SearchIndex: idx,
		terms:       make([]TermId, 0, len(query.Terms)),
		postings:    make([]int, 0, len(query.Terms)),
		cur:         0}

	for term, _ := range query.Terms {
		if _, ok := idx.Ivt[term]; ok {
			// NOTE: Initialziing Frontier.postings to 0 implies postings lists has minimal length 1.
			f.postings = append(f.postings, 0)
			f.terms = append(f.terms, term)
		}
	}
	sort.Sort(f)
	return f
}

// sort.Sort(f) sorts f.terms and f.postings.
func (f *Frontier) Len() int { return len(f.postings) }
func (f *Frontier) Less(i, j int) bool {
	return f.docId(i) < f.docId(j)
}
func (f *Frontier) Swap(i, j int) {
	f.terms[i], f.terms[j] = f.terms[j], f.terms[i]
	f.postings[i], f.postings[j] = f.postings[j], f.postings[i]
}

func (f *Frontier) docId(frontierIdx int) DocId {
	term := f.terms[frontierIdx]
	post := f.postings[frontierIdx]
	plist := f.Ivt[term]
	if post >= len(plist) {
		return EndOfPostingList
	}
	return plist[post].DocId
}

func scan(f *Frontier, threshold func() float64, emit chan *Posting, vocab *Vocab, debug bool) {
	for {
		if debug {
			PrettyPrint(NewPlotTable(os.Stdout), f.SearchIndex, f.terms, f.postings, f.cur)
		}

		pivotTermIdx := f.findPivotTerm(threshold())
		if pivotTermIdx < 0 {
			return // No more docs
		}

		pivotDocIdx := f.postings[pivotTermIdx]
		if pivotDocIdx >= len(f.Ivt[f.terms[pivotTermIdx]]) {
			return // No more docs
		}

		pivot := f.Ivt[f.terms[pivotTermIdx]][pivotDocIdx].DocId
		if pivot < f.cur {
			// pivot has been considerred, advance one of the preceeding terms.
			f.postings[f.pickTerm(pivotTermIdx)]++
		} else {
			if p := &f.Ivt[f.terms[0]][f.postings[0]]; p.DocId == pivot {
				// Success, all terms preceeding pTerm belong to the pivot.
				f.cur = pivot
				emit <- p
				f.postings[0]++
			} else {
				// Not enough mass yet on pivot, advance one of the preceeding terms.
				f.postings[f.pickTerm(pivotTermIdx)]++
			}
		}

		sort.Sort(f)
	}
}

func (f *Frontier) findPivotTerm(threshold float64) int {
	// TODO(y): Implement this.
	return 0
}

// pickTerm returns a value in range [0, pivotTermIdx), or -1 for error.
func (f *Frontier) pickTerm(pivotTermIdx int) int {
	// TODO(y): Implement this.
	return -1
}

func (f *Frontier) score(query *Document, post *Posting) float64 {
	return jaccardCoefficient(query, f.Fwd[post.DocId])
}

func jaccardCoefficient(q, d *Document) float64 {
	inters := 0
	for t, f := range q.Terms {
		inters += min(d.Terms[t], f)
	}
	return float64(inters) / float64(q.Len+d.Len-inters)
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func Search(query *Document, cap int, idx *SearchIndex, debug bool) []Result {
	results := NewResultHeap(cap)
	threshold := func() float64 {
		if results.Len() < cap {
			return 0.0
		}
		return results.rank[0].s // TODO(y): Introduce factor F.
	}

	f := newFrontier(query, idx)
	candidates := make(chan *Posting)
	go func() {
		scan(f, threshold, candidates, idx.Vocab, debug)
		close(candidates)
	}()

	for post := range candidates {
		results.Grow(Result{
			p: post,
			s: f.score(query, post)})
	}

	return results.Sort()
}
