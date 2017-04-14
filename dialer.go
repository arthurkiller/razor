package razor

import (
	"net"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type razor struct {
	Addr  [4]byte
	Port  int
	RAddr [4]byte
	RPort int
	fd    int
}

func Dial(host string) (Razor, error) {
	rc := &razor{}

	ss := strings.Split(host, ":")
	addr := net.ParseIP(ss[0]).To4()
	if addr == nil {
		return nil, ErrParseHost
	}

	copy(rc.Addr[:], addr)

	port, err := strconv.Atoi(ss[1])
	if err != nil {
		return nil, err
	}
	rc.Port = port

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	rc.fd = fd

	return rc, nil
}

func (r *razor) Read(b []byte) (n int, err error) {
	n, err = syscall.Read(r.fd, b)
	return
}

func (r *razor) Write(b []byte) (n int, err error) {
	sa := &syscall.SockaddrInet4{Addr: r.Addr, Port: r.Port}
	// TODO this need to be discuessed
	n = len(b)
	err = syscall.Sendto(r.fd, b, syscall.MSG_FASTOPEN, sa)
	return
}

func (r *razor) Close() error {
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

func (r *razor) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   r.RAddr[:],
		Port: r.RPort,
	}
}
func (r *razor) LocalAddr() net.Addr {
	return &net.TCPAddr{
		IP:   r.Addr[:],
		Port: r.Port,
	}
}
func (r *razor) SetDeadline(t time.Time) error      { return nil }
func (r *razor) SetReadDeadline(t time.Time) error  { return nil }
func (r *razor) SetWriteDeadline(t time.Time) error { return nil }
