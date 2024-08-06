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

func newClientMemDataHandler(path string, blockSize int) (*dataHandler, error) {
	fh, err := NewFileReaderHelper(path, blockSize)
	if err != nil {
		return nil, nil
	}

	data, err := fh.ReadAll()
	if err != nil {
		return nil, err
	}

	blockCount := calculateBlockCount(len(data), blockSize)
	bds := make([]blockData, blockCount)
	for i := 0; i < blockCount; i++ {
		bds[i] = blockData{
			Raw:       data[i*blockSize : min((i+1)*blockSize, len(data))],
			Index:     i,
			Count:     blockCount,
			BlockSize: blockSize,
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
	fh, err := NewFileReaderHelper(path, blockSize)
	if err != nil {
		return nil, nil
	}

	blockCount := calculateBlockCount(fh.Size(), blockSize)

	// init block data
	// no need to fill raw data, because it will be filled when read
	bds := make([]blockData, blockCount)
	for i := 0; i < blockCount; i++ {
		bds[i] = blockData{
			Index:     i,
			Count:     blockCount,
			BlockSize: blockSize,
		}
	}

	return &dataHandler{
		blockSize: blockSize,
		fh:        fh,
		pb:        NewProgressBar(blockCount),
		mode:      fileMode,
		bds:       bds,
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

func calculateBlockCount(fileSize, blockSize int) int {
	return (fileSize + blockSize - 1) / blockSize
}
