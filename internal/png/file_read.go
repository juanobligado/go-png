package png

import (
	"bytes"
	"encoding/binary"
	"os"
)







func check_signature(file *os.File) bool{
    PNG_SIGNATURE := []byte{137,80,78,71,13,10,26,10}
    signature := make([]byte,len(PNG_SIGNATURE))

    _,err := file.Read(signature)
    if err != nil{
        return false;
    }
    return bytes.Compare(PNG_SIGNATURE,signature) == 0 ;
  
}




var buffer_1 []byte  = make([]byte, 1)
var buffer_2 []byte  = make([]byte, 2)
var buffer_4 []byte  = make([]byte, 4)


func read_uint32(file *os.File) (uint32, error) {

    _,err := file.Read(buffer_4)
    if err != nil{
        return 0,err
    }
    return to_uint32(buffer_4),nil
}

func to_uint32(slice []byte) uint32{
    return binary.BigEndian.Uint32(slice)
}

func to_uint16(slice []byte) uint16{
    return binary.BigEndian.Uint16(slice)
}

func read_uint16(file *os.File) (uint16, error) {

    _,err := file.Read(buffer_2)
    if err != nil{
        return 0,err
    }
    return to_uint16(buffer_2),nil
}


func read_uint8(file *os.File) (byte, error) {

    _,err := file.Read(buffer_1)
    if err != nil{
        return 0,err
    }
    return buffer_1[0],nil
}




