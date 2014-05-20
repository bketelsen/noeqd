package main

import (
	"fmt"

	"github.com/bketelsen/noeqd"
)

func main() {
	var generator *noeqd.Generator
	id, _ := generator.Get()
	fmt.Println("Generated: ", id)
}
