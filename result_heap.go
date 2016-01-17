package weakand

import "container/heap"

// ResultHeap is a mean-heap with size cap.
type ResultHeap []Result
type Result struct {
	Posting
	Score float64
}

func (mh *ResultHeap) Len() int           { return len(*mh) }
func (mh *ResultHeap) Less(i, j int) bool { return (*mh)[i].Score < (*mh)[j].Score } // TODO(y): Make sure it is a MIN-heap.
func (mh *ResultHeap) Swap(i, j int)      { (*mh)[i], (*mh)[j] = (*mh)[j], (*mh)[i] }
func (mh *ResultHeap) Push(x interface{}) { (*mh) = append(*mh, x.(Result)) }
func (mh *ResultHeap) Pop() interface{}   { l := len(*mh); r := (*mh)[:l-1]; *mh = (*mh)[:l-1]; return r }

func (mh *ResultHeap) Grow(x Result, cap int) {
	if len(*mh) < cap {
		mh.Push(x)
		heap.Fix(mh, len(*mh)-1)
	} else if (*mh)[0].Score < x.Score {
		(*mh)[0] = x
		heap.Fix(mh, 0)
	}
}

// Sort() sorts mh in descending order of Result.Score, given that mh is in heapified status.
func (mh *ResultHeap) Sort() []Result {
	l := len(*mh)
	for len(*mh) > 0 {
		mh.Swap(0, len(*mh)-1)
		mh.Pop()
		heap.Fix(mh, 0)
	}
	(*mh) = (*mh)[0:l]
	return []Result(*mh)
}
