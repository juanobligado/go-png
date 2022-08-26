package png

import (
	"errors"
	"os"
)

type ImageData struct {
	Metadata *IHDRChunk
	Pixels [][]Pixel
}


func ReadPng(filename string) (*ImageData , error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		return nil,err
	}
	defer file.Close()

	chunks, err := ReadChunks(file)
	if err != nil{
		return nil,err
	}
	filtered,err := Decompress(chunks)
	if err != nil{
		return nil,err
	}
	ihdr,err := chunks.IHDR()
	if err != nil{
		return nil,err
	}
	if ihdr.Interlace != 0{
		return nil,errors.New("Interlacing not supported yet")
	}
	plte,err := chunks.PLTE()
	if err != nil{
		return nil,err
	}
	unfiltered_bytes,err := Unfilter(filtered,ihdr)
	if err != nil{ 
		return nil,err
	}
	
	image_builder,err := createImageBuilder(ihdr,plte)
	if err !=nil{
		return nil,err
	}
	image_data,err := image_builder.BuildPixels(unfiltered_bytes)
	if err !=nil{
		return nil,err
	}
	return image_data,nil
}