package file_tranport

import (
	"fmt"
	"io"
	"net"
	"path/filepath"
)

type FileTransportServer struct {
	downloadDir string
	adviser     *transportAdviser
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
	path := filepath.Join(s.downloadDir, "server_test.txt")

	// 创建数据处理器
	handler, err := newUninitServerDataHandler(path)
	if err != nil {
		fmt.Println("newUninitServerDataHandler err:", err)
		return
	}

	// 读取客户端发送的数据
	err = s.receiveData(conn, handler)
	if err != nil {
		fmt.Println("receiveData err:", err)
		return
	}

	// 持久化
	if err = handler.Persist(); err != nil {
		fmt.Println("handler.Persist err:", err)
		return
	}

	fmt.Println("文件接收完毕")
}

func (s *FileTransportServer) receiveData(conn net.Conn, handler *dataHandler) error {
	for {
		size, err := GetDataSizeFromConn(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		buf := make([]byte, size)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		block, err := newBlockData(buf[:n])
		if err != nil {
			return err
		}
		if err = handler.WriteBlock(block); err != nil {
			return err
		}
	}

	if !handler.IsAllBlockCompleted() {
		return errDataBroken
	}

	return nil
}
