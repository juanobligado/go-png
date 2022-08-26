package png

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadPng(t *testing.T) {
	filename := "../../test_cases/test_mean_0_bw.png"
	data,err := ReadPng(filename)
	assert.Nil(t,err)

	mean := data.Mean()
	assert.NotNil(t,mean)
	
}

func TestImageWithFilters(t *testing.T) {
	filename := "../../test_cases/deer.png"
	data,err := ReadPng(filename)
	assert.Nil(t,err)

	mean := data.Mean()
	assert.NotNil(t,mean)
	
}

func TestAllRed(t *testing.T) {

	filename := "../../test_cases/test_all_red.png"
	data,err := ReadPng(filename)
	assert.Nil(t,err)

	mean := data.Mean()
	assert.Equal(t,uint64(255),mean[0])
	assert.Equal(t,uint64(0),mean[1])
	assert.Equal(t,uint64(0),mean[2])
	assert.Equal(t,uint64(255),mean[3])
	assert.NotNil(t,mean)

}

func TestAllBlack(t *testing.T) {

	filename := "../../test_cases/test_all_black.png"
	data,err := ReadPng(filename)
	assert.Nil(t,err)

	mean := data.Mean()
	assert.NotNil(t,mean)

}