package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/sarthak21-negi/distributed-file-storage/p2p"
)

type FileServerOpts struct{
	StorageRoot string
	PathTransformFunc PathTransformFunc
	Transport p2p.Transport
	BootstrapNodes []string
}

type FileServer struct{
	FileServerOpts

	peerLock sync.Mutex
	peers map[string]p2p.Peer

	store *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer{
	storeOpts := StoreOpts{
		Root: opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store: NewStore(storeOpts),
		quitch: make(chan struct{}),
		peers: make(map[string]p2p.Peer),
	}
}

type Payload struct{
	Key string
	Data []byte
}

func (s *FileServer) broadcast(p *Payload) error{
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(p)
}

func (s *FileServer) StoreData(key string, r io.Reader) error{

	if err := s.store.Write(key, r); err != nil {
        return err
    }

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r)
	if err != nil {
		return err
	}
	
	p := &Payload{
		Key: key,
		Data: buf.Bytes(),
	}

    fmt.Println(buf.Bytes())

    return s.broadcast(p)
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error{
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p 

	log.Printf("connected with remote %s", p.RemoteAddr())
	
	return nil
}

func (s *FileServer) loop() {
	defer func(){
		log.Println("file server stopped due to user quit action")
		s.Transport.Close()
	}()
	for{
		select{
		case msg := <-s.Transport.Consume():
		fmt.Println(msg)
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) bootstrapNetwork() error{
	for _, addr := range s.BootstrapNodes{
		if len(addr) == 0{
			continue
		}
		go func(addr string) {
			fmt.Println("attempting to connect with remote: ", addr)

			if err := s.Transport.Dial(addr); err != nil{
				log.Println("dial error: ", err)
			}
		}(addr)
	}
	return nil
}
 
func (s *FileServer) Start() error{
	if err := s.Transport.ListenAndAccept(); err != nil{
		return err
	}

	s.bootstrapNetwork()

	s.loop()

	return nil
}

