package file_tranport

type transportAdviser struct {
	bandwidth        float64
	totalSize        int
	blockTransportMS int
}

func newTransportAdviser(bandwidth float64, totalSize, blockTransportMS int) *transportAdviser {
	return &transportAdviser{
		bandwidth:        bandwidth,
		totalSize:        totalSize,
		blockTransportMS: blockTransportMS,
	}
}

func (t *transportAdviser) GoodBlockSize() int {
	maxBlockSize := int(t.bandwidth * float64(t.blockTransportMS) / 1000)
	if maxBlockSize > t.totalSize {
		return t.totalSize
	}

	// for redundancy
	return maxBlockSize * 8 / 10
}

func (t *transportAdviser) UpdateTotalSize(total int) {
	t.totalSize = total
}

func (t *transportAdviser) UpdateBandwidth(bandwidth float64) {
	t.bandwidth = bandwidth
}
