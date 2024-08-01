package file_tranport

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

type FileTransportServer struct {
	downloadDir string
}

func NewFileTransportServer(downloadDir string) *FileTransportServer {
	return &FileTransportServer{
		downloadDir: downloadDir,
	}
}

func (s *FileTransportServer) server(ip, port string) error {
	// 服务端
	// 1. 创建监听
	ln, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	defer ln.Close()

	// 2. 阻塞等待客户端连接
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		// 3. 处理客户端请求
		go s.handleConn(conn)
	}
}

func (s *FileTransportServer) handleConn(conn net.Conn) {
	defer conn.Close()

	// 创建接收文件
	file, err := os.Create(filepath.Join(s.downloadDir, "server_test.txt"))
	if err != nil {
		fmt.Println("os.Create err:", err)
		return
	}

	// 读取客户端发送的数据
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("conn.Read err:", err)
			return
		}

		// 写入文件
		_, err = file.Write(buf[:n])
		if err != nil {
			fmt.Println("file.Write err:", err)
			return
		}
	}

	fmt.Println("文件接收完毕")
}
