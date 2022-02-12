package conn

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"runtime"
	"time"
)

const (
	DeviceConnectTimeout = 1 * time.Minute
	TcpSocketAddress     = "127.0.0.1:27015"
	UnixSocketAddress    = "/var/run/usbmuxd"
)

type DeviceConnection struct {
	encryptedConnect   net.Conn
	unencryptedConnect net.Conn
}

func NewDeviceConnection() (*DeviceConnection, error) {
	conn := &DeviceConnection{}
	return conn, conn.connect()
}

func (conn *DeviceConnection) connect() (err error) {
	var network, address string
	switch runtime.GOOS {
	case "windows":
		network, address = "tcp", TcpSocketAddress
	case "darwin", "android", "linux":
		network, address = "unix", UnixSocketAddress
	default:
		return fmt.Errorf("unsupported system: %s, please report to https://github.com/SonicCloudOrg/sonic-ios-bridge",
			runtime.GOOS)
	}
	d := net.Dialer{
		Timeout: DeviceConnectTimeout,
	}
	conn.encryptedConnect, err = d.Dial(network, address)
	if err != nil {
		return fmt.Errorf("fail to connect socket: %w", err)
	}
	return
}

type DeviceConnectInterface interface {
	DisableSessionSSL()
	Send(message []byte) error
	Reader() io.Reader
	Writer() io.Writer
	Connect() net.Conn
	Close() error
}

type sslConnection struct {
	conn    net.Conn
	sslConn *tls.Conn
	timeout time.Duration
}

func (conn *DeviceConnection) DisableSessionSSL() {
	panic("implement me")
}

func (conn *DeviceConnection) Send(message []byte) error {
	panic("implement me")
}

func (conn *DeviceConnection) Reader() io.Reader {
	panic("implement me")
}

func (conn *DeviceConnection) Writer() io.Writer {
	panic("implement me")
}

func (conn *DeviceConnection) Connect() net.Conn {
	panic("implement me")
}

func (conn *DeviceConnection) Close() error {
	panic("implement me")
}
