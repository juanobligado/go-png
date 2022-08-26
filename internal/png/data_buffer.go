package png

import "bytes"

// Abstract chunk data to handle big chunks in the future
type IByteBuffer interface{
	GetRawData() []byte
	GetReader() *bytes.Reader
}

type MemoryBuffer struct{
	data []byte
}


func (memoryData *MemoryBuffer) GetRawData() []byte{
	return memoryData.data
}

func (self *MemoryBuffer) GetReader() *bytes.Reader{
	return bytes.NewReader(self.data)
}

