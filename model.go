package file_tranport

import (
	"encoding/binary"
	"encoding/json"
	"net"
)

func MarshallTransportModel(a any) ([]byte, error) {
	bs, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	size := uint32(len(bs))

	return append(int32ToBytes(size), bs...), nil
}

func GetDataSizeFromConn(conn net.Conn) (int, error) {
	var size uint32
	err := binary.Read(conn, binary.BigEndian, &size)
	if err != nil {
		return 0, err
	}

	return int(size), nil
}

func int32ToBytes(n uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, n)
	return bytes
}

type blockData struct {
	Raw       []byte `json:"raw"`
	Index     int    `json:"index"`
	Count     int    `json:"count"`
	BlockSize int    `json:"blockSize"`
}
