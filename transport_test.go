package file_tranport

import (
	"testing"
	"time"
)

func TestSimpleServer(t *testing.T) {
	server := NewFileTransportServer("testdata")
	go server.server("localhost", "8888")

	t.Log("server started")

	client := NewFileTransportClient("testdata/test.txt")
	err := client.request("localhost", "8888")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("client request done")

	time.Sleep(10 * time.Second)
}
