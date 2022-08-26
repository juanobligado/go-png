package png

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestReadRBG(m *testing.T) {

	 buffer := MemoryBuffer{}
	 buffer.data =[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	 
	 reader := RgbReader[uint8]{HasAlpha: false, Bits: 8}
	 data,err := reader.ReadPixels( buffer.GetReader(),3)

	 assert.Nil(m,err)
	 assert.Equal(m,3,len(data))
	 assert.Equal(m,3,len(data[0].Uint16()))

}




func TestGrayscale1(m *testing.T) {

	buffer := MemoryBuffer{}
	buffer.data =[]byte{0x0f}
	
	reader := GrayscaleReader[uint8]{HasAlpha: false, Bits: 4}
	data,err := reader.ReadPixels( buffer.GetReader(),2)

	assert.Nil(m,err)
	assert.Equal(m,2,len(data))

}

func TestGrayscale2(m *testing.T) {

	buffer := MemoryBuffer{}
	buffer.data =[]byte{0xf0}
	
	reader := GrayscaleReader[uint8]{HasAlpha: false, Bits: 4}
	data,err := reader.ReadPixels( buffer.GetReader(),2)

	assert.Nil(m,err)
	assert.Equal(m,2,len(data))

}
