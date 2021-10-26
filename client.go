package main

import (
	"fmt"
	"io"
	"net"
)

type Client struct {
	signalServerAddress string
	serverAddress       string
	socksServerAddress  string
	key                 []byte
}

func NewClient(
	signalServerAddress string,
	serverAddress string,
	socksAddress string,
	key []byte,
) *Client {
	return &Client{
		signalServerAddress: signalServerAddress,
		serverAddress:       serverAddress,
		socksServerAddress:  socksAddress,
		key:                 key,
	}
}

func (client *Client) Start() {
	go client.startMoniter()
}

func (client *Client) startMoniter() {
	for {

		signalServer, err := net.Dial("tcp", client.signalServerAddress)
		if err != nil {
			fmt.Printf("connect server failed:%v\n", err)
			continue
		}

		fmt.Printf("monitoring :%v\n", client.signalServerAddress)
		client.monitor(signalServer)
	}

}

func (client *Client) monitor(signalServer net.Conn) {

	for {
		buf := make([]byte, 1)

		n, err := io.ReadFull(signalServer, buf[:1])
		if err != nil {
			fmt.Printf("monitor failed:%v\n", err)
			return
		}

		fmt.Printf("monitor read:%v\n", n)

		signal := Signal(buf[0])
		switch signal {
		case Create:
			{
				go client.doCreateConnect()
			}
		}

	}
}

func (client *Client) doCreateConnect() {
	serverConn, err := client.connectServer()
	if err != nil {
		fmt.Printf("connect to server failed:%v\n", err)
		return
	}

	fmt.Printf("connect to server :%v\n", serverConn.RemoteAddr().String())

	socksConn, err := client.connectSocksServer()
	if err != nil {
		fmt.Printf("connect to socks failed:%v\n", err)

		return
	}

	fmt.Printf("connect to socks :%v\n", socksConn.RemoteAddr().String())

	fmt.Printf("bridge:%v-%v\n", serverConn.RemoteAddr().String(), socksConn.RemoteAddr().String())

	wrapServerConn, err := NewWrapConn(client.key, serverConn)
	if err != nil {
		fmt.Printf("wrap conn failed:%v\n", err)
		return
	}

	client.bridge(wrapServerConn, socksConn)
}

func (client *Client) bridge(src, dst net.Conn) {
	forward := func(src, dst net.Conn) {
		defer src.Close()
		defer dst.Close()
		n, err := io.Copy(dst, src)
		if err != nil {
			fmt.Printf("forward failed:%v.copy bytes:%v\n", err, n)
		}

		fmt.Printf("close bridge(from %v to %v).copy bytes:%v\n", src.RemoteAddr().String(), dst.RemoteAddr().String(), n)

	}
	go forward(src, dst)
	go forward(dst, src)
}

func (client *Client) connectServer() (net.Conn, error) {
	return net.Dial("tcp", client.serverAddress)
}

func (client *Client) connectSocksServer() (net.Conn, error) {
	return net.Dial("tcp", client.socksServerAddress)
}
