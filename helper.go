package main

import (
	"errors"
	"os"
	"strconv"
)

func GetEnvVariables() (email string, password string, path string, yearsInt int, err error) {
	email = os.Getenv("email")
	if email == "" {
		return "", "", "", 0, errors.New("no email")
	}

	password = os.Getenv("password")
	if password == "" {
		return "", "", "", 0, errors.New("no password")
	}

	path = os.Getenv("path")
	if path == "" {
		return "", "", "", 0, errors.New("path is empty")
	}

	years := os.Getenv("years")
	if years == "" {
		return "", "", "", 0, errors.New("no years threshold set")
	}

	yearsInt, err = strconv.Atoi(years)
	if err != nil {
		return "", "", "", 0, errors.New("failed parsing year")
	}

	return email, password, path, yearsInt, err
}
