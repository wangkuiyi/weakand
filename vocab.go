package weakand

type Vocab struct {
	TermIndex map[string]int
	Terms     []string
}

// Id is not thread-safe.
func (v Vocab) Id(term string) TermId {
	id, ok := v.TermIndex[term]
	if !ok {
		v.Terms = append(v.Terms, term)
		id = len(v.Terms) - 1
		v.TermIndex[term] = id
		return TermId(id)
	}
	return TermId(id)
}

func (v Vocab) Term(id TermId) string {
	return v.Terms[int(id)]
}
