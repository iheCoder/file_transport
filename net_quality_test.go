package file_tranport

import (
	"io"
	"net"
	"testing"
	"time"
)

type mockNetConn struct {
	data  []byte
	index int
}

func (m *mockNetConn) Read(b []byte) (n int, err error) {
	if m.index >= len(m.data) {
		return 0, io.EOF
	}

	n = copy(b, m.data[m.index:])
	m.index += n
	return n, nil
}

func (m *mockNetConn) Write(b []byte) (n int, err error) {
	m.data = append(m.data, b...)
	return len(b), nil
}

func (m *mockNetConn) Close() error {
	return nil
}

func (m *mockNetConn) LocalAddr() net.Addr {
	return nil
}

func (m *mockNetConn) RemoteAddr() net.Addr {
	return nil
}

func (m *mockNetConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockNetConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockNetConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestNetQualityMeasurer_StartMeasure(t *testing.T) {
	m := NewNetQualityMeasurer(1024)
	conn := &mockNetConn{}
	err := m.StartMeasure(conn)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	err = m.EndMeasure(conn)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := m.GetMeasureResp(conn)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Bandwidth)
}
