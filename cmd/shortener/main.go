package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var urls []string
var mutex sync.Mutex

func main() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r.URL.Path)
	case http.MethodPost:
		handlePost(w, r.URL.Path, r.Body)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handleGet(w http.ResponseWriter, path string) {
	if len(path) <= 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(path[1:])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id >= len(urls) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", urls[id])
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func handlePost(w http.ResponseWriter, path string, body io.ReadCloser) {
	if path != "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := readBody(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mutex.Lock()
	urls = append(urls, url)
	id := len(urls) - 1
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "http://localhost:8080/%d", id)
}

func readBody(body io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(body); err != nil {
		return "", err
	}

	return buf.String(), nil
}
