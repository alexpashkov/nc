package main

import (
	"flag"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	var (
		addr        string
		startServer bool
		timeout     time.Duration
	)
	flag.StringVar(&addr, "addr", "localhost:2222", "addr to communicate with")
	flag.BoolVar(&startServer, "serve", false, "set to true if you want to start a server, takes precedence over write")
	flag.DurationVar(&timeout, "timeout", time.Second*10, "client dial timeout")
	flag.Parse()
	if startServer {
		log.Printf("starting server on %s", addr)
		if err := serve(addr); err != nil {
			log.Fatal(errors.Wrap(err, "serving error"))
		}
	} else {
		log.Printf("writing to %s", addr)
		if err := write(addr, timeout); err != nil {
			log.Fatal(errors.Wrap(err, "writing error"))
		}
	}
}

func serve(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			return errors.Wrap(err, "failed to accept")
		}
		if err := handleConn(conn); err != nil {
			return errors.Wrap(err, "connection handling error")
		}
	}
}

func handleConn(conn net.Conn) error {
	defer conn.Close()
	_, err := io.Copy(os.Stdout, conn)
	return err
}

func write(address string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return errors.Wrap(err, "dial error")
	}
	defer conn.Close()
	_, err = io.Copy(conn, os.Stdin)
	return errors.Wrap(err, "failed to copy input into the connectino")
}
