package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)
type PathKey struct{
	PathName string
	Original string
} 

func CASPathTransformFunc(key string) PathKey{
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		Original: hashStr,
	}
}
type PathTransformFunc func(string) PathKey

type StoreOpts struct{
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) string{
	return key
}
type Store struct{
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (p PathKey) Filename() string{
	return fmt.Sprintf ("%s/%s", p.PathName, p.Original)
}

func(s *Store) writeStream(key string, r io.Reader) error{
	pathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	pathAndFileName := pathKey.Filename()

	f,err := os.Create(pathAndFileName)
	if err != nil{
		return err
	}

	n,err := io.Copy(f, r)
	if err != nil{
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, pathAndFileName)

	return nil
}