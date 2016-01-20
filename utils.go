package weakand

import (
	"log"
	"os"
	"unicode"

	"github.com/huichen/sego"
)

func OpenOrDie(file string) *os.File {
	f, e := os.Open(file)
	if e != nil {
		log.Panic(e)
	}
	return f
}

func WithFile(file string, fn func(f *os.File)) {
	f := OpenOrDie(file)
	defer f.Close()
	fn(f)
}

func CreateOrDie(file string) *os.File {
	f, e := os.Create(file)
	if e != nil {
		log.Panic(e)
	}
	return f
}

func AllPunctOrSpace(s string) bool {
	for _, u := range s {
		if !unicode.IsPunct(u) && !unicode.IsSpace(u) {
			return false
		}
	}
	return true
}

func Tokenize(doc string, sgmt *sego.Segmenter) []string {
	var terms []string
	for _, seg := range sgmt.Segment([]byte(doc)) {
		term := seg.Token().Text()
		if !AllPunctOrSpace(term) {
			terms = append(terms, term)
		}
	}
	return terms
}
