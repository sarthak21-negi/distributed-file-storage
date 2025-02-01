package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/sarthak21-negi/distributed-file-storage/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer{
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr: listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		EncKey: newEncryptionKey(),
		StorageRoot: strings.ReplaceAll(listenAddr, ":", "_") + "_network",  // Replace ':' with '_'
        PathTransformFunc: CASPathTransformFunc,
        Transport: tcpTransport,
        BootstrapNodes: nodes,
	}

	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main(){
	s1 := makeServer(":3000","")
	s2 := makeServer(":4000", ":3000")

	go func(){
		log.Fatal(s1.Start())
	}()

	time.Sleep(2 * time.Second)

	go s2.Start()
	time.Sleep(2 * time.Second)

	key := "coolPicture.jpg"
	data := bytes.NewReader([]byte("my big data file here!"))
	s2.StoreData("coolPicture.jpg", data)
	
	if err := s2.store.Delete(key); err != nil{
		log.Fatal(err)
	}

	r, err := s2.Get(key)
	if err != nil{
		log.Fatal(err)
	}

	b, err := io.ReadAll(r)
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println(string(b))
}