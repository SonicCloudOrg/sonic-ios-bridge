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
	//SSL
	sslConnect net.Conn
	//no SSL
	unencryptedConnect net.Conn
}

type DeviceConnectInterface interface {
	DisableSessionSSL()
	EnableSessionSSL(version []int, pairRecord PairRecord) error
	Reader() io.Reader
	Writer() io.Writer
	GetConn() net.Conn
	Close() error
}

func (conn *DeviceConnection) DisableSessionSSL() {
	if conn.sslConnect != nil {
		conn.sslConnect = nil
	}
	return
}

func (conn *DeviceConnection) EnableSessionSSL(version []int, pairRecord PairRecord) error{
	tlsConn, err := conn.createTLSClient(version,pairRecord)
	if err != nil {
		return err
	}
	conn.unencryptedConnect = conn.sslConnect
	conn.sslConnect = net.Conn(tlsConn)
	return nil
}

func (conn *DeviceConnection) createTLSClient(version []int, pairRecord PairRecord) (*tls.Conn, error) {
	minVersion := uint16(tls.VersionTLS11)
	maxVersion := uint16(tls.VersionTLS11)
	if version[0] > 10 {
		minVersion = tls.VersionTLS11
		maxVersion = tls.VersionTLS13
	}
	cert5, err := tls.X509KeyPair(pairRecord.HostCertificate, pairRecord.HostPrivateKey)
	if err != nil {
		return nil, err
	}
	conf := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert5},
		ClientAuth:         tls.NoClientCert,
		MinVersion:         minVersion,
		MaxVersion:         maxVersion,
	}
	tlsConn := tls.Client(conn.sslConnect, conf)
	err = tlsConn.Handshake()
	if err = tlsConn.Handshake(); err != nil {
		return nil, err
	}
	return tlsConn, nil
}

func (conn *DeviceConnection) Reader() io.Reader {
	return conn.sslConnect
}

func (conn *DeviceConnection) Writer() io.Writer {
	return conn.sslConnect
}

func (conn *DeviceConnection) GetConn() net.Conn {
	if conn.sslConnect != nil {
		return conn.sslConnect
	}
	return conn.unencryptedConnect
}

func (conn *DeviceConnection) Close() error {
	if conn.sslConnect != nil {
		if err := conn.sslConnect.Close(); err != nil {
			return fmt.Errorf("close connect error: %s", err)
		}
	}
	return nil
}

func NewDeviceConnection() (*DeviceConnection, error) {
	conn := &DeviceConnection{}
	return conn, conn.connectSocket()
}

func (conn *DeviceConnection) connectSocket() error {
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
	c, err := d.Dial(network, address)
	if err != nil {
		return fmt.Errorf("error connect socket:%w", err)
	}
	conn.sslConnect = c
	return nil
}
