package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type Server struct {
	innerListenAddress string
	outerListenAddress string
	outerConnChan      chan net.Conn
	signalServer       *SignalServer
	key                []byte
}

func NewServer(
	innerListenAddress string,
	outerListenAddress string,
	signalServer *SignalServer,
	key []byte,
) *Server {
	return &Server{
		innerListenAddress: innerListenAddress,
		outerListenAddress: outerListenAddress,
		signalServer:       signalServer,
		outerConnChan:      make(chan net.Conn, 1),
		key:                key,
	}
}

func (s *Server) Start() {
	go s.startListenInner()
	go s.startListenOuter()
}

func (s *Server) startListenInner() {
	server, err := net.Listen("tcp", s.innerListenAddress)
	if err != nil {
		fmt.Printf("listen inner failed: %v\n", err)
		return
	}

	fmt.Printf("server listening inner: %v\n", s.innerListenAddress)

	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Printf("inner accept failed: %v\n", err)
			continue
		}
		fmt.Printf("server accept inner: %v\n", client.RemoteAddr().String())

		go s.processInnerRequest(client)
	}
}

func (s *Server) processInnerRequest(client net.Conn) {

	remote, err := s.connectRemote()
	if err != nil {
		fmt.Printf("connect remote failed:%v\n", err)
		client.Close()
		return
	}

	fmt.Printf("bridge:%v-%v\n", client.RemoteAddr().String(), remote.RemoteAddr().String())

	wrapRemoteConn, err := NewWrapConn(s.key, remote)
	if err != nil {
		fmt.Printf("wrap conn failed:%v\n", err)
		return
	}
	s.bridge(client, wrapRemoteConn)

}

func (s *Server) startListenOuter() {
	server, err := net.Listen("tcp", s.outerListenAddress)
	if err != nil {
		fmt.Printf("listen outer failed: %v\n", err)
		return
	}

	fmt.Printf("server listening outer: %v\n", s.outerListenAddress)

	for {
		outerConn, err := server.Accept()
		if err != nil {
			fmt.Printf("outer accept failed: %v\n", err)
			continue
		}

		fmt.Printf("outer accept client: %v\n", outerConn.RemoteAddr().String())
		go s.processOuterRequest(outerConn)
	}
}

func (s *Server) processOuterRequest(outerConn net.Conn) {

	s.outerConnChan <- outerConn

}

func (s *Server) bridge(src, dst net.Conn) {
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

func (s *Server) connectRemote() (net.Conn, error) {

	err := s.signalServer.RequestConnect()
	if err != nil {
		return nil, err
	}

	select {
	case outer := <-s.outerConnChan:
		{
			return outer, nil
		}

	case <-time.After(time.Second * 30):
		{
			return nil, errors.New("timeout")
		}
	}
}
