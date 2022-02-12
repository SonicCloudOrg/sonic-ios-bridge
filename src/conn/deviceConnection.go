package conn

import (
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

type DeviceConnectInterface interface {
	DisableSessionSSL()
	Send(message []byte) error
	Reader() io.Reader
	Writer() io.Writer
	Connect() net.Conn
	Close() error
}

type DeviceConnection struct {
	encryptedConnect   net.Conn
	unencryptedConnect net.Conn
}

func NewDeviceConnection(socketToConnectTo string) (*DeviceConnection, error) {
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
