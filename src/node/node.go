package node

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/tomaszkoziara/btc-handshake/src/protocol"
	"github.com/tomaszkoziara/btc-handshake/src/protocol/payload"
)

type Node struct {
	nodeVersion  int32
	networkMagic uint32
}

func (n *Node) Connect(address string) error {
	fmt.Println("connecting to the node")
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("error while establishing tcp connection: %v", err)
	}
	defer conn.Close()

	if err := n.handshake(conn); err != nil {
		return err
	}

	return nil
}

func (n *Node) handshake(conn net.Conn) error {
	fmt.Println("sending version")
	err := n.sendVersion(conn)
	if err != nil {
		return err
	}

	fmt.Println("reading version")
	err = n.readVersion(conn)
	if err != nil {
		return err
	}

	fmt.Println("sending verack")
	err = n.sendVerack(conn)
	if err != nil {
		return err
	}

	fmt.Println("reading verack")
	err = n.readVerack(conn)
	if err != nil {
		return err
	}

	// Check if connection is still open
	// TODO: that is not relevant for the node, but it can be part of some e2e test
	if ok, err := isConnectionAlive(conn); err == nil && ok {
		fmt.Println("connection is still alive!")
	}

	return nil
}

func (n *Node) sendVersion(conn net.Conn) error {
	now := time.Now().Unix()
	recvIP := net.ParseIP("192.168.1.1")
	fromIP := net.ParseIP("192.168.1.1")

	// TODO: fill the fields appropriately
	versionPayload := payload.NewVersion(n.nodeVersion).
		WithServices(0).
		WithTimestamp(now).
		WithAddrRecv(*payload.NewNetworkAddr(n.nodeVersion, true).
			WithServices(0).
			WithIP(recvIP).
			WithPort(1)).
		WithAddrFrom(*payload.NewNetworkAddr(n.nodeVersion, true).
			WithServices(0).
			WithIP(fromIP).
			WithPort(1)).
		WithNonce(0).
		WithUserAgent("").
		WithStartHeight(0).
		WithRelay(true)

	versionHeader := new(protocol.Header).
		WithMagic(n.networkMagic).
		WithCommand("version")

	versionMsg := protocol.Message{
		Header:  *versionHeader,
		Payload: versionPayload,
	}

	return versionMsg.Encode(conn)
}

func (n *Node) readVersion(conn net.Conn) error {
	receivedHeader := new(protocol.Header)
	receivedHeader.Decode(conn)
	if receivedHeader.Command != "version" {
		return fmt.Errorf("error, received command %v", receivedHeader.Command)
	}
	if receivedHeader.Magic != n.networkMagic {
		return fmt.Errorf("error, received magic %v", receivedHeader.Magic)
	}
	// TODO: decode version and perform appropriate checks
	if _, err := conn.Read(make([]byte, receivedHeader.Length)); err != nil {
		return fmt.Errorf("error, reading version content %v", err)
	}

	return nil
}

func (n *Node) sendVerack(conn net.Conn) error {
	verackHeader := new(protocol.Header).
		WithMagic(n.networkMagic).
		WithCommand("verack")
	verackMsg := protocol.Message{
		Header: *verackHeader,
	}
	return verackMsg.Encode(conn)
}

func (n *Node) readVerack(conn net.Conn) error {
	receivedHeader := new(protocol.Header)
	receivedHeader.Decode(conn)
	if receivedHeader.Command != "verack" {
		return fmt.Errorf("error, received command %v", receivedHeader.Command)
	}
	if receivedHeader.Magic != n.networkMagic {
		return fmt.Errorf("error, received magic %v", receivedHeader.Magic)
	}

	return nil
}

func isConnectionAlive(conn net.Conn) (bool, error) {
	err := conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		return false, err
	}
	one := []byte{0}
	if _, err := conn.Read(one); err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return true, nil
		}
		return false, nil
	}
	return true, nil
}

func NewNode(nodeVersion int32, networkMagic uint32) *Node {
	return &Node{
		nodeVersion:  nodeVersion,
		networkMagic: networkMagic,
	}
}
