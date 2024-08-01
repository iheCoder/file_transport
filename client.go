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

	// 3. 读取文件内容并发送
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("file.Read err:", err)
			return err
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("conn.Write err:", err)
			return err
		}
	}

	return nil
}
