package weakand

import (
	"bufio"
	"io"
	"log"
	"strings"
)

type Vocab struct {
	TermIndex map[string]int // term to term-Id.
	Terms     []string       // term-Id to term.
}

func NewVocab(in io.Reader) *Vocab {
	v := &Vocab{
		TermIndex: make(map[string]int),
		Terms:     make([]string, 0),
	}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		fs := strings.Fields(scanner.Text())
		// Assumes that each line has multiple fields, and the last one is the term.
		if len(fs) > 0 {
			v.Id(fs[len(fs)-1]) // Add fs[1] to v.
		}
	}
	if e := scanner.Err(); e != nil {
		log.Panicf("Parsing vocab error %v", e)
	}

	return v
}

// Id returns TermId of a term.  If the term is not in v, add it.  Id is not thread-safe.
func (v *Vocab) Id(term string) TermId {
	id, ok := v.TermIndex[term]
	if !ok {
		v.Terms = append(v.Terms, term)
		id = len(v.Terms) - 1
		v.TermIndex[term] = id
		return TermId(id)
	}
	return TermId(id)
}

func (v *Vocab) Term(id TermId) string {
	return v.Terms[int(id)]
}
