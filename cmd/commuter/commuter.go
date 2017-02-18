package main

import (
	"fmt"
	"os"

	"github.com/KyleBanks/commuter/pkg/geo"
)

func main() {
	key := os.Getenv("GMAPS")
	r, err := geo.NewRouter(key)
	if err != nil {
		panic(err)
	}

	d, err := r.Duration(geo.Route{
		From: "123 Main St.",
		To:   "321 Maple Ave.",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(d)
}
