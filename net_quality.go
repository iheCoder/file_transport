package file_tranport

import (
	"encoding/json"
	"io"
	"net"
	"time"
)

var bufferSize = 1024

type netQualityTestData struct {
	TestData []byte    `json:"test_data"`
	SendTime time.Time `json:"send_time"`
	dataSize int
}

type netQualityTestResponse struct {
	Bandwidth float64 `json:"bandwidth"`
}

type netQualityMeasurer struct {
	data netQualityTestData
}

func NewNetQualityMeasurer(dataSize int) *netQualityMeasurer {
	qtd := netQualityTestData{
		TestData: make([]byte, dataSize),
	}
	return &netQualityMeasurer{
		data: qtd,
	}
}

func (q *netQualityMeasurer) StartMeasure(conn net.Conn) error {
	q.wrapTestData()

	return q.sendTestData(conn)
}

func (q *netQualityMeasurer) sendTestData(conn net.Conn) error {
	dbytes, err := json.Marshal(q.data)
	if err != nil {
		return err
	}

	_, err = conn.Write(dbytes)
	return err
}

func (q *netQualityMeasurer) wrapTestData() {
	q.data.SendTime = time.Now()
}

func (q *netQualityMeasurer) EndMeasure(conn net.Conn) error {
	if err := q.receiveTestData(conn); err != nil {
		return err
	}

	resp, err := q.measureQuality()
	if err != nil {
		return err
	}

	return q.sendMeasureResp(conn, resp)
}

func (q *netQualityMeasurer) receiveTestData(conn net.Conn) error {
	// read data
	dataBytes := make([]byte, 0)
	buf := make([]byte, bufferSize)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		dataBytes = append(dataBytes, buf[:n]...)
	}

	// unmarshal data and fill dataSize
	err := json.Unmarshal(dataBytes, &q.data)
	q.data.dataSize = len(dataBytes)

	return err
}

func (q *netQualityMeasurer) measureQuality() (*netQualityTestResponse, error) {
	endTime := time.Now()
	duration := endTime.Sub(q.data.SendTime)
	resp := netQualityTestResponse{
		Bandwidth: float64(q.data.dataSize) / duration.Seconds(),
	}

	return &resp, nil
}

func (q *netQualityMeasurer) sendMeasureResp(conn net.Conn, resp *netQualityTestResponse) error {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = conn.Write(respBytes)
	return err
}

func (q *netQualityMeasurer) GetMeasureResp(conn net.Conn) (*netQualityTestResponse, error) {
	respBytes := make([]byte, 0)
	buf := make([]byte, bufferSize)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		respBytes = append(respBytes, buf[:n]...)
	}

	resp := &netQualityTestResponse{}
	err := json.Unmarshal(respBytes, resp)
	return resp, err
}
