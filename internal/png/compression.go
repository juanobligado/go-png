package png

import (
	"bytes"
	"compress/zlib"
	"io"
	"io/ioutil"
)



type IDatMemoryReader struct{
	buffer []byte
	
}

func (idat *IDatMemoryReader) Reader() io.Reader  {
	return bytes.NewReader(idat.buffer)
}


func Decompress( chunk_map ChunkMap) ([]byte,error) {

	buffer := make([]byte, 0)

	// Join IDATS and Decode everything into the buffer
	for _,chunk := range chunk_map[IDAT]{

		zreader,err := zlib.NewReader(chunk.Reader())
		if err != nil{
			return nil,err
		}
		chunk_bytes,err := ioutil.ReadAll(zreader)	
		if err != nil{
			return nil,err
		}
		buffer = append(buffer, chunk_bytes...)
	}
	return buffer,nil
}


