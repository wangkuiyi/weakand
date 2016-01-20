package weakand

import "container/heap"

// ResultHeap is a mean-heap with size cap.
type ResultHeap struct {
	rank  []Result
	index map[DocId]int // DocId->index_in_rank, also help dedup.
	cap   int
}

type Result struct {
	p *Posting
	s float64
}

func NewResultHeap(cap int) *ResultHeap {
	return &ResultHeap{
		rank:  make([]Result, 0),
		index: make(map[DocId]int),
		cap:   cap,
	}
}

func (h *ResultHeap) Len() int           { return len(h.rank) }
func (h *ResultHeap) Less(i, j int) bool { return h.rank[i].s < h.rank[j].s } // TODO(y): Make sure it is a MIN-heap.
func (h *ResultHeap) Swap(i, j int)      { h.rank[i], h.rank[j] = h.rank[j], h.rank[i] }

func (h *ResultHeap) Push(x interface{}) {
	docId := x.(Result).p.DocId
	if i, ok := h.index[docId]; ok {
		h.rank[i] = x.(Result)
	} else {
		h.rank = append(h.rank, x.(Result))
		h.index[docId] = h.Len() - 1
	}
}
func (h *ResultHeap) Pop() interface{} {
	l := h.Len()
	r := h.rank[l-1]
	h.rank = h.rank[:l-1]
	delete(h.index, r.p.DocId)
	return r
}

func (h *ResultHeap) Grow(x Result) {
	docId := x.p.DocId
	if i, ok := h.index[docId]; ok {
		h.rank[i] = x
	} else if h.Len() < h.cap {
		h.Push(x)
		heap.Fix(h, h.Len()-1)
	} else if h.rank[0].s < x.s {
		oldDocId := h.rank[0].p.DocId
		h.rank[0] = x
		delete(h.index, oldDocId)
		h.index[docId] = 0
		heap.Fix(h, 0)
	}
}

// Sort() sorts h in descending order of Result.Score, given that h is in heapified status.
func (h *ResultHeap) Sort() []Result {
	r := make([]Result, h.Len())
	i := h.Len() - 1
	for h.Len() > 0 {
		h.Swap(0, h.Len()-1)
		r[i] = h.Pop().(Result)
		i--
		heap.Fix(h, 0)
	}
	return r
}
