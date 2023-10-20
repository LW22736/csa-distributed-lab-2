package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	//	"net/rpc"
	//	"fmt"
	//	"time"
	//	"net"
)

var nextAddr string
var shutDownSent = false

var SingVerseHandler = "BuddyOperations.SingVerse"
var ShutDownHandler = "BuddyOperations.ShutDown"

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

func (b *BuddyOperations) ShutDown(req Request, baseRes *Response) (err error) {
	if shutDownSent == false {
		res := new(Response)
		res.Message = true
		client, err := rpc.Dial("tcp", nextAddr)
		if err != nil {
			panic(err)
		}
		client.Go(ShutDownHandler, Request{0}, res, nil)
	}
	os.Exit(0)
	return
}

func doVerse(n int) {
	if n > 0 {
		fmt.Println(n, "bottles of beer on the wall,", n, "bottles of beer. Take one down, pass it around...")
		res := new(Response)
		client, err := rpc.Dial("tcp", nextAddr)
		if err != nil {
			panic(err)
		}
		client.Go(SingVerseHandler, Request{n - 1}, res, nil)
	} else {
		fmt.Println("The crippling alcoholism has hit...")
		res := new(Response)
		client, err := rpc.Dial("tcp", nextAddr)
		if err != nil {
			panic(err)
		}
		shutDownSent = true
		client.Go(ShutDownHandler, Request{0}, res, nil)
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
