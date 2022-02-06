package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/favicon", doNothing)
	http.HandleFunc("/", hello)
	http.HandleFunc("/goodbye", goodbye)
	http.ListenAndServe("127.0.0.1:9090", nil)
}

func doNothing(w http.ResponseWriter, r *http.Request) {
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello World !")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "oops", http.StatusBadRequest)
	}

	fmt.Fprintf(w, "Hello %s", data)
}

func goodbye(w http.ResponseWriter, r *http.Request) {
	log.Println("Goodbye world!")
}
