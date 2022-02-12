package conn

import (
	"crypto/tls"
	"fmt"
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
	HandShake(version []int, pairRecord PairRecord) (err error)
	Reader(length int) (data []byte, err error)
	Writer(data []byte) (err error)
	GetConn() net.Conn
	Close() error
}

func (conn *DeviceConnection) DisableSessionSSL() {
	if conn.sslConnect != nil {
		conn.sslConnect = nil
	}
	return
}

func (conn *DeviceConnection) HandShake(version []int, pairRecord PairRecord) (err error) {
	minVersion := uint16(tls.VersionTLS11)
	maxVersion := uint16(tls.VersionTLS11)
	if version[0] > 10 {
		minVersion = tls.VersionTLS11
		maxVersion = tls.VersionTLS13
	}
	cert5, err := tls.X509KeyPair(pairRecord.HostCertificate, pairRecord.HostPrivateKey)
	if err != nil {
		return err
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
		return err
	}
	return
}

func (conn *DeviceConnection) Reader(length int) (data []byte, err error) {
	c := conn.GetConn()
	err = c.SetReadDeadline(time.Now().Add(DeviceConnectTimeout))
	if err != nil {
		return nil, err
	}
	data = make([]byte, 0, length)
	for len(data) < length {
		buf := make([]byte, length-len(data))
		chunk, err1 := 0, error(nil)
		if chunk, err1 = c.Read(buf); err1 != nil && chunk == 0 {
			return nil, err1
		}
		//append chunks
		data = append(data, buf[:chunk]...)
	}
	return
}

func (conn *DeviceConnection) Writer(data []byte) (err error) {
	c := conn.GetConn()
	err = c.SetWriteDeadline(time.Now().Add(DeviceConnectTimeout))
	if err != nil {
		return err
	}
	for length := 0; length < len(data); {
		var s int
		if s, err = c.Write(data[length:]); err != nil {
			return err
		}
		if s == 0 {
			return err
		}
		//computed
		length += s
	}
	return
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
	if conn.unencryptedConnect != nil {
		if err := conn.unencryptedConnect.Close(); err != nil {
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
