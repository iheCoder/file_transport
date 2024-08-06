package file_tranport

import (
	"fmt"
	"io"
	"net"
	"os"
)

type FileTransportClient struct {
	path string
}

const (
	fixedBlockSize = 1024
)

func NewFileTransportClient(path string) *FileTransportClient {
	return &FileTransportClient{
		path: path,
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

	// 2. 打开文件
	file, err := os.Open(c.path)
	if err != nil {
		fmt.Println("os.Open err:", err)
		return err
	}

	// 3. 读取所有数据，并构建数据处理器
	// TODO: 由于数据可能会很大，所以这里不适合一次性读取所有数据
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("io.ReadAll err:", err)
		return err
	}
	handler, err := newClientMemDataHandler(c.path, data, fixedBlockSize)
	if err != nil {
		fmt.Println("newClientMemDataHandler err:", err)
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
