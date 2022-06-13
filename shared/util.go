package shared

import (
	"crypto/tls"
	"errors"
	"flag"
	"os"
	"os/exec"
	"strings"
	_ "unsafe"
)

var tlsMinVersion = flag.String("tls-min-version", "1.2", "Minimum TLS version")
var tlsMaxVersion = flag.String("tls-max-version", "1.3", "Maximum TLS version")
var tlsCipherPreference = flag.String("tls-cipher-preference", "", "Prefer AES ciphers (AES), or ChaCha ciphers (CHACHA), don't specify for default behaviour")

func ExecCmd(cmd string, arg ...string) error {
	cmdO := exec.Command(cmd, arg...)
	cmdO.Stdout = os.Stdout
	cmdO.Stderr = os.Stderr
	return cmdO.Run()
}

type MacAddr [6]byte

func GetSrcMAC(packet []byte) MacAddr {
	var mac MacAddr
	copy(mac[:], packet[6:12])
	return mac
}

func GetDestMAC(packet []byte) MacAddr {
	var mac MacAddr
	copy(mac[:], packet[0:6])
	return mac
}

func MACIsUnicast(mac MacAddr) bool {
	return (mac[0] & 1) == 0
}

func TlsVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "1.0"
	case tls.VersionTLS11:
		return "1.1"
	case tls.VersionTLS12:
		return "1.2"
	case tls.VersionTLS13:
		return "1.3"
	}
	return "Invalid"
}

func TlsVersionNum(version string) uint16 {
	switch version {
	case "1.0":
		return tls.VersionTLS10
	case "1.1":
		return tls.VersionTLS11
	case "1.2":
		return tls.VersionTLS12
	case "1.3":
		return tls.VersionTLS13
	}
	return 0
}

func TlsUseFlags(tlsConfig *tls.Config) {
	tlsConfig.MinVersion = TlsVersionNum(*tlsMinVersion)
	tlsConfig.MaxVersion = TlsVersionNum(*tlsMaxVersion)

	switch strings.ToUpper(*tlsCipherPreference) {
	case "AES":
		TlsSetCipherAESPreference(true)
	case "CHACHA":
		TlsSetCipherAESPreference(false)
	case "":
		break
	default:
		panic(errors.New("invalid TLS preference. Must be blank, AES or CHACHA"))
	}
}

//go:linkname hasAESGCMHardwareSupport crypto/tls.hasAESGCMHardwareSupport
var hasAESGCMHardwareSupport bool

func TlsSetCipherAESPreference(preferAES bool) {
	hasAESGCMHardwareSupport = preferAES
}
