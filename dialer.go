package razor

import (
	"context"
	"net"
	"syscall"
	"time"
)

type razor struct {
	ServerAddr      [4]byte
	ServerPort      int
	RAddr           [4]byte
	RPort           int
	fd              int
	ctx, rctx, wctx context.Context
}

func (r *razor) Read(b []byte) (n int, err error) {
	n, err = syscall.Read(r.fd, b)
	return

}

func (r *razor) Write(b []byte) (n int, err error) {
	sa := &syscall.SockaddrInet4{Addr: r.ServerAddr, Port: r.ServerPort}
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

func (r *razor) LocalAddr() net.Addr
func (r *razor) RemoteAddr() net.Addr
func (r *razor) SetDeadline(t time.Time) error
func (r *razor) SetReadDeadline(t time.Time) error
func (r *razor) SetWriteDeadline(t time.Time) error
