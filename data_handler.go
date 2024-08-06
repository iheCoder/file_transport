package file_tranport

import (
	"encoding/json"
)

type dataHandler struct {
	blockSize int
	bds       []blockData
	pb        *progressBar

	// for server
	initialized bool
}

func newBlockData(raw []byte) (*blockData, error) {
	var bd blockData
	err := json.Unmarshal(raw, &bd)
	if err != nil {
		return nil, err
	}
	return &bd, nil
}

func newClientDataHandler(data []byte, blockSize int) *dataHandler {
	blockCount := (len(data) + blockSize - 1) / blockSize
	bds := make([]blockData, blockCount)
	for i := 0; i < blockCount; i++ {
		bds[i] = blockData{
			Raw:   data[i*blockSize : min((i+1)*blockSize, len(data))],
			Index: i,
			Count: blockCount,
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

func newServerDataHandler() *dataHandler {
	return &dataHandler{}
}

func (d *dataHandler) initServerDataHandler(blockCount int) {
	d.bds = make([]blockData, blockCount)
	d.pb = NewProgressBar(blockCount)
}

func (d *dataHandler) WriteBlock(bd *blockData) {
	if !d.initialized {
		d.initServerDataHandler(bd.Count)
		d.initialized = true
	}

	d.bds[bd.Index] = *bd
	d.pb.Set(bd.Index)
}

func (d *dataHandler) CombinedOne() []byte {
	total := make([]byte, 0)
	for _, bd := range d.bds {
		total = append(total, bd.Raw...)
	}
	return total
}
