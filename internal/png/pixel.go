package png

import (
	"bytes"
	"math/bits"
	"reflect"
)

type Uint interface {
     uint | uint8 | uint16 | uint32 | uint64 
}

type Palette struct{}



type PixelValue []uint32
type Pixel interface{
	Uint16() []uint16
}


type GrayscalePixel[T Uint ] struct {
	data []T
}
func (self *GrayscalePixel[T])  Uint16()[]uint16{
	data := make([]uint16, 0)
	for _,v := range self.data{
		data= append(data,uint16(v) )
	}
	return data
} 

type RGBPixel[T Uint ] struct {
	rgb_data []T
}
func (self *RGBPixel[T])  Uint16()[]uint16{
	data := make([]uint16, 0)
	for _,v := range self.rgb_data{
		data= append(data,uint16(v) )
	}
	return data
} 

func CreateRGBA[T Uint]( r T,g T,b T,a T) *RGBPixel[T]{
	rgba := RGBPixel[T]{}
	rgba.rgb_data = append(rgba.rgb_data,r,g,b,a)
	return &rgba
}

func CreateRGB[T Uint]( r T,g T,b T) *RGBPixel[T]{
	rgba := RGBPixel[T]{}
	rgba.rgb_data = append(rgba.rgb_data,r,g,b)
	return &rgba
}

func CreateGrayscalePixel[T Uint]( g T) *GrayscalePixel[T]{
	pixel := GrayscalePixel[T]{}
	pixel.data = append(pixel.data,g)
	return &pixel
}

func CreateGrayscaleWithAlphaPixel[T Uint]( g T,a T) *GrayscalePixel[T]{
	pixel := GrayscalePixel[T]{}
	pixel.data = append(pixel.data,g,a)
	return &pixel
}



type PixelReader[T Uint] interface{
	ReadPixels(reader *bytes.Buffer, n uint32) ([]Pixel,error)
}



type RgbReader[T Uint] struct {
	HasAlpha bool
	Bits uint
}
func    (self *RgbReader[T])   ReadPixels( reader *bytes.Reader, n uint32) ([]Pixel,error){
	
	pixels := make([]Pixel, 0)
	for i:=uint32(0);i<n;i++{
		r,err :=  read[T](reader)
		if err != nil{
			return nil, err
		}
		g,err := read[T](reader)
		if err != nil{
			return nil, err
		}
		
		b,err := read[T](reader)
		if err != nil{
			return nil, err
		}
		if !self.HasAlpha{
			pixel := CreateRGB[T](r,g,b)
			pixels = append(pixels, pixel)
		}else{
			a,err := read[T](reader)
			if err != nil{
				return nil, err
			}
			pixels = append(pixels, CreateRGBA[T](r,g,b,a))
		}
	}
	return pixels,nil

} 


type GrayscaleReader[T Uint] struct {
	HasAlpha bool
	Bits uint
}
func    (self *GrayscaleReader[T])   ReadPixels( reader *bytes.Reader, n uint32) ([]Pixel,error){
	
	pixels := make([]Pixel, 0)

	mask := byte(0) 
	for i := 0 ; i < int(self.Bits) ; i++{
		mask = (mask << 1) + 1
	}
	// MSB first
	mask = bits.Reverse8(mask)
	for i:=uint32(0);i<n;i++{
		p,err :=  read[T](reader)
		if err != nil{
			return nil, err
		}
		if self.Bits >= 8 {
			if !self.HasAlpha{
				pixel := CreateGrayscalePixel(p)
				pixels = append(pixels, pixel)
			}else{
				a,err := read[T](reader)
				if err != nil{
					return nil, err
				}
				pixel := CreateGrayscaleWithAlphaPixel(p,a)
				pixels = append(pixels, pixel)
			}
	
		}else{
			byte_pixels := 8 / self.Bits
			encoded := byte(p)
			for s:=0 ; s < int(byte_pixels); s++{
				pixel_value :=   encoded & mask
				pixel_value =  bits.RotateLeft8(pixel_value,int(self.Bits))   
				pixel := CreateGrayscalePixel(pixel_value)
				pixels = append(pixels, pixel)
				encoded = encoded << uint8(self.Bits)

				i++
				if i  == n{
					return pixels,nil
				}
			}
		}
	}
	return pixels,nil

} 


func  read[T Uint](reader *bytes.Reader) (T, error) {
    zero := T(0)
    tType := reflect.TypeOf(zero)
	
	if tType.Size() == 1 {
		b , err := reader.ReadByte()
		if err != nil{
			return 0, err
		} 
		return T(b),nil
	}

	buffer := make([]byte, tType.Size() )

    _,err := reader.Read(buffer)
    if err != nil{
        return 0,err
    }
	return T(to_uint16(buffer)),nil
}


