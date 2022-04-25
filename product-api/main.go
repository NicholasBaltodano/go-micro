package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NicholasBaltodano/go-micro/product-api/handlers"
	"github.com/gorilla/mux"
)

//var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind Address for the server")

func main() {
	// Set logger
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	// Define handlers
	productHandler := handlers.NewProduct(logger)

	// ServeMux creation and handler set up
	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", productHandler.GetProducts)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.Use(productHandler.MiddlewareValidateProduct)
	putRouter.HandleFunc("/{id:[0-9]+}", productHandler.UpdateProducts)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.Use(productHandler.MiddlewareValidateProduct)
	postRouter.HandleFunc("/", productHandler.AddProduct)

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
