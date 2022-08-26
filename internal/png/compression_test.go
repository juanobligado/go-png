package png

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestDecode(m *testing.T) {

	chunkmap := ChunkMap{}
	 buffer := MemoryBuffer{}
	 buffer.data =[]byte{0x18, 0x57, 0x63, 0x20, 0x19, 0x30, 0x30, 0x00, 0x00, 0x00, 0x34, 0x00,0x01}
	 chunk := Chunk{}
	 chunk.data = &buffer
	 chunkmap["IDAT"] = append(chunkmap["IDAT"], chunk)

	 decompressed,err := Decompress(chunkmap) 

	 assert.Nil(m,err )

	 assert.True(m,len(decompressed) == 52 )

	//	assert.Equal(m,uint32(16),chunk.size)

}
