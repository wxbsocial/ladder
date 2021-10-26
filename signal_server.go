package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
)

type SignalServer struct {
	clients map[string]net.Conn
	address string
	lock    sync.RWMutex
}

func NewSignalServer(
	address string,
) *SignalServer {
	return &SignalServer{
		clients: make(map[string]net.Conn),
		address: address,
	}
}

func (server *SignalServer) Start() {
	go server.startListen()
}

func (server *SignalServer) startListen() {

	s, err := net.Listen("tcp", server.address)
	if err != nil {
		fmt.Printf("signal server listen failed: %v\n", err)
		return
	}

	fmt.Printf("signal server listening: %v\n", server.address)

	for {
		client, err := s.Accept()
		if err != nil {
			fmt.Printf("signal server accept failed: %v\n", err)
			continue
		}

		key := server.getClientKey(client)

		fmt.Printf("signal server accept client:%v\n", key)

		server.lock.Lock()
		server.clients[key] = client
		server.lock.Unlock()
	}
}

func (server *SignalServer) getClientKey(client net.Conn) string {
	return client.RemoteAddr().String()
}

func (server *SignalServer) removeClient(client net.Conn) {
	defer client.Close()

	key := server.getClientKey(client)

	server.lock.Lock()
	defer server.lock.Unlock()

	delete(server.clients, key)

	fmt.Printf("signal server remove client:%v\n", key)

}

func (server *SignalServer) selectClient(clients map[string]net.Conn) (net.Conn, error) {
	server.lock.RLock()
	defer server.lock.RUnlock()

	if len(clients) == 0 {
		return nil, errors.New("not found available client")
	}

	keys := make([]string, len(clients))
	idx := 0
	for key := range clients {
		keys[idx] = key
		idx++
	}

	randKey := keys[rand.Intn(len(keys))]

	client := clients[randKey]

	return client, nil
}

func (server *SignalServer) RequestConnect() error {
	for {
		client, err := server.selectClient(server.clients)
		if err != nil {
			return err
		}

		if err := server.sendSignal(client, Create); err != nil {
			server.removeClient(client)
			continue
		}

		return nil

	}

}

func (server *SignalServer) sendSignal(client net.Conn, cmd Signal) error {
	_, err := client.Write([]byte{byte(cmd)})
	if err != nil {
		return err
	}
	return nil
}
