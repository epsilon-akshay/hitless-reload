package main

import "net"

func main() {
	l, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	l.Write([]byte("learnnn"))
	l.Close()
}
