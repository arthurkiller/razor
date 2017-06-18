package main

import (
	"fmt"

	"github.com/arthurkiller/razor"
)

func main() {
	client, err := razor.Dial("127.0.0.1:1888")
	if err != nil {
		panic(err)
	}

	n, err := client.Write([]byte("foo"))
	if err != nil {
		fmt.Println(n)
		panic(err)
	}

	fmt.Println("client say foo done!", n)
}
