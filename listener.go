package razor

import (
	"context"
	"fmt"
	"net"
	"syscall"
)

//type Listener interface {
//	// Accept waits for and returns the next connection to the listener.
//	Accept() (Conn, error)
//
//	// Close closes the listener.
//	// Any blocked Accept operations will be unblocked and return errors.
//	Close() error
//
//	// Addr returns the listener's network address.
//	Addr() Addr
//}

type razorListener struct {
	ServerAddr [4]byte
	ServerPort int
	fd         int
	ctx        context.Context
	isInit     bool
}

func (r *razorListener) Accept() (net.Conn, error) {
	if r.isInit {
		sa := &syscall.SockaddrInet4{Addr: r.ServerAddr, Port: r.ServerPort}

		fmt.Println("init the socket")
		fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
		if err != nil {
			return nil, err
			//if err == syscall.ENOPROTOOPT {
			//	fmt.Println("TCP Fast Open server support is unavailable (unsupported kernel)")
			//}
		}
		r.fd = fd

		fmt.Println("init the set the socket")
		err = syscall.SetsockoptInt(r.fd, syscall.SOL_TCP, TCP_FASTOPEN, 1)
		if err != nil {
			return nil, err
		}

		err = syscall.Bind(r.fd, sa)
		if err != nil {
			return nil, err
		}

		fmt.Println("listen the socket")
		err = syscall.Listen(r.fd, LISTEN_BACKLOG)
		if err != nil {
			return nil, err
		}
		r.isInit = false
	}

	rc := &razor{}
	fmt.Println("accept the conn")

	cfd, sockaddr, err := syscall.Accept(r.fd)
	if err != nil {
		return nil, err
	}
	rc.fd = cfd
	rc.ServerAddr = r.ServerAddr
	rc.ServerPort = r.ServerPort

	if raddr, ok := sockaddr.(*syscall.SockaddrInet4); ok {
		rc.RAddr = raddr.Addr
		rc.RPort = raddr.Port
	}
	return rc, nil
}

func (r *razorListener) Close() error {
	err := syscall.Shutdown(r.fd, syscall.SHUT_RDWR)
	if err != nil {
		return err
	}
	err = syscall.Close(r.fd)
	if err != nil {
		return err
	}
	return nil
}

func (r *razorListener) Addr() net.Addr {
	return nil
}