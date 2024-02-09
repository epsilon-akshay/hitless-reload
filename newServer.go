package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

func main() {
	commonC := make(chan net.Conn, 1)
	go func() {
		const addr = "akshay.sock"
		//syscall.Unlink(addr)

		l, err := net.Listen("unix", addr)
		if err != nil {
			log.Fatal(err)
		}
		for {
			c, err := l.Accept()
			if err != nil {
				panic(err)
			}

			buf := make([]byte, 1)
			oob := make([]byte, 1024)
			_, oobn, _, _, err := c.(*net.UnixConn).ReadMsgUnix(buf, oob)
			if err != nil {
				panic(err)
			}
			scms, err := syscall.ParseSocketControlMessage(oob[:oobn])
			if err != nil {
				panic(err)
			}

			fd, err := syscall.ParseUnixRights(&scms[0])
			if err != nil {
				panic(err)
			}
			cc, err := net.FileConn(os.NewFile(uintptr(fd[0]), ""))
			if err != nil {
				panic(err)
			}

			commonC <- cc
		}
		defer l.Close()
	}()

	go func() {
		l, err := net.Listen("tcp", ":8082")
		if err != nil {
			panic(err)
		}

		go func() {
			a := <-commonC
			b := make([]byte, 1092)
			a.Read(b)
			fmt.Println("handling b", string(b))
		}()

		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}
			b := make([]byte, 1024)
			_, err = conn.Read(b)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		}
	}()
	time.Sleep(1 * time.Hour)
}
