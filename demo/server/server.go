package main

import "github.com/arthurkiller/razor"
import "fmt"

func main() {
	server,err := razor.Listen("127.0.0.1:1888")
	if err != nil {
		panic(err)
	}

	conn,err := server.Accept()
	if err != nil {
		panic(err)
	}

	buf := make([]byte,10)
	n,err := conn.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println("server read from client",string(buf[:n]),n)
}
