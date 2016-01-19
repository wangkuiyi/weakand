package weakand

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/olekukonko/tablewriter"
)

func PrettyPrint(w io.Writer, fwd ForwardIndex, ivt InvertedIndex, vocab *Vocab, terms []TermId, postings []int, currentDoc DocId) {
	table := tablewriter.NewWriter(w)

	// Convert terms and postings into a map TermdId->DocId
	termDoc := make(map[TermId]DocId)
	for i, t := range terms {
		if postings[i] < len(ivt[t]) {
			termDoc[t] = ivt[t][postings[i]].DocId
		} else {
			termDoc[t] = EndOfPostingList
		}
	}

	// Construct a posting list containing all documents, and use
	// PostList's sortablility to sort them and get docId->index mapping.
	ps := make(PostingList, 0, len(fwd))
	for d, _ := range fwd {
		ps = append(ps, Posting{DocId: d, TF: 0})
	}
	ps = append(ps, Posting{DocId: EndOfPostingList, TF: 0})
	sort.Sort(ps)

	docIdx := make(map[DocId]int)
	for i, p := range ps {
		docIdx[p.DocId] = i
	}
	docIdx[EndOfPostingList] = len(ps)

	row := []string{"Term"}
	for _, p := range ps {
		s := fmt.Sprintf("%016x", p.DocId)
		if currentDoc == p.DocId {
			s = "●" + s
		} else {
			s = " " + s
		}
		row = append(row, s)
	}
	table.SetHeader(row)

	// NOTE: Do not range over ivt, which is a map and range is random.
	for termId, term := range vocab.Terms {
		pl := ivt[TermId(termId)]
		row = make([]string, len(fwd))
		for _, p := range pl {
			mark := "○"
			if p.DocId == termDoc[TermId(termId)] {
				mark = "●"
			}
			row[docIdx[p.DocId]] = mark
		}
		table.Append(append(append([]string{term}, row...), ""))
	}

	row = make([]string, len(fwd))
	for d, c := range fwd {
		var buf bytes.Buffer
		for t, n := range c.Terms {
			if n > 1 {
				fmt.Fprintf(&buf, "%dx", n)
			}
			fmt.Fprintf(&buf, "%s ", vocab.Term(t))
		}
		row[docIdx[d]] = buf.String()
	}
	table.SetFooter(append(append([]string{" "}, row...), " "))

	table.Render() // Send output
}
