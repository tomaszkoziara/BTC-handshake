package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tomaszkoziara/btc-handshake/src/node"
)

func main() {
	nodeVersion := int32(70015)
	regtestNetworkMagic := uint32(0xDAB5BFFA)

	node := node.NewNode(nodeVersion, regtestNetworkMagic)
	// TODO: insert node discovery and connect to nodes in a non-blocking way
	go func() {
		if err := node.Connect("localhost:18444"); err != nil {
			log.Fatalf("error connecting to node: %v", err)
		}
	}()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	_ = <-sigChannel
	fmt.Println("Bye bye.")
}
