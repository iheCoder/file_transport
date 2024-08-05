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
	defer file.Close()

	// 创建数据处理器
	handler := newServerDataHandler(blockCount)

	// 读取客户端发送的数据
	data, err := s.receiveData(conn, handler)
	if err != nil {
		fmt.Println("receiveData err:", err)
		return
	}

	// 写入文件
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("file.Write err:", err)
		return
	}

	fmt.Println("文件接收完毕")
}

func (s *FileTransportServer) receiveData(conn net.Conn, handler *dataHandler) ([]byte, error) {
	for {
		size, err := GetDataSizeFromConn(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		buf := make([]byte, size)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		block, err := newBlockData(buf[:n])
		if err != nil {
			return nil, err
		}
		handler.WriteBlock(block)
	}

	if !handler.IsAllBlockCompleted() {
		return nil, errDataBroken
	}

	return handler.CombinedOne(), nil
}
