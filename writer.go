package ilbc

import (
	"fmt"
	"io"
)

func writeILBCHeader(data []byte, mode frameMode) []byte {
	return append([]byte(fmt.Sprintf(iLBCHeaderFormat, mode)), data...)
}

type Options struct {
	Mode frameMode
}

// Encode writes the WAVE []byte p to w in iLBC format with the given options.
// Default parameters are used if a nil *Options is passed.
func Encode(w io.Writer, p []byte, o *Options) error {
	// The input to the encoder SHOULD be 16 bit uniform PCM sampled at 8 kHz.
	// specified in section 2.1.
	h, err := readWaveHeader(p)
	if err != nil {
		return fmt.Errorf("error decoding header, err: %s", err)
	}
	if h.NumChannels != 1 {
		return fmt.Errorf("invalid channel %d, should be mono", h.NumChannels)
	}
	if h.SampleRate != 8000 {
		return fmt.Errorf("invalid sample rate %d, should be 16", h.SampleRate)
	}
	if h.BitsPerSample != 16 {
		return fmt.Errorf("invalid bits per sample %d, should be 8000", h.BitsPerSample)
	}

	mode := FrameMode20
	if o != nil {
		mode = o.Mode
	}
	e := codec{mode: int(mode)}
	data := e.encode(p[44:])

	// Add iLBC header
	data = writeILBCHeader(data, mode)

	_, err = w.Write(data)
	return err
}
