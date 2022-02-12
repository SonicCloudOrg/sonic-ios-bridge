package conn

type PairRecord struct {
	HostID            string
	SystemBUID        string
	HostCertificate   []byte
	HostPrivateKey    []byte
	DeviceCertificate []byte
	EscrowBag         []byte
	WiFiMACAddress    string
	RootCertificate   []byte
	RootPrivateKey    []byte
}
