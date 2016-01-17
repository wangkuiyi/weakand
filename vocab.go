package weakand

type Vocab struct {
	TermIndex map[string]int
	Terms     []string
}

// Id is not thread-safe.
func (v Vocab) Id(term string) int {
	id, ok := v.TermIndex[term]
	if !ok {
		v.Terms = append(v.Terms, term)
		id = len(v.Terms) - 1
		v.TermIndex[term] = id
		return id
	}
	return id
}

func (v Vocab) Term(id int) string {
	return v.Terms[id]
}
