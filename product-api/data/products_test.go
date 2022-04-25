package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Nicks",
		Price: 5.45,
		SKU:   "abs-asdf-sad"}

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
