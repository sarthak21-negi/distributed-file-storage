package p2p

import( 
	"io"
	"encoding/GOB"
)
type Decoder interface{
	Decode(io.Reader, any) error
}

type GOBDecoder struct{}

func(dec GOBDecoder) Decode(r io.Reader, v any) error{
	return gob.NewDecoder(r).Decode(v)
}