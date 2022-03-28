package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	Deleted     string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products {
	return ProductList
}

func UpdateProduct(id int, p *Product) error {
	fp, pos, err := FindProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	fmt.Println(fp)
	ProductList[pos] = p

	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func FindProduct(id int) (*Product, int, error) {
	for i, p := range ProductList {
		if p.ID == id {
			return p, i, nil
		}
	}
	fmt.Println("hello")
	return nil, -1, ErrProductNotFound
}
func AddProduct(p *Product) {
	p.ID = getNextID()
	ProductList = append(ProductList, p)
}

func getNextID() int {
	lp := ProductList[len(ProductList)-1]
	return lp.ID + 1
}

var ProductList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy Milk Cofee",
		Price:       2.45,
		SKU:         "234asd",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Esspresso",
		Description: "Strong coffee no milk",
		Price:       2.85,
		SKU:         "268sasd",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
