package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NicholasBaltodano/go-micro/handlers"
)

func main() {
	// Set logger
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	// Define handlers
	//helloHandler := handlers.NewHello(logger)
	goodbyeHandler := handlers.NewGoodbye(logger)
	productHandler := handlers.NewProduct(logger)

	// ServeMux creation and handler set up
	serveMux := http.NewServeMux()
	//serveMux.Handle("/", helloHandler)
	//serveMux.Handle("/goodbye", goodbyeHandler)
	serveMux.Handle("/goodbye", goodbyeHandler)
	//serveMux.Handle("/products", productHandler)
	serveMux.Handle("/", productHandler)
	serveMux.HandleFunc("/favicon", doNothing)

	// Create server with custom servemux
	server := &http.Server{
		Addr:         "localhost:9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	// Non-block listen and serve
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	logger.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}

func doNothing(w http.ResponseWriter, r *http.Request) {
}
