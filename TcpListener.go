package main

import (
	"fmt"
	"net"
)

type TcpListener struct {
	OnAccept func(con net.Conn)
	OnError  func(msg *string)
	listener net.Listener
	NoDelay  bool
}

func (tcp *TcpListener) handleAccept() {
	for {
		client, err := tcp.listener.Accept()
		if err != nil {
			msg := fmt.Sprintf("Accept error: %v", err)
			if tcp.OnError != nil {
				tcp.OnError(&msg)
			}

			continue
		}

		if tcp.NoDelay {
			if tcpConn, ok := client.(*net.TCPConn); ok {
				err := tcpConn.SetNoDelay(true)
				if err != nil {
					fmt.Println("TCP_NODELAY set error:", err)
					tcpConn.Close()
					continue
				}
			}
		}

		if tcp.OnAccept != nil {
			go func() {
				tcp.OnAccept(client)
				fmt.Println("client done")
			}()
		}
	}
}

func (tcp *TcpListener) Start(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tcp.listener = listener

	go tcp.handleAccept()
	return nil
}

func (tcp *TcpListener) Stop() {
	if tcp.listener != nil {
		tcp.listener.Close()
	}
}
