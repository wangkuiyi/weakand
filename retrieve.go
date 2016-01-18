package weakand

import "sort"

type Frontier struct {
	terms      []TermId
	postings   []int // indexing posting list of InvertedIndex[terms[i]].
	currentDoc DocId
	ivtIdx     InvertedIndex
	fwdIdx     ForwardIndex
}

func newFrontier(query *Document, ivtIdx InvertedIndex, fwdIdx ForwardIndex) *Frontier {
	f := &Frontier{
		terms:      make([]TermId, 0, len(query.Terms)),
		postings:   make([]int, 0, len(query.Terms)),
		currentDoc: 0,
		ivtIdx:     ivtIdx,
		fwdIdx:     fwdIdx}

	for term, _ := range query.Terms {
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

func (f *Frontier) findPivotTerm(threshold float64) int {
	// TODO(y): Implement this.
	return 0
}

// pickTerm returns a value in range [0, pivotTermIdx), or -1 for error.
func (f *Frontier) pickTerm(pivotTermIdx int) int {
	// TODO(y): Implement this.
	return -1
}

func (f *Frontier) score(query *Document, post Posting) float64 {
	return jaccardCoefficient(query, f.fwdIdx[post.DocId])
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

func Retrieve(query *Document, cap int, ivtIdx InvertedIndex, fwdIdx ForwardIndex) []Result {
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
		results.Grow(Result{
			Posting: post,
			Score:   f.score(query, post)}, cap)
	}

	results.Sort()
	return results
}
