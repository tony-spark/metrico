package main

import (
	"io"
	"log"
	"net/http"
)

func loggingHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(r.Method, r.RequestURI, err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println(r.Method, r.RequestURI, string(b))
}

func main() {
	http.HandleFunc("/", loggingHandler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
