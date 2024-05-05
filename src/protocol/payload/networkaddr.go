package payload

import (
	"bytes"
	"io"
	"net"

	"github.com/tomaszkoziara/btc-handshake/src/utils"
)

type NetworkAddr struct {
	Version    int32
	VersionMsg bool

	// net_addr properties
	Time     uint32
	Services uint64
	IP       net.IP
	Port     uint16
}

func (na *NetworkAddr) WithTime(time uint32) *NetworkAddr {
	na.Time = time
	return na
}

func (na *NetworkAddr) WithServices(services uint64) *NetworkAddr {
	na.Services = services
	return na
}

func (na *NetworkAddr) WithIP(ip net.IP) *NetworkAddr {
	na.IP = ip
	return na
}

func (na *NetworkAddr) WithPort(port uint16) *NetworkAddr {
	na.Port = port
	return na
}

func (na *NetworkAddr) Bytes() ([]byte, error) {
	bw := utils.NewBinaryWriter(new(bytes.Buffer))

	if na.Version >= 31402 && na.VersionMsg {
		bw.WriteLittleEndian(na.Time)
	}
	bw.WriteLittleEndian(na.Services)
	bw.WriteBigEndian(ipTo16Bytes(na.IP.To16()))
	bw.WriteBigEndian(na.Port)

	if bw.Err != nil {
		return nil, bw.Err
	}

	return bw.Buffer.Bytes(), nil
}

func (v *NetworkAddr) Decode(reader io.Reader, version int32, versionMsg bool) error {
	br := utils.NewBinaryReader(reader)

	v.Version = version
	v.VersionMsg = versionMsg

	if version >= 31402 && versionMsg {
		v.Time = br.ReadUint32LittleEndian()
	}

	v.Services = br.ReadUint64LittleEndian()
	v.IP = bytesToIP(br.ReadBytes(16))
	v.Port = br.ReadUint16BigEndian()

	if br.Err != nil {
		return br.Err
	}

	return nil
}

func ipTo16Bytes(ip net.IP) [16]byte {
	ip16 := ip.To16()
	if ip16 == nil {
		return [16]byte{}
	}
	var ipArray [16]byte
	copy(ipArray[:], ip16)
	return ipArray
}

func bytesToIP(ip []byte) net.IP {
	return net.IP(ip[:])
}

func NewNetworkAddr(version int32, versionMsg bool) *NetworkAddr {
	return &NetworkAddr{
		Version:    version,
		VersionMsg: versionMsg,
	}
}
