package sphinx

import (
	"encoding/binary"
	"errors"
	"io"
)

// ErrVarIntNotCanonical indicates the decoded varint was not minimally encoded.
var ErrVarIntNotCanonical = errors.New("decoded varint is not canonical")

// ReadVarInt reads a variable-length integer from `r` and returns it as a uint64.
// It uses a buffer `buf` to read data efficiently.
func ReadVarInt(r io.Reader, buf *[8]byte) (uint64, error) {
	// Read the first byte to determine the discriminant.
	_, err := io.ReadFull(r, buf[:1])
	if err != nil {
		return 0, err
	}
	discriminant := buf[0]

	var value uint64
	switch {
	case discriminant < 0xfd:
		value = uint64(discriminant)

	case discriminant == 0xfd:
		// Read the next 2 bytes for uint16.
		_, err := io.ReadFull(r, buf[:2])
		if err == io.EOF {
			return 0, io.ErrUnexpectedEOF
		} else if err != nil {
			return 0, err
		}
		value = uint64(binary.BigEndian.Uint16(buf[:2]))

		// Ensure canonical encoding.
		if value < 0xfd {
			return 0, ErrVarIntNotCanonical
		}

	case discriminant == 0xfe:
		// Read the next 4 bytes for uint32.
		_, err := io.ReadFull(r, buf[:4])
		if err == io.EOF {
			return 0, io.ErrUnexpectedEOF
		} else if err != nil {
			return 0, err
		}
		value = uint64(binary.BigEndian.Uint32(buf[:4]))

		// Ensure canonical encoding.
		if value <= 0xffff {
			return 0, ErrVarIntNotCanonical
		}

	default:
		// Read the next 8 bytes for uint64.
		_, err := io.ReadFull(r, buf[:8])
		if err == io.EOF {
			return 0, io.ErrUnexpectedEOF
		} else if err != nil {
			return 0, err
		}
		value = binary.BigEndian.Uint64(buf[:8])

		// Ensure canonical encoding.
		if value <= 0xffffffff {
			return 0, ErrVarIntNotCanonical
		}
	}

	return value, nil
}

// WriteVarInt writes a uint64 `val` to `w` using a variable number of bytes.
// It uses the buffer `buf` for efficient serialization.
func WriteVarInt(w io.Writer, val uint64, buf *[8]byte) error {
	var length int

	switch {
	case val < 0xfd:
		buf[0] = uint8(val)
		length = 1

	case val <= 0xffff:
		buf[0] = 0xfd
		binary.BigEndian.PutUint16(buf[1:3], uint16(val))
		length = 3

	case val <= 0xffffffff:
		buf[0] = 0xfe
		binary.BigEndian.PutUint32(buf[1:5], uint32(val))
		length = 5

	default:
		buf[0] = 0xff
		binary.BigEndian.PutUint64(buf[1:9], val)
		length = 9
	}

	_, err := w.Write(buf[:length])
	return err
}
