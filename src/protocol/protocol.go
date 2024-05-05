package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/tomaszkoziara/btc-handshake/src/utils"
)

var (
	ErrInvalidChecksum error = errors.New("invalid sumcheck")
)

type Message struct {
	Header  Header
	Payload Payload
}

type Payload interface {
	Encode(writer io.Writer) error
	Decode(reader io.Reader) error
}

type Header struct {
	Magic    uint32
	Command  string
	Length   uint32
	Checksum uint32
}

func (h *Header) WithMagic(magic uint32) *Header {
	h.Magic = magic
	return h
}

func (h *Header) WithCommand(command string) *Header {
	h.Command = command
	return h
}

func (h *Header) WithLength(length uint32) *Header {
	h.Length = length
	return h
}

func (h *Header) WithChecksum(checksum uint32) *Header {
	h.Checksum = checksum
	return h
}

func (h *Header) Decode(reader io.Reader) error {
	// br := utils.NewBinaryReader(reader)

	// v.Version = version
	// v.VersionMsg = versionMsg

	// if version >= 31402 && versionMsg {
	// 	v.Time = br.ReadUint32LittleEndian()
	// }

	// v.Services = br.ReadUint64LittleEndian()
	// v.IP = bytesToIP(br.ReadBytes(16))
	br := utils.NewBinaryReader(reader)

	h.Magic = br.ReadUint32LittleEndian()
	h.Command = br.ReadString(12)
	h.Length = br.ReadUint32LittleEndian()
	h.Checksum = br.ReadUint32LittleEndian()

	if br.Err != nil {
		return br.Err
	}

	return nil
	// b, err := read(reader, 4)
	// if err != nil {
	// 	return nil, err
	// }
	// magic := binary.LittleEndian.Uint32(b)

	// b, err = read(reader, 12)
	// if err != nil {
	// 	return nil, err
	// }
	// command := strings.Trim(string(b), "\x00")

	// b, err = read(reader, 4)
	// if err != nil {
	// 	return nil, err
	// }
	// length := binary.LittleEndian.Uint32(b)

	// b, err = read(reader, 4)
	// if err != nil {
	// 	return nil, err
	// }
	// checksum := binary.LittleEndian.Uint32(b)

	// payload := []byte{}
	// if length > 0 {
	// 	payload, err = read(reader, int(length))
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	computedChecksum := binary.LittleEndian.Uint32(Checksum(payload))
	// 	if computedChecksum != checksum {
	// 		return nil, ErrInvalidChecksum
	// 	}
	// }

	// header := Header{
	// 	Magic:    magic,
	// 	Command:  command,
	// 	Length:   length,
	// 	Checksum: checksum,
	// }

	// return &header, nil
}

func (h *Header) Encode(writer io.Writer) error {
	bw := utils.NewBinaryWriter(new(bytes.Buffer))

	bw.WriteLittleEndian(h.Magic)
	bw.Write(utils.StringToBytes(h.Command, 12))
	bw.WriteLittleEndian(h.Length)
	bw.WriteLittleEndian(h.Checksum)

	if bw.Err != nil {
		return bw.Err
	}

	if _, err := writer.Write(bw.Buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (msg *Message) Encode(writer io.Writer) error {
	payload := []byte{}
	length := uint32(0)
	checksum := binary.LittleEndian.Uint32(utils.Checksum(payload))

	if msg.Payload != nil {
		payloadBuf := new(bytes.Buffer)
		if err := msg.Payload.Encode(payloadBuf); err != nil {
			return err
		}
		payload = payloadBuf.Bytes()
		length = uint32(payloadBuf.Len())
		checksum = binary.LittleEndian.Uint32(utils.Checksum(payload))
	}

	msg.Header.Length = length
	msg.Header.Checksum = checksum

	if err := msg.Header.Encode(writer); err != nil {
		return err
	}

	if _, err := writer.Write(payload); err != nil {
		return err
	}

	return nil
}

// func (h *Header) EncodeHeader(writer io.Writer) error {
// 	payload := []byte{}
// 	if h.Payload != nil {
// 		payload = h.Payload
// 	}

// 	binary.Write(writer, binary.LittleEndian, msg.Magic)
// 	binary.Write(writer, binary.BigEndian, StringToBytes(msg.Command, 12))
// 	binary.Write(writer, binary.LittleEndian, uint32(len(payload)))

// 	payloadHash := sha256.Sum256(payload)
// 	payloadHash = sha256.Sum256(payloadHash[:])
// 	binary.Write(writer, binary.LittleEndian, payloadHash[0:4])
// 	binary.Write(writer, binary.BigEndian, payload)

// 	return message.Bytes()
// }

// func read(reader io.Reader, bytes int) ([]byte, error) {
// 	b := make([]byte, bytes)
// 	if _, err := reader.Read(b); err != nil {
// 		return nil, err
// 	}

// 	return b, nil
// }
