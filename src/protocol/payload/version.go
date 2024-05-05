package payload

import (
	"bytes"
	"io"

	"github.com/tomaszkoziara/btc-handshake/src/utils"
)

type Version struct {
	version     int32
	services    uint64
	timestamp   int64
	addrRecv    NetworkAddr
	addrFrom    NetworkAddr
	nonce       uint64
	userAgent   string
	startHeight int32
	relay       bool
}

func (v *Version) WithServices(services uint64) *Version {
	v.services = services
	return v
}

func (v *Version) WithTimestamp(timestamp int64) *Version {
	v.timestamp = timestamp
	return v
}

func (v *Version) WithAddrRecv(addrRecv NetworkAddr) *Version {
	v.addrRecv = addrRecv
	return v
}

func (v *Version) WithAddrFrom(addrFrom NetworkAddr) *Version {
	v.addrFrom = addrFrom
	return v
}

func (v *Version) WithNonce(nonce uint64) *Version {
	v.nonce = nonce
	return v
}

func (v *Version) WithUserAgent(userAgent string) *Version {
	v.userAgent = userAgent
	return v
}

func (v *Version) WithStartHeight(startHeight int32) *Version {
	v.startHeight = startHeight
	return v
}

func (v *Version) WithRelay(relay bool) *Version {
	v.relay = relay
	return v
}

func (v *Version) Decode(reader io.Reader) error {
	panic("not implemented")
}

func (v *Version) Encode(writer io.Writer) error {
	bw := utils.NewBinaryWriter(new(bytes.Buffer))

	bw.WriteLittleEndian(v.version)
	bw.WriteLittleEndian(v.services)
	bw.WriteLittleEndian(v.timestamp)

	addrRecv, err := v.addrRecv.Bytes()
	if err != nil {
		return err
	}
	bw.WriteBigEndian(addrRecv)

	if v.version < 106 {
		if bw.Err != nil {
			return bw.Err
		}

		if _, err := writer.Write(bw.Buffer.Bytes()); err != nil {
			return err
		}

		return nil
	}

	addrFrom, err := v.addrFrom.Bytes()
	if err != nil {
		return err
	}
	bw.WriteBigEndian(addrFrom)
	bw.WriteLittleEndian(v.nonce)
	bw.WriteString(v.userAgent)
	bw.WriteLittleEndian(v.startHeight)

	if v.version < 70001 {
		if bw.Err != nil {
			return bw.Err
		}

		if _, err := writer.Write(bw.Buffer.Bytes()); err != nil {
			return err
		}

		return nil
	}

	relay := uint8(0)
	if v.relay {
		relay = 1
	}
	bw.WriteLittleEndian(relay)

	if bw.Err != nil {
		return bw.Err
	}

	if _, err := writer.Write(bw.Buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

func NewVersion(version int32) *Version {
	return &Version{
		version: version,
	}
}
