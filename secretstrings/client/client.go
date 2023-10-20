package main

import (
	"bufio"
	"flag"
	"net/rpc"
	"os"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
	//	"bufio"
	//	"os"
	"fmt"
)

func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)

	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	file, _ := os.Open("../wordlist")
	defer file.Close()

	reader := bufio.NewReader(file)
	word, err := reader.ReadString('\n')
	fmt.Println(word, err)
	for err == nil {
		request := stubs.Request{Message: word}
		response := new(stubs.Response)
		client.Call(stubs.PremiumReverseHandler, request, response)
		fmt.Println("Responded:", response.Message)
		word, err = reader.ReadString('\n')
	}
}
