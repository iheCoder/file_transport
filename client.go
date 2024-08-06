package file_tranport

import (
	"fmt"
	"net"
)

var (
	defaultBandWidth        = 1024
	defaultBlockTransportMS = 1000
)

type FileTransportClient struct {
	path    string
	adviser *transportAdviser
}

const (
	fixedBlockSize = 1024
)

func NewFileTransportClient(path string) *FileTransportClient {
	return &FileTransportClient{
		path:    path,
		adviser: newTransportAdviser(float64(defaultBandWidth), defaultBlockTransportMS),
	}
}

func (c *FileTransportClient) request(ip, port string) error {
	// 客户端
	// 1. 连接服务端
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return err
	}
	defer conn.Close()

	// 3. 读取所有数据，并构建数据处理器
	handler, err := c.adviser.GetClientDataHandler(c.path)
	if err != nil {
		fmt.Println("GetClientDataHandler err:", err)
		return err
	}

	// 3. 读取文件内容并发送
	return c.sendData(conn, handler)
}

func (c *FileTransportClient) sendData(conn net.Conn, handle *dataHandler) error {
	for {
		ok, bd, _ := handle.ReadNextBlock()
		if !ok {
			break
		}

		bytes, err := MarshallTransportModel(bd)
		if err != nil {
			fmt.Println("json.Marshal err:", err)
			return err
		}

		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println("conn.Write err:", err)
			return err
		}
	}

	return nil
}
