package file_tranport

import (
	"encoding/json"
)

type blockDataMode int32

const (
	memoryMode blockDataMode = iota
	fileMode
)

type dataHandler struct {
	blockSize int
	bds       []blockData
	pb        *progressBar
	mode      blockDataMode

	// for server
	initialized bool
	// for file mode
	fh *fileHelper
}

func newBlockData(raw []byte) (*blockData, error) {
	var bd blockData
	err := json.Unmarshal(raw, &bd)
	if err != nil {
		return nil, err
	}
	return &bd, nil
}

func newClientMemDataHandler(path string, data []byte, blockSize int) (*dataHandler, error) {
	fh, err := NewFileHelper(path, blockSize)
	if err != nil {
		return nil, nil
	}

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
		mode:      memoryMode,
		fh:        fh,
	}, nil
}

func newClientFileDataHandler(path string, blockSize int) (*dataHandler, error) {
	fh, err := NewFileHelper(path, blockSize)
	if err != nil {
		return nil, nil
	}

	blockCount := (fh.fileSize + blockSize - 1) / blockSize
	return &dataHandler{
		blockSize: blockSize,
		fh:        fh,
		pb:        NewProgressBar(blockCount),
		mode:      fileMode,
	}, nil
}

func (d *dataHandler) ReadNextBlock() (bool, *blockData, error) {
	index := d.pb.FindNextUnset()
	if index == -1 {
		return false, nil, nil
	}

	if err := d.fillRawBlock(index); err != nil {
		return false, nil, err
	}

	d.pb.Set(index)
	return true, &d.bds[index], nil
}

func (d *dataHandler) fillRawBlock(index int) error {
	if d.mode != fileMode {
		return nil
	}

	raw, err := d.fh.ReadBlock(index)
	if err != nil {
		return err
	}

	d.bds[index].Raw = raw
	return nil
}

func (d *dataHandler) SetBlockCompleted(index int) {
	d.pb.Set(index)
	// release memory
	d.bds[index].Raw = nil
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

func (d *dataHandler) WriteBlock(bd *blockData) error {
	if !d.initialized {
		d.initServerDataHandler(bd.Count)
		d.initialized = true
	}

	return d.writeBlock(bd)
}

func (d *dataHandler) writeBlock(bd *blockData) error {
	switch d.mode {
	case fileMode:
		return d.fh.WriteBlock(bd.Index, bd.Raw)
	default:
		d.bds[bd.Index] = *bd
		d.pb.Set(bd.Index)
	}

	return nil
}

func (d *dataHandler) CombinedOne() []byte {
	total := make([]byte, 0)
	for _, bd := range d.bds {
		total = append(total, bd.Raw...)
	}
	return total
}

func (d *dataHandler) Persist() error {
	if d.mode == fileMode {
		return nil
	}

	return d.fh.WriteBlock(0, d.CombinedOne())
}
