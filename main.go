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
	var dirArg = flag.String("dir", "", "directory to serve")
	var host = flag.String("host", "0.0.0.0", "host to bind to")
	var port = flag.Int("port", 3000, "port to bind to")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not determine directory, wtf is wrong with your computer: %v", err)
	}

	var dir string
	if *dirArg != "" {
		dir = path.Join(wd, *dirArg)
	} else {
		dir = wd
	}

	// Support passing the serving directory as an argument, not a flag
	// We'll only listen to the first directory passed
	if len(flag.Args()) > 0 {
		dir = path.Join(wd, flag.Args()[0])
	}

	log.Printf("serving from %s\n", dir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle which file to serve based on the method and presence of a variant
		var file string
		if variant := r.URL.Query().Get("variant"); variant != "" {
			file = fmt.Sprintf("%s-%s.json", r.Method, variant)
		} else {
			file = fmt.Sprintf("%s.json", r.Method)
		}

		p := path.Join(dir, r.URL.Path[1:], file)
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
