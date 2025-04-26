package format

import (
	"encoding/binary"
	"errors"
	"io"
)

func checkType(type_ byte) bool {
	switch type_ {
	case CodeType, DocsType, TreeType:
		return true
	default:
		return false
	}
}

func FormHeader(type_ byte) ([]byte, error) {
	if !checkType(type_) {
		return nil, errors.New("unknown header type")
	}

	header := make([]byte, HeaderSize)
	header[0] = type_

	return header, nil
}

func PutSizeInHeader(header []byte, size int) error {
	if len(header) != HeaderSize {
		return errors.New("invalid header size")
	}

	binary.BigEndian.PutUint32(header[typeBytesN:], uint32(size))
	return nil
}

func parseHeader(header []byte) (byte, int, error) {
	if len(header) != HeaderSize {
		return 0, 0, errors.New("invalid header")
	}

	type_ := header[0]
	if !checkType(type_) {
		return 0, 0, errors.New("unknown header type")
	}
	size := int(binary.BigEndian.Uint32(header[typeBytesN:]))

	return type_, size, nil
}

func ParseHeaderFromData(data []byte) (byte, int, error) {
	header := data[:HeaderSize]
	return parseHeader(header)
}

func ParseHeaderFromFile(file io.Reader) (byte, int, error) {
	header := make([]byte, HeaderSize)
	n, err := file.Read(header)
	if err != nil {
		return 0, 0, err
	}
	if n != HeaderSize {
		return 0, 0, errors.New("invalid header")
	}

	return parseHeader(header)
}
