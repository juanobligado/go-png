package png

import (
	"bytes"
	"errors"
)

const (
	None  =byte(0)
	Sub = byte(1)
	Up = byte(2)
	Average = byte(3)
	Paeth = byte(4)
)



type Filter interface {
	Process([]byte,[][]byte) ([]byte, error)
}

func get_a(i uint32,minBytesPerPixel uint32,data []byte ) byte{
	if i < minBytesPerPixel{
		return 0
	}else{
		return data[i - minBytesPerPixel]
	}
}

func get_b(i uint32,line uint32,minBytesPerPixel uint32,data [][]byte) byte{
		return data[line-1][i]
}

func get_c(i uint32,line uint32,minBytesPerPixel uint32,data [][]byte) byte{
	if i < minBytesPerPixel{
		return 0
	}else{
		before_bLine := line-1
		before_bColumn := i - minBytesPerPixel
		return get_b(before_bColumn,before_bLine,minBytesPerPixel,data)
	}
}

func (ihdr *IHDRChunk) FilterStep() uint32 {
	samples,_ := ihdr.GetSamples()
	if ihdr.Bitdepth == 16{
		samples = samples*2
	}
	return samples
}

type IdentityFilter struct{}

func (f *IdentityFilter) Process(data []byte,processed [][]byte) ([]byte, error) {
	return data, nil
}

type UndoSubFilter struct{
	ihdr *IHDRChunk
}

func (f *UndoSubFilter) Process(data []byte,processed [][]byte ) ([]byte, error) {
	if  len(data) == 0{
		return nil, errors.New("No Data To Filter")
	}
	step := f.ihdr.FilterStep()
	recon := make([]byte, len(data))
	for i:=uint32(0) ; i < uint32(len(data)) ; i++{
		a := get_a(i,step,recon)
		recon[i] = data[i] + a  
	}
	return recon, nil
}

type UndoUpFilter struct{
	ihdr *IHDRChunk
}
func (self *UndoUpFilter) Process(data []byte,processed [][]byte) ([]byte, error) {
	if  len(data) == 0{
		return nil, errors.New("No Data To Filter")
	}
	step := self.ihdr.FilterStep()
	recon := make([]byte, len(data))	
	line := uint32(len(processed))
	for i:=uint32(0) ; i< uint32(len(data)) ; i++{
		
		reconb := get_b(i,line,step,processed)		
		recon[i] = data[i] + reconb  
	}
	return recon, nil
}

type UndoAverageFilter struct{
	ihdr *IHDRChunk

}
func (self *UndoAverageFilter) Process(data []byte,processed [][]byte) ([]byte, error) {
	if  len(data) == 0{
		return nil, errors.New("No Data To Filter")
	}
	recon := make([]byte, len(data))
	step := self.ihdr.FilterStep()
	line := uint32(len(processed))

	for i:= uint32(0) ; i< uint32(len(data)) ; i++{
		recona := get_a(i,step,data)
		reconb := get_b(i,line,step,processed)	
		recon[i] = data[i] +  (recona  + reconb)/2
	}
	return recon, nil

}


type UndoPaethFilter struct{
	ihdr *IHDRChunk
}

func abs( a int) int{
	if a < 0{
		return -a
	}
	return a
}

func (self *UndoPaethFilter) PaethPredictor(a byte,b byte,c byte) byte{
	p := int(a) + int(b) - int(c)
    pa := abs( p - int(a))
    pb := abs(p - int(b))
    pc := abs(p - int(c))
    Pr := c
	if pa <= pb && pa <= pc { 
		Pr = a
	}else if pb <= pc {
		Pr = b
	} 
    return Pr
}
func (self *UndoPaethFilter) Process(data []byte,processed [][]byte) ([]byte, error) {
	if  len(data) == 0{
		return nil, errors.New("No Data To Filter")
	}
	recon := make([]byte, len(data))
	step := self.ihdr.FilterStep()
	line := uint32(len(processed))

	for i:=uint32(0) ; i< uint32(len(data)) ; i++{
		a := get_a(i,step,data)
		b := get_b(i,line,step,processed)	
		c := get_c(i,line,step,processed)
		recon[i] = data[i] +  self.PaethPredictor(a,b,c)
	}
	return recon, nil

}



func createInverseFilter(b byte,ihdr *IHDRChunk) (Filter, error) {

	switch b {
		// No Filtering
		case None:
			return &IdentityFilter{}, nil
		case Sub:
			return &UndoSubFilter{ihdr: ihdr},nil// Sub
		case Up:
			return &UndoUpFilter{ihdr: ihdr},nil
		case Average:
			return &UndoAverageFilter{ihdr: ihdr},nil
		case Paeth:
			return &UndoPaethFilter{ihdr: ihdr},nil
			
	}
	return nil, errors.New("Not Supported Filtering Type")
}


func Unfilter( filtered_bytes []byte,  ihdr *IHDRChunk ) ([][]byte,error) {

	
	line_size,err :=  ihdr.GetBytesPerLine() 
	if err != nil{
		return nil,err
	}
	// Add Filter line header byte
	line_size++

	buffer := bytes.NewBuffer(filtered_bytes)	
	nLines := len(filtered_bytes) / int(line_size)
	if len(filtered_bytes) % int(line_size) != 0 {
		nLines++
	}

	image := make([][]byte, 0)
	for iLine:=0 ; iLine < nLines;iLine++ {
		filter_type,_ := buffer.ReadByte()	
		filtered_line := make([]byte,line_size-1)
		_,err := buffer.Read(filtered_line)
		if err != nil{
			return nil,err
		}
		
		reverse_filter,err := createInverseFilter(filter_type,ihdr)
		if err != nil{
			return nil,err
		}

		unfiltered_line,err := reverse_filter.Process(filtered_line,image)		
		if err != nil{
			return nil,err
		}
		image = append(image, unfiltered_line)
	}
	return image,nil
}