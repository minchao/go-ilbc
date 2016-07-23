// based on https://gist.github.com/brendaningram/67c1761b9559cbc6afa3
package ilbc

import (
	"bytes"
	"encoding/binary"
)

// golang types
// uint8       the set of all unsigned  8-bit (1 byte) integers (0 to 255)
// uint16      the set of all unsigned 16-bit (2 byte) integers (0 to 65535)
// uint32      the set of all unsigned 32-bit (4 byte) integers (0 to 4294967295)
// uint64      the set of all unsigned 64-bit (8 byte) integers (0 to 18446744073709551615)

// WAVE PCM soundfile format http://soundfile.sapp.org/doc/WaveFormat/
// The default byte ordering assumed for WAVE data files is little-endian.
// Files written using the big-endian byte ordering scheme have the identifier RIFX instead of RIFF.
type waveHeader struct {
	// Offset 0
	// Contains the letters "RIFF" in ASCII form
	RiffTag [4]byte // ChunkID 4 bytes

	// Offset 4
	// This is the size of the entire file in bytes minus 8 bytes for the two fields not included in this count: ChunkID and ChunkSize.
	// Also = 4 + (8 + SubChunk1Size) + (8 + SubChunk2Size)
	RiffLength uint32 // ChunkSize 4 bytes

	// Offset 8
	// Contains the letters "WAVE"
	WaveTag [4]byte // Format 4 bytes

	// Offset 12
	// Contains the letters "fmt "
	FmtTag [4]byte // Subchunk1ID 4 bytes

	// Offset 16
	// 16 for PCM
	FmtLength uint32 // Subchunk1Size 4 bytes

	// Offset 20
	// PCM = 1 (i.e. Linear quantization)
	// Values other than 1 indicate some form of compression.
	AudioFormat uint16 // AudioFormat 2 bytes

	// Offset 22
	// Mono = 1, Stereo = 2
	NumChannels uint16 // NumChannels 2 bytes

	// Offset 24
	// 44100, 96000 etc
	SampleRate uint32 // SampleRate 4 bytes

	// Offset 28
	// = SampleRate * NumChannels * BitsPerSample/8
	ByteRate uint32 // ByteRate 4 bytes

	// Offset 32
	// The number of bytes for one sample including all channels.
	// = NumChannels * BitsPerSample/8
	BlockAlign uint16 // BlockAlign 2 bytes

	// Offset 34
	// 8 bits = 8, 16 bits = 16
	BitsPerSample uint16 // BitsPerSample 2 bytes

	// Offset 36
	// Contains the letters "data"
	DataTag [4]byte // Subchunk2ID 4 bytes

	// Offset 40
	// This is the number of bytes in the data.
	// = NumSamples * NumChannels * BitsPerSample/8
	DataLength uint32 // Subchunk2Size 4 bytes
}

// Convert the waveHeader struct to a byte slice
func (h *waveHeader) toBytes() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, h)

	return buffer.Bytes()
}

func readWaveHeader(data []byte) (waveHeader, error) {
	var h waveHeader
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &h)
	if err != nil {
		return h, err
	}

	return h, nil
}

func writeWaveHeader(data []byte) []byte {
	var (
		waveHeaderSize uint32 = 44 // bytes
		fileLength        uint32 = uint32(len(data)) + waveHeaderSize
		riffLength        uint32 = fileLength - 8
		dataLength        uint32 = fileLength - waveHeaderSize
	)

	header := &waveHeader{}
	copy(header.RiffTag[:], "RIFF")
	header.RiffLength = riffLength
	copy(header.WaveTag[:], "WAVE")
	copy(header.FmtTag[:], "fmt ")
	header.FmtLength = 16
	header.AudioFormat = 1
	header.NumChannels = 1
	header.SampleRate = 8000
	header.ByteRate = header.SampleRate * uint32(header.NumChannels) * uint32(header.SampleRate/8)
	header.BlockAlign = header.NumChannels * uint16(16/8)
	header.BitsPerSample = uint16(16)
	copy(header.DataTag[:], "data")
	header.DataLength = dataLength

	data = append(header.toBytes(), data...)

	return data
}
