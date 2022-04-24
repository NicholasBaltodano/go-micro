package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/NicholasBaltodano/go-micro/product-api/data"
	"github.com/gorilla/mux"
)

type Products struct {
	logger *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		p.logger.Println(err)
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Post product called")

	// Validation done by middleware
	newProduct := r.Context().Value(KeyProduct{}).(data.Product)

	p.logger.Printf("Prod: %#v", newProduct)
	data.AddProduct(&newProduct)
}

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	// Gorilla ID goes into mux.Vars(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}

	p.logger.Println("PUT product called - ID = ", id)
	// Validation done by middleware
	newProduct := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &newProduct)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Failed to update product", http.StatusNotFound)
	}
}

type KeyProduct struct{}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.logger.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
