package png

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestFilterFactory(m *testing.T) {

	 f,err := createInverseFilter(0,nil)

	 assert.Nil(m,err)
	 assert.NotNil(m,f)

}