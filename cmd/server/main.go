package main

import (
	api "github.com/AndrewDonelson/go-zeroknowledge-proof/pkg/api"
)

func main() {
	api.RunAPI(":8080", "../../certs/")
}
