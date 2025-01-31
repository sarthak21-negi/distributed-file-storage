package main

import (
	//"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type PathKey struct{
	PathName string
	FileName string
} 

const defaultRootFolderName = "ggnetwork"

func CASPathTransformFunc(key string) PathKey{
	hash := sha1.Sum([]byte(key))
    hashStr := hex.EncodeToString(hash[:])

    // Replace colons or other invalid characters
    hashStr = strings.ReplaceAll(hashStr, ":", "_")

    blocksize := 5
    sliceLen := len(hashStr) / blocksize
    paths := make([]string, sliceLen)

    for i := 0; i < sliceLen; i++ {
        from, to := i*blocksize, (i*blocksize)+blocksize
        paths[i] = hashStr[from:to]
    }

    return PathKey{
        PathName: strings.Join(paths, "/"),
        FileName: hashStr,
    }

}

type PathTransformFunc func(string) PathKey

type StoreOpts struct{
	Root string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey{
	return PathKey{
		PathName: key,
		FileName: key,
	}
}
type Store struct{
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil{
		opts.PathTransformFunc = DefaultPathTransformFunc
	
	}
	if len(opts.Root) == 0{
		opts.Root = defaultRootFolderName
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (p PathKey) FirstPathName() string{
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0{
		return ""
	}
	return paths[0]
} 

func (p PathKey) FullPath() string{
	return fmt.Sprintf ("%s/%s", p.PathName, p.FileName)
}

func (s *Store) Has(key string) bool{
	pathKey := s.PathTransformFunc(key)

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)

	return !errors.Is(err, os.ErrNotExist)

}

func (s *Store) Clear() error{
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(key string) error{
	pathKey := s.PathTransformFunc(key)
	
	defer func(){
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	firstPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

func (s *Store) Read(key string) (int64, io.Reader, error){
	return s.readStream(key)
}

func (s *Store) readStream(key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
    
	file, err:= os.Open(fullPathWithRoot)
	if err != nil{
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil{
		return 0, nil, err
	}

	return fi.Size(), file, nil
}

func(s *Store) writeStream(key string, r io.Reader) (int64, error) {
	pathKey := s.PathTransformFunc(key)
    pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, strings.ReplaceAll(pathKey.PathName, ":", "_"))
    
    fmt.Println("Creating directory:", pathNameWithRoot) // Debugging line
    
    if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
        fmt.Println("Error creating directories:", err)
        return 0, err
    } else {
        fmt.Println("Directories created successfully")
    }

    fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
    fmt.Println("Creating file at:", fullPathWithRoot) // Debugging line
    f, err := os.Create(fullPathWithRoot)
    if err != nil {
        return 0, err
    }
    defer f.Close()

    n, err := io.Copy(f, r)
    if err != nil {
        return 0, err
    }

    return n, nil
}