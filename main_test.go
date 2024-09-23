package main

import (
	"os"
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

func TestGetEnvVariables(t *testing.T) {
	os.Setenv("password", "test")
	os.Setenv("years", "abc")

	_, _, _, _, err := GetEnvVariables()
	if err == nil {
		t.Error("Email is supposed to be invalid")
	}

	os.Clearenv()

	os.Setenv("email", "test@gmail.com")
	os.Setenv("years", "abc")

	_, _, _, _, err = GetEnvVariables()
	if err == nil {
		t.Error("Password is supposed to be invalid")
	}

	os.Clearenv()

	os.Setenv("email", "test@gmail.com")
	os.Setenv("password", "test")

	_, _, _, _, err = GetEnvVariables()
	if err == nil {
		t.Error("Year is supposed to be invalid")
	}

	os.Clearenv()

	os.Setenv("email", "test@gmail.com")
	os.Setenv("password", "test")
	os.Setenv("years", "1")

	_, _, _, _, err = GetEnvVariables()
	if err == nil {
		t.Error("Path is supposed to be invalid")
	}

	os.Clearenv()

	os.Setenv("email", "test@gmail.com")
	os.Setenv("password", "test")
	os.Setenv("path", "test")
	os.Setenv("years", "abc")

	_, _, _, _, err = GetEnvVariables()
	if err == nil {
		t.Error("Year is supposed to be invalid")
	}

	os.Clearenv()

	os.Setenv("email", "test@gmail.com")
	os.Setenv("password", "test")
	os.Setenv("path", "test")
	os.Setenv("years", "1")

	_, _, _, _, err = GetEnvVariables()
	if err != nil {
		t.Error("All environment variables are supposed to be valid")
	}
}
