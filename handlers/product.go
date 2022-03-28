package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/NicholasBaltodano/go-micro/product-api/data"
)

type Products struct {
	logger *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.logger.Println("getProducts called")
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.logger.Println("PUT", r.URL.Path)
		// expect ID in URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		results := reg.FindAllSubmatch([]byte(r.URL.Path), -1)

		// Check for multiple IDs
		if len(results) != 1 {
			p.logger.Println("Put, URI multiple ID's")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(results[0]) != 2 {
			p.logger.Println("Put, invalid URI, multiple capture groups", string(results[0][0]))
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := results[0][1]
		id, err := strconv.Atoi(string(idString))
		if err != nil {
			p.logger.Println("Failed to convert id , ", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
		}

		p.logger.Println("got id", id)
		p.updateProducts(id, rw, r)
	}

	// Catch All
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		p.logger.Println(err)
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Post product called")
	newProduct := &data.Product{}

	err := newProduct.FromJSON(r.Body)
	if err != nil {
		p.logger.Println(err)
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
	}

	p.logger.Printf("Prod: %#v", newProduct)
	data.AddProduct(newProduct)
}

func (p Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("PUT product called")
	newProduct := &data.Product{}

	err := newProduct.FromJSON(r.Body)
	if err != nil {
		p.logger.Println(err)
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, newProduct)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Failed to update product", http.StatusNotFound)
	}
}
