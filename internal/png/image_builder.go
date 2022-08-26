package png

import (
	"bytes"
	"errors"
)

type ImageBuilder interface {
	BuildPixels([][]byte) (*ImageData,error)
}

type RGBImageBuilder[T Uint] struct {
	IHDR *IHDRChunk 
}

func (self *RGBImageBuilder[T]) BuildPixels(data [][]byte) (*ImageData,error) {

	reader := RgbReader[T]{HasAlpha:  self.IHDR.HasAlpha(),Bits: uint(self.IHDR.Bitdepth)}
	imageData := ImageData{ Metadata: self.IHDR}

	for _,row := range data{		
		new_pixels,err := reader.ReadPixels(bytes.NewReader(row),self.IHDR.Width)
		if err != nil{
			return nil,err
		}
		imageData.Pixels = append(imageData.Pixels,new_pixels)

	}
	return &imageData,nil
}

type GrayScaleImageBuilder[T Uint] struct {
	IHDR  *IHDRChunk
}

func (self *GrayScaleImageBuilder[T]) BuildPixels(data [][]byte) (*ImageData,error) {

	reader := GrayscaleReader[T]{HasAlpha:  self.IHDR.HasAlpha(),Bits: uint(self.IHDR.Bitdepth)}
	imageData := ImageData{ Metadata: self.IHDR}
	for _,row := range data{		
		new_pixels,err := reader.ReadPixels(bytes.NewReader(row),self.IHDR.Width)
		if err != nil{
			return nil,err
		}
		imageData.Pixels = append(imageData.Pixels,new_pixels)

	}
	return &imageData,nil

}

type PaletteImageBuilder struct {
	IHDR  *IHDRChunk
	PLTE  *Palette
}

func (self *PaletteImageBuilder) BuildPixels(data [][]byte) (*ImageData,error) {
	return nil,errors.New("Not Implemented")
}

func createImageBuilder(ihdr *IHDRChunk, palette *Palette) (ImageBuilder, error) {


	switch ihdr.Colortype {
		case IndexedColor:{
			return &PaletteImageBuilder{IHDR: ihdr,PLTE: palette},nil
		}
		case TrueColorWithAlpha:{
			switch ihdr.Bitdepth{
				case 8:{
					return &RGBImageBuilder[uint8]{ IHDR: ihdr},nil
				}
				case 16:{
					return &RGBImageBuilder[uint16]{IHDR: ihdr},nil
				}
			}

		}
		case TrueColor:
			{
				switch ihdr.Bitdepth{
					case 8:{
						return &RGBImageBuilder[uint8]{ IHDR: ihdr},nil
					}
					case 16:{
						return &RGBImageBuilder[uint16]{IHDR: ihdr},nil
					}
				}
			}
		case GreyscaleWithAlpha:
		case Greyscale:
			{
				switch ihdr.Bitdepth{
					case 8:{
						return &GrayScaleImageBuilder[uint8]{IHDR: ihdr},nil
					}
					case 16:{
						return &GrayScaleImageBuilder[uint16]{IHDR: ihdr},nil
					}
				}			
			}
	}

	return nil, errors.New("Not Supported Image Type")
}