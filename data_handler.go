package file_tranport

type dataHandler struct {
	blockSize int
	bds       []blockData
	pb        *progressBar
}

type blockData struct {
	raw   []byte
	index int
}

func newClientDataHandler(data []byte, blockSize int) *dataHandler {
	blockCount := (len(data) + blockSize - 1) / blockSize
	bds := make([]blockData, blockCount)
	for i := 0; i < blockCount; i++ {
		bds[i] = blockData{
			raw:   data[i*blockSize : min((i+1)*blockSize, len(data))],
			index: i,
		}
	}

	return &dataHandler{
		blockSize: blockSize,
		bds:       bds,
		pb:        NewProgressBar(blockCount),
	}
}

func (d *dataHandler) ReadNextBlock() (bool, *blockData) {
	index := d.pb.FindNextUnset()
	if index == -1 {
		return false, nil
	}

	d.pb.Set(index)
	return true, &d.bds[index]
}

func (d *dataHandler) SetBlockCompleted(index int) {
	d.pb.Set(index)
}

func (d *dataHandler) IsAllBlockCompleted() bool {
	return d.pb.IsAllSet()
}

func newServerDataHandler(blockCount int) *dataHandler {
	return &dataHandler{
		bds: make([]blockData, blockCount),
		pb:  NewProgressBar(blockCount),
	}
}

func (d *dataHandler) WriteBlock(bd *blockData) {
	d.bds[bd.index] = *bd
	d.pb.Set(bd.index)
}

func (d *dataHandler) CombinedOne() []byte {
	total := make([]byte, 0)
	for _, bd := range d.bds {
		total = append(total, bd.raw...)
	}
	return total
}
