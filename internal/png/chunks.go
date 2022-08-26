package png

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
)

const (
	// Critical Chunks
	IHDR string = "IHDR"
	PLTE string = "PLTE"
	IDAT string = "IDAT"
	IEND string = "IEND"

	// Histogram chunk if and only if having PLTE
	hIST string ="hIST"
	bKGD string ="bKGD"
	gAMA string ="gAMA"
	sRGB string ="sRGB"
	pHYs string ="pHYs"
	Byte     uint32 = 1
	Kibibyte      = 1024 * Byte
	Mebibyte      = 1024 * Kibibyte
	Gibibyte      = 1024 * Mebibyte
	MaxChunkSize =  100*Mebibyte

	Greyscale = byte(0)
	TrueColor = byte(2)
	IndexedColor = byte(3)
	GreyscaleWithAlpha = byte(4)
	TrueColorWithAlpha = byte(6)

)



type Chunk struct {
	offset int64
	size   uint32
	code   []byte
	data   IByteBuffer
	crc32  uint32	
}
func (self *Chunk) Reader() io.Reader{
	return  bytes.NewReader(self.data.GetRawData())
}

func (self *Chunk) Iscritical() bool{
	str := string(self.code)	
	Criticals   := []string{IHDR,PLTE,IDAT,IEND}
	for _,v := range Criticals{
		if v == str {
			return true
		}
	}
	return false
	
}



type IHDRChunk struct {
	Width       uint32
	Height      uint32
	Bitdepth    byte
	Colortype   byte
	Compression byte
	Filter      byte
	Interlace   byte
}


func (self *IHDRChunk) HasAlpha() bool{
	if self.Colortype == TrueColorWithAlpha || self.Colortype == GreyscaleWithAlpha{
		return true
	}
	return false
}


func (self *IHDRChunk) GetSamples() (uint32,error){
	
	switch self.Colortype{
		case Greyscale:
		case IndexedColor:
			{
				return uint32(1),nil
			}
		case GreyscaleWithAlpha:
			{
				return uint32(2),nil
			}
		case TrueColor:
			{
				return uint32(3),nil
			}
		case TrueColorWithAlpha:
			{
				return uint32(4),nil
			}
	}
	return uint32(0),errors.New("Color Type not supported")
	
}


func (self *IHDRChunk) GetBytesPerLine() (uint32,error){
	data_points_per_pixel,err := self.GetSamples()
	if err != nil{
		return 0,err
	}
	bits_per_pixel := uint64(self.Bitdepth*byte(data_points_per_pixel))
	line_bits := bits_per_pixel * uint64(self.Width)
	line_bytes := line_bits / 8
	if line_bits % 8 != 0{
		line_bytes++
	}
	return uint32(line_bytes),nil
}


type ChunkMap map[string][]Chunk

func (self ChunkMap) IHDR() (*IHDRChunk, error){
	chunks := self[IHDR]
	if chunks == nil || len(chunks) == 0{
		return nil,errors.New("Can find IHDR Chunk in Map")
	}
	chunk := chunks[0]
	ihdr := IHDRChunk{}
	chunkValue := reflect.ValueOf(&ihdr)
	err := read_struct(chunk.data.GetRawData()  ,chunkValue)
	return &ihdr , err
}

func (self ChunkMap) PLTE() (*Palette, error){
	return nil,nil
	// TODO: Read Palette
}



func read_data_buffer(file *os.File,  chunksize uint32) (IByteBuffer,error){
	if chunksize < MaxChunkSize {
		buffer := MemoryBuffer{data: make([]byte, chunksize)}
		_,err := file.Read(buffer.data)
		return &buffer,err
		
	}else{
		// create lazy load for chunks
		return nil,errors.New("Big Data Not Supported")
	}
}

// Reads Struct using reflection
func read_struct(data []byte,chunkValue reflect.Value) error{
	index := uint32(0)
	for i:=0 ; i < chunkValue.Elem().NumField();i++ {
		field := chunkValue.Elem().Field(i)
		increment := uint32( field.Type().Size() )

		switch field.Kind(){

			case reflect.Uint32:{
				a :=  to_uint32(data[index:index+increment])
				field.SetUint(uint64(a))
				break;			
			}
			case reflect.Uint16:{
				increment = 2
				a := to_uint16(data[index:index+increment])
				field.SetUint(uint64(a))
				break;
			}
			case reflect.Uint8:{
				increment = 1
				field.SetUint(uint64(data[index]))
				break;
			}
		}
		index = index + increment
	} 
		

	return nil
}

func  read_chunk( file *os.File  ) (*Chunk,error) {

	chunk := Chunk{}
	chunksize,err := read_uint32(file)
    if err!= nil{
        return nil,err
    }
	chunk.size = chunksize
	chunk.offset,_ =  file.Seek(0,1)

	// read chunk-code
	chunk.code = make([]byte,4)
    _,err = file.Read(chunk.code)
    if err != nil {
        return nil,err
    }
	// Check Chunk Header
	for _,a := range chunk.code{
		if (a >= byte('a') && a <= byte('z')) || (a >= byte('A') && a <= byte('Z')){
			continue
		   }else{
			return &chunk,errors.New(fmt.Sprintf("Invalid Character %x found while processing chunk header '%x'",a,chunk.code))
		   }
	}
	//cache up to 100 mb chunks
	chunk.data,err = read_data_buffer(file,chunk.size)
	if err !=nil{
		return nil,err
	}
	chunk.crc32,_ = read_uint32(file)
    return &chunk,nil   
}



func ReadChunks(file *os.File) (ChunkMap,error){


	is_png := check_signature(file)
	if !is_png{
		return nil,errors.New("NOT a valid PNG File")
	}
		
	chunks := make(map[string][]Chunk,0)
	
	for len(chunks[IEND]) == 0{
		current_chunk,err := read_chunk(file)
		if err != nil{
			return nil,err
		}
	
		chunks[string(current_chunk.code)] = append(chunks[string(current_chunk.code)], *current_chunk)

	}
	return chunks,nil
}







