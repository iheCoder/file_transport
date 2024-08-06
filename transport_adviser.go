package file_tranport

type transportAdviser struct {
	bandwidth        float64
	totalSize        int
	blockTransportMS int
}

var (
	// big file boundary is 100MB
	bigFileBoundary = 100 * 1024 * 1024
	// max transport ms per block is 10s
	maxTransportSecondsPerBlock = 10
)

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

func (t *transportAdviser) GetClientDataHandler(path string) (*dataHandler, error) {
	// get file size
	size, err := GetFileSize(path)
	if err != nil {
		return nil, err
	}

	// get good block size
	blockSize := t.GoodBlockSize()

	// create new client data handler
	if float64(size)/t.bandwidth > float64(maxTransportSecondsPerBlock) {
		return newClientFileDataHandler(path, blockSize)
	} else {
		return newClientMemDataHandler(path, blockSize)
	}
}
