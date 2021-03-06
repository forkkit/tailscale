// Copyright (c) 2020 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safesocket

import (
	"context"
	"fmt"
	"net"
	"syscall"
)

func path(vendor, name string, port uint16) string {
	return fmt.Sprintf("127.0.0.1:%v", port)
}

func ConnCloseRead(c net.Conn) error {
	return c.(*net.TCPConn).CloseRead()
}

func ConnCloseWrite(c net.Conn) error {
	return c.(*net.TCPConn).CloseWrite()
}

// TODO(apenwarr): handle magic cookie auth
func Connect(cookie, vendor, name string, port uint16) (net.Conn, error) {
	p := path(vendor, name, port)
	pipe, err := net.Dial("tcp", p)
	if err != nil {
		return nil, err
	}
	return pipe, err
}

func setFlags(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET,
			syscall.SO_REUSEADDR, 1)
	})
}

// TODO(apenwarr): use named pipes instead of sockets?
//   I tried to use winio.ListenPipe() here, but that code is a disaster,
//   built on top of an API that's a disaster. So for now we'll hack it by
//   just always using a TCP session on a fixed port on localhost. As a
//   result, on Windows we ignore the vendor and name strings.
// TODO(apenwarr): handle magic cookie auth
func Listen(cookie, vendor, name string, port uint16) (net.Listener, uint16, error) {
	lc := net.ListenConfig{
		Control: setFlags,
	}
	p := path(vendor, name, port)
	pipe, err := lc.Listen(context.Background(), "tcp", p)
	if err != nil {
		return nil, 0, err
	}
	return pipe, uint16(pipe.Addr().(*net.TCPAddr).Port), err
}
