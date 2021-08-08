package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	var host = flag.String("host", "0.0.0.0", "host to bind to")
	var port = flag.Int("port", 3000, "port to bind to")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("could not determine directory, wtf is wrong with your computer: %v", err)
		}

		// Handle which file to serve based on the method and presence of a variant
		var file string
		if variant := r.URL.Query().Get("variant"); variant != "" {
			file = fmt.Sprintf("%s-%s.json", r.Method, variant)
		} else {
			file = fmt.Sprintf("%s.json", r.Method)
		}

		p := path.Join(wd, r.URL.Path[1:], file)
		f, err := os.Open(p)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("JSON mock not found at %s", p)))
			return
		}
		defer f.Close()

		w.Header().Add("Content-Type", "application/json")
		io.Copy(w, f)
	})

	log.Printf("server listening at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
