package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"strings"
)

func StringToBytes(s string, maxBytes int) []byte {
	b := []byte(s)

	if len(b) > maxBytes {
		b = b[:maxBytes]
	}
	if len(b) < maxBytes {
		padding := make([]byte, maxBytes-len(b))
		b = append(b, padding...)
	}

	return b
}

func Checksum(payload []byte) []byte {
	hash := sha256.Sum256(payload)
	hash = sha256.Sum256(hash[:])

	return hash[:4]
}

type BinaryWriter struct {
	Buffer *bytes.Buffer
	Err    error
}

func (bw *BinaryWriter) WriteLittleEndian(data interface{}) {
	if bw.Err != nil {
		return
	}

	bw.Err = binary.Write(bw.Buffer, binary.LittleEndian, data)
}

func (bw *BinaryWriter) WriteBigEndian(data interface{}) {
	if bw.Err != nil {
		return
	}

	bw.Err = binary.Write(bw.Buffer, binary.BigEndian, data)
}

func (bw *BinaryWriter) WriteString(data string) {
	if bw.Err != nil {
		return
	}

	_, err := bw.Buffer.WriteString(data)
	bw.Err = err
}

func (bw *BinaryWriter) Write(data []byte) {
	if bw.Err != nil {
		return
	}

	_, err := bw.Buffer.Write(data)
	bw.Err = err
}

func NewBinaryWriter(buffer *bytes.Buffer) *BinaryWriter {
	return &BinaryWriter{
		Buffer: buffer,
	}
}

type BinaryReader struct {
	Reader io.Reader
	Err    error
}

func (br *BinaryReader) ReadUint32LittleEndian() uint32 {
	if br.Err != nil {
		return uint32(0)
	}

	b, err := br.readBytes(4)
	if err != nil {
		br.Err = err
	}

	return binary.LittleEndian.Uint32(b)
}

func (br *BinaryReader) ReadUint64LittleEndian() uint64 {
	if br.Err != nil {
		return uint64(0)
	}

	b, err := br.readBytes(8)
	if err != nil {
		br.Err = err
	}

	return binary.LittleEndian.Uint64(b)
}

func (br *BinaryReader) ReadUint16BigEndian() uint16 {
	if br.Err != nil {
		return uint16(0)
	}

	b, err := br.readBytes(2)
	if err != nil {
		br.Err = err
	}

	return binary.BigEndian.Uint16(b)
}

func (br *BinaryReader) ReadString(bytes int) string {
	if br.Err != nil {
		return ""
	}

	b, err := br.readBytes(bytes)
	if err != nil {
		br.Err = err
	}

	return strings.Trim(string(b), "\x00")
}

func (br *BinaryReader) ReadBytes(bytes int) []byte {
	if br.Err != nil {
		return []byte{}
	}

	b, err := br.readBytes(bytes)
	if err != nil {
		br.Err = err
	}

	return b
}

func (br *BinaryReader) readBytes(bytes int) ([]byte, error) {
	b := make([]byte, bytes)
	_, err := br.Reader.Read(b)
	br.Err = err

	return b, nil
}

func NewBinaryReader(reader io.Reader) *BinaryReader {
	return &BinaryReader{
		Reader: reader,
	}
}
