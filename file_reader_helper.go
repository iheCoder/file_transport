package file_tranport

import "os"

type fileHelper struct {
	f         *os.File
	blockSize int
	fileSize  int
}

// NewFileReaderHelper create a new file helper
func NewFileReaderHelper(path string, blockSize int) (*fileHelper, error) {
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

func GetFileSize(path string) (int, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return int(stat.Size()), nil
}

// Size get file size
func (f *fileHelper) Size() int {
	return f.fileSize
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

// Close close file
func (f *fileHelper) Close() error {
	return f.f.Close()
}

// ReadAll read all data from file
func (f *fileHelper) ReadAll() ([]byte, error) {
	data := make([]byte, f.fileSize)
	_, err := f.f.Read(data)
	return data, err
}
