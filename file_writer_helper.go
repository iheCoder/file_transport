package file_tranport

import "os"

// NewFileWriterHelper create a new file writer helper
func NewFileWriterHelper(path string) (*fileHelper, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &fileHelper{
		f: f,
	}, nil
}

func (f *fileHelper) InitWriterFileHelper(blockSize int) {
	f.blockSize = blockSize
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

// WriteAll write all data to file
func (f *fileHelper) WriteAll(data []byte) error {
	_, err := f.f.Write(data)
	return err
}
