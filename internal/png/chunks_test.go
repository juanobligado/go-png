package png

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Complete the 'extraLongFactorials' function below.
 *
 * The function accepts INTEGER n as parameter.
 */

func TestCheckSignature(m *testing.T) {
	filename := "../../test_cases/deer.png"
	file, err := os.OpenFile(filename,os.O_RDONLY , 0444)
	if err != nil{
		m.Error(err)
	}
	defer file.Close()
	is_png := check_signature(file)
	assert.True(m,is_png);

}

func TestReadChunk(m *testing.T) {
	filename := "../../test_cases/test_all_black.png"
	file, err := os.OpenFile(filename,os.O_RDONLY , 0444)
	if err != nil{
		m.Error(err)
	}
	defer file.Close()

	chunks,err := ReadChunks(file)


	assert.True(m,len(chunks["IHDR"]) > 0)
	assert.True(m,len(chunks["IDAT"]) > 0)
	assert.True(m,len(chunks["IEND"]) > 0)

	//	assert.Equal(m,uint32(16),chunk.size)

}


