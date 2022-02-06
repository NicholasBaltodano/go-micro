package main

import (
	"log"
	"net/http"
	"os"

	"github.com/NicholasBaltodano/go-micro/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)
	helloHandler := handlers.NewHello(logger)
	goodbyeHandler := handlers.NewGoodbye(logger)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", helloHandler)
	serveMux.Handle("/goodbye", goodbyeHandler)

	serveMux.HandleFunc("/favicon", doNothing)

	http.ListenAndServe("127.0.0.1:9090", serveMux)
}

func doNothing(w http.ResponseWriter, r *http.Request) {
}
