package file_tranport

import "testing"

func TestGenPartialFileHash(t *testing.T) {
	fh, err := NewFileReaderHelper("testdata/min.txt", fixedBlockSize)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := fh.GenPartialFileHash()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
}
