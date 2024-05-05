package payload_test

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tomaszkoziara/btc-handshake/src/protocol/payload"
)

func TestICanSerializeAddressCorrectly(t *testing.T) {
	now := time.Now().Unix()
	ipv4 := net.ParseIP("192.168.1.1")

	testCases := []struct {
		name        string
		networkaddr payload.NetworkAddr
		version     int32
		versionMsg  bool
	}{
		{
			name:       "IPv4",
			version:    99999,
			versionMsg: true,
			networkaddr: func() payload.NetworkAddr {
				na := payload.NewNetworkAddr(99999, true)
				na.WithTime(uint32(now)).
					WithServices(1000).
					WithIP(ipv4).
					WithPort(3000)
				return *na
			}(),
		},
		// TODO: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encodedNetworkAddr, err := tc.networkaddr.Bytes()
			assert.NoError(t, err)

			decodedNetworkAddr := new(payload.NetworkAddr)
			err = decodedNetworkAddr.Decode(bytes.NewReader(encodedNetworkAddr), tc.version, tc.versionMsg)
			assert.NoError(t, err)

			assert.Equal(t, tc.networkaddr.Time, decodedNetworkAddr.Time)
			assert.Equal(t, tc.networkaddr.Services, decodedNetworkAddr.Services)
			assert.Equal(t, tc.networkaddr.IP.To4(), decodedNetworkAddr.IP.To4())
			assert.Equal(t, tc.networkaddr.IP.To16(), decodedNetworkAddr.IP.To16())
			assert.Equal(t, tc.networkaddr.Port, decodedNetworkAddr.Port)
		})
	}
}

// TODO: add more tests for errors and unhappy path
