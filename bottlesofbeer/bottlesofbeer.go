package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	//	"net/rpc"
	//	"fmt"
	//	"time"
	//	"net"
)

var nextAddr string

var SingVerseHandler = "BuddyOperations.SingVerse"

type Response struct {
	Message bool
}

type Request struct {
	Message int
}

type BuddyOperations struct{}

func (b *BuddyOperations) SingVerse(req Request, res *Response) (err error) {
	doVerse(req.Message)

	res.Message = true
	return
}

func doVerse(n int) {
	if n > 0 {
		fmt.Println(n, "bottles of beer on the wall,", n, "bottles of beer. Take one down, pass it around...")
		res := new(Response)
		client, _ := rpc.Dial("tcp", nextAddr)
		client.Go(SingVerseHandler, Request{n - 1}, res, nil)
	} else {
		fmt.Println("The crippling alcoholism has hit...")
	}
}

func main() {
	thisPort := flag.String("this", "8030", "Port for this process to listen on")
	flag.StringVar(&nextAddr, "next", "localhost:8040", "IP:Port string for next member of the round.")
	bottles := flag.Int("n", 0, "Bottles of Beer (launches song if not 0)")
	flag.Parse()

	err := rpc.Register(&BuddyOperations{})
	if err != nil {
		panic(err)
	}

	listener, _ := net.Listen("tcp", ":"+*thisPort)
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			panic(err)
		}
	}(listener)

	if *bottles != 0 {
		go doVerse(*bottles)
	}
	fmt.Println("Starting to run rpc")
	rpc.Accept(listener)
}
