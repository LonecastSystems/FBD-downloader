package main

import (
	"log"
)

func main() {
	email, password, path, yearsInt, err := GetEnvVariables()
	if err != nil {
		log.Panic(err.Error())
	}

	Download(email, password, path, yearsInt)
}
