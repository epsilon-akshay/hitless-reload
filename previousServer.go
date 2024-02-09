package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var CommonConn uintptr
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-sigc
		const addr = "akshay.sock"
		//syscall.Unlink(addr)

		//TRY changing  laddress
		l, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: addr, Net: "unix"})
		if err != nil {
			log.Fatal(err)
		}

		buf := []byte{1}
		_, _, err = l.WriteMsgUnix(buf, syscall.UnixRights(int(CommonConn)), nil)
		if err != nil {
			panic(err)
		}

		fmt.Println("writing all existing connections to neeraj.sock")

		l.Close()

		panic("closing server ")

	}()

	go func() {
		l, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}

			tcpC, ok := conn.(*net.TCPConn)
			if !ok {
				panic("error no  conn  is  tcp")
			}

			f, err := tcpC.File()
			if err != nil {
				panic(err)
			}

			CommonConn = f.Fd()

			time.Sleep(10 * time.Minute)

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
