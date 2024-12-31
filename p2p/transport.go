package p2p

import "net"

type Peer interface {
	Send([]byte) error
	RemoteAddr() net.Addr
	Close() error
}

type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}