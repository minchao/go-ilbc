package ilbc

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strconv"
	"fmt"
)

// Decode reads a iLBC from r and returns it as WAVE []byte.
func Decode(r io.Reader) ([]byte, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if len(buf) < 9 || !bytes.HasPrefix(buf, []byte(iLBCHeaderPrefix)) {
		return nil, errors.New("iLBC header not found")
	}
	// Get iLBC frame mode
	header := buf[:9]
	m, _ := strconv.ParseInt(string(header[6:8]), 10, 64)
	mode := frameMode(m)
	if !(mode == FrameMode20 || mode == FrameMode30) {
		return nil, fmt.Errorf("invalid frame mode %d", mode)
	}

	d := codec{mode: int(mode)}
	data := d.decode(buf[9:])

	// Add WAVE header
	data = writeWaveHeader(data)

	return data, nil
}
