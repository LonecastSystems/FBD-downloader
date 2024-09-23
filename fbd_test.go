package main

import (
	"testing"
)

func TestCountries(t *testing.T) {
	for country, code := range Countries {
		_, exists := Leagues[code]
		if !exists {
			t.Errorf("%v does not exist with a league", country)
		}
	}
}
