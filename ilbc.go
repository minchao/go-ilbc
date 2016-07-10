package ilbc

/*
#cgo CFLAGS: -I./lib/iLBC_rfc3951

#include <stdlib.h>
#include <math.h>
#include <string.h>
#include "ilbc.h"

void init(int mode);

int encode(unsigned char *samples, int sampleOffset, int sampleLength,
	unsigned char *data, int dataOffset);

int decode(unsigned char *data, int dataOffset, int dataLength,
	unsigned char *samples, int sampleOffset);
*/
import "C"
import (
	"bytes"
	"unsafe"
)

func Encode(samples []byte, mode int) (data []byte) {
	data = make([]byte, len(samples))

	C.init(C.int(mode))
	encoded := C.encode(
		(*C.uchar)(unsafe.Pointer(&samples[0])),
		C.int(0),
		C.int(len(samples)),
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(0),
	)

	return data[:encoded]
}

func Decode(data []byte, mode int) (samples []byte) {
	samples = make([]byte, len(data)*10)

	C.init(C.int(mode))
	decoded := C.decode(
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(0),
		C.int(len(data)),
		(*C.uchar)(unsafe.Pointer(&samples[0])),
		C.int(0),
	)

	return samples[:decoded]
}
