package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"golang.org/x/crypto/chacha20"
)

type WrapConn struct {
	innerConn net.Conn
	key       []byte
	encoder   *chacha20.Cipher
	decoder   *chacha20.Cipher
}

func NewWrapConn(key []byte, conn net.Conn) (*WrapConn, error) {
	s := &WrapConn{
		key:       key, // should be exactly 32 bytes
		innerConn: conn,
	}

	var err error
	nonce := make([]byte, chacha20.NonceSizeX)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// s.encoder, err = chacha20.NewUnauthenticatedCipher(s.key, nonce)
	// if err != nil {
	// 	return nil, err
	// }

	s.encoder, err = chacha20.NewUnauthenticatedCipher(s.key, s.key[:24])
	if err != nil {
		return nil, err
	}

	// if n, err := s.innerConn.Write(nonce); err != nil || n != len(nonce) {
	// 	return nil, errors.New("write nonce failed: " + err.Error())
	// }

	return s, nil
}

func (c *WrapConn) Read(b []byte) (int, error) {

	if c.decoder == nil {
		// nonce := make([]byte, chacha20.NonceSizeX)
		// if n, err := io.ReadAtLeast(c.innerConn, nonce, len(nonce)); err != nil || n != len(nonce) {
		// 	return n, fmt.Errorf("can't read nonce from stream:%v,read bytes:%v ", err, n)
		// }
		// decoder, err := chacha20.NewUnauthenticatedCipher(c.key, nonce)
		// if err != nil {
		// 	return 0, errors.New("generate decoder failed: " + err.Error())
		// }
		// c.decoder = decoder

		decoder, err := chacha20.NewUnauthenticatedCipher(c.key, c.key[:24])
		if err != nil {
			return 0, errors.New("generate decoder failed: " + err.Error())
		}
		c.decoder = decoder
		fmt.Print("decoder generated\n")
	}

	n, err := c.innerConn.Read(b)
	if err != nil || n == 0 {
		return n, err
	}

	dst := make([]byte, n)
	pn := b[:n]
	c.decoder.XORKeyStream(dst, pn)
	copy(pn, dst)
	return n, nil

}

func (c *WrapConn) Write(b []byte) (n int, err error) {
	dst := make([]byte, len(b))
	c.encoder.XORKeyStream(dst, b)
	return c.innerConn.Write(dst)

}

func (c *WrapConn) Close() error {
	return c.innerConn.Close()

}
func (c *WrapConn) LocalAddr() net.Addr {
	return c.innerConn.LocalAddr()

}
func (c *WrapConn) RemoteAddr() net.Addr {
	return c.innerConn.RemoteAddr()

}
func (c *WrapConn) SetDeadline(t time.Time) error {
	return c.innerConn.SetDeadline(t)

}
func (c *WrapConn) SetReadDeadline(t time.Time) error {
	return c.innerConn.SetReadDeadline(t)
}
func (c *WrapConn) SetWriteDeadline(t time.Time) error {
	return c.innerConn.SetWriteDeadline(t)
}
