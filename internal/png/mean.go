package png

type Histogram map[uint16]uint64

func calculateHistogram(imageData *ImageData) []Histogram {
	channels, _ := imageData.Metadata.GetSamples()
	histograms := make([]Histogram, channels)
	for _, row := range imageData.Pixels {
		for _, pixel := range row {
			values := pixel.Uint16()
			for iValue, value := range values {
				if histograms[iValue] == nil {
					histograms[iValue] = make(Histogram)
				}
				histograms[iValue][value]++
			}
		}
	}
	return histograms
}

func (self Histogram) mean() uint64 {
	sum := uint64(0)
	n := uint64(0)
	for k, v := range self {
		sum += uint64(k) * v
		n += v
	}
	return sum / n
}

func (self *ImageData) Mean() []uint64 {
	histograms := calculateHistogram(self)
	means := make([]uint64, 0)
	for _, h := range histograms {
		means = append(means, h.mean())
	}
	return means
}