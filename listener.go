package razor

import (
	"context"
	"net"
	"strconv"
	"strings"
	"syscall"
)

// razorListener implement the RazorListener can also be used as the net.Listener
type razorListener struct {
	ServerAddr [4]byte
	ServerPort int
	fd         int
	ctx        context.Context
	isInit     bool
}

func Listen(host string) (RazorListener, error) {
	r := &razorListener{}

	ss := strings.Split(host, ":")
	addr := net.ParseIP(ss[0]).To4()
	if addr == nil {
		return nil, ErrParseHost
	}
	copy(r.ServerAddr[:], addr)
	port, err := strconv.Atoi(ss[1])
	if err != nil {
		return nil, err
	}
	r.ServerPort = port

	sa := &syscall.SockaddrInet4{Addr: r.ServerAddr, Port: r.ServerPort}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		if err == syscall.ENOPROTOOPT {
			return nil, ErrTFONotSupport
		}
		return nil, err
	}
	r.fd = fd

	err = syscall.SetsockoptInt(r.fd, syscall.SOL_TCP, TCP_FASTOPEN, 1)
	if err != nil {
		return nil, err
	}

	err = syscall.Bind(r.fd, sa)
	if err != nil {
		return nil, err
	}

	err = syscall.Listen(r.fd, LISTEN_BACKLOG)
	if err != nil {
		return nil, err
	}
	r.isInit = false

	return r, nil
}

func (r *razorListener) Accept() (net.Conn, error) {
	if r.isInit {
		sa := &syscall.SockaddrInet4{Addr: r.ServerAddr, Port: r.ServerPort}

		fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
		if err != nil {
			return nil, err
		}
		r.fd = fd

		err = syscall.SetsockoptInt(r.fd, syscall.SOL_TCP, TCP_FASTOPEN, 1)
		if err != nil {
			return nil, err
		}

		err = syscall.Bind(r.fd, sa)
		if err != nil {
			return nil, err
		}

		err = syscall.Listen(r.fd, LISTEN_BACKLOG)
		if err != nil {
			return nil, err
		}
		r.isInit = false
	}

	rc := &razor{}

	cfd, sockaddr, err := syscall.Accept(r.fd)
	if err != nil {
		return nil, err
	}
	rc.fd = cfd
	rc.Addr = r.ServerAddr
	rc.Port = r.ServerPort

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
	return &net.TCPAddr{
		IP:   r.ServerAddr[:],
		Port: r.ServerPort,
	}
}
