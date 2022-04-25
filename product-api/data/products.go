package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	Deleted     string  `json:"-"`
}

// Struct functions
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// SKU is of format abc-adss-sdesf
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

// Types and Vars
type Products []*Product

var ErrProductNotFound = fmt.Errorf("Product not found")

// Package Functions
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

// Data
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
