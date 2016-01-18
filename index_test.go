package weakand

import (
	"bufio"
	"log"
	"os"
	"path"
	"strings"
	"testing"
)

func TestBuildIndex(t *testing.T) {
	corpus := [][]string{
		{"apple", "pie"},
		{"apple", "iphone"},
		{"iphone", "jailbreak"}}

	ch := make(chan []string)

	go func() {
		for _, d := range corpus {
			ch <- d
		}
		close(ch)
	}()

	// TODO(y): Finish this.
}

func TestDocumentHashCollision(t *testing.T) {
	withFile(path.Join(gosrc(), "github.com/wangkuiyi/weakand/testdata/internet-zh.num"),
		func(f *os.File) {
			dict := make(map[DocId][][]string)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				fs := strings.Fields(scanner.Text())
				if len(fs) == 2 {
					content := fs[1:]
					did := documentHash(content)
					if _, ok := dict[did]; ok {
						t.Errorf("Collision between %v and %v", content, dict[did])
					}
					dict[did] = append(dict[did], content)
				}
			}
			if e := scanner.Err(); e != nil {
				t.Errorf("Reading %s error: %v", f.Name(), e)
			}
		})
}

func gosrc() string {
	return path.Join(os.Getenv("GOPATH"), "src")
}

func openOrDie(file string) *os.File {
	f, e := os.Open(file)
	if e != nil {
		log.Panic(e)
	}
	return f
}

func withFile(file string, fn func(f *os.File)) {
	f := openOrDie(file)
	defer f.Close()
	fn(f)
}
