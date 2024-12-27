package main

import (
	"log"
	"time"

	"github.com/sarthak21-negi/distributed-file-storage/p2p"
)

func main(){
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot: "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport: tcpTransport,
	}

	s := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()
	
	if err := s.Start(); err != nil{
		log.Fatal(err)
	}

}