package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)
type TCPPeer struct{
	net.Conn
	outbound bool
}

func (p *TCPPeer) Send(b []byte) error{
	_,err := p.Conn.Write(b)
	return err
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer{
	return &TCPPeer{
		Conn: conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct{
	ListenAddr string
	HandshakeFunc HandshakeFunc
	Decoder Decoder
	OnPeer func(Peer) error 
}
type TCPTransport struct{
	TCPTransportOpts
	listener net.Listener
	rpcch chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport{
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch: make(chan RPC),
	}
}
//Consume implements transport interface, will return read-only channel
// for reading the incoming message received from incoming peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) Close() error{
	return t.listener.Close()
}

// Dial implements the transport interface
func (t *TCPTransport) Dial(addr string) error{
	conn, err := net.Dial("tcp", addr)
	if err != nil{
		return err
	}

	fmt.Println(conn)

	go t.handleConn(conn,true)

	return nil
}

func(t *TCPTransport) ListenAndAccept() error{
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	log.Printf("TCP transport listening on port: %s\n", t.ListenAddr)
	
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for{
		conn, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed){
			return
		}
		if err != nil{
			fmt.Printf("TCP accept error: %s\n", err)
		}
		fmt.Printf("new incoming connection %+v\n",conn)

		go t.handleConn(conn,true)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool){
	var err error

	defer func() {
		fmt.Printf("Dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

	if err = t.HandshakeFunc(peer); err != nil{
		return 
	}
	
	if t.OnPeer != nil{
		if err = t.OnPeer(peer); err != nil{
			return
		}
	}

	rpc := RPC{}
	for{
	    err = t.Decoder.Decode(conn, &rpc) 
		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}
}