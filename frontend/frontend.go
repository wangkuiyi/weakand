package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/rpc"

	"github.com/wangkuiyi/weakand"
)

var (
	bend *rpc.Client

	inputAndSubmit = `<html>
  <body>
    <form action="/text">
      <input type="text" name="query">
      <input type="submit" value="Search">
    </form>
  </body>
</html>
`
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, inputAndSubmit)

	if q := r.FormValue("text"); len(q) > 0 {
		var rs []weakand.Result
		if e := bend.Call("SearchServer.Search", q, &rs); e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
		for _, r := range rs {
			fmt.Fprintf(w, "%s\n", r.Literal) // TODO(y): Print r.Score
		}
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, inputAndSubmit)

	if q := r.FormValue("text"); len(q) > 0 {
		if e := bend.Call("SearchServer.Add", q, nil); e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	backend := flag.String("backend", ":18082", "weak-and backserver listening address")
	addr := flag.String("addr", ":18081", "frontend server listening address")
	flag.Parse()

	var e error
	bend, e = rpc.DialHTTP("tcp", *backend)
	if e != nil {
		log.Fatalf("Cannot dial backend RPC server: %v", e)
	}

	http.HandleFunc("/", searchHandler)
	http.HandleFunc("/add/", addHandler)

	http.ListenAndServe(*addr, nil)
}
