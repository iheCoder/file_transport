package file_tranport

import "os"

type fileHelper struct {
	f         *os.File
	blockSize int
	fileSize  int
}

// NewFileHelper create a new file helper
func NewFileHelper(path string, blockSize int) (*fileHelper, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := int(stat.Size())
	return &fileHelper{
		f:         f,
		blockSize: blockSize,
		fileSize:  fileSize,
	}, nil
}

// Size get file size
func (f *fileHelper) Size() (int, error) {
	stat, err := f.f.Stat()
	if err != nil {
		return 0, err
	}

	return int(stat.Size()), nil
}

// ReadBlock read block from file
func (f *fileHelper) ReadBlock(index int) ([]byte, error) {
	raw := make([]byte, f.blockSize)
	_, err := f.f.Seek(int64(index*f.blockSize), 0)
	if err != nil {
		return nil, err
	}

	n, err := f.f.Read(raw)
	if err != nil {
		return nil, err
	}

	return raw[:n], nil
}

// WriteBlock write block to file
func (f *fileHelper) WriteBlock(index int, raw []byte) error {
	_, err := f.f.Seek(int64(index*f.blockSize), 0)
	if err != nil {
		return err
	}

	_, err = f.f.Write(raw)
	if err != nil {
		return err
	}
	return nil
}
