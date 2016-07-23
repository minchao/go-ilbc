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
	"unsafe"
)

var iLBCHeaderPrefix = "#!iLBC"
var iLBCHeaderFormat = iLBCHeaderPrefix + "%d\n"

type frameMode int

const (
	FrameMode20 frameMode = 20
	FrameMode30 frameMode = 30
)

type codec struct {
	mode int
}

func (d *codec) encode(samples []byte) (data []byte) {
	data = make([]byte, len(samples))

	C.init(C.int(d.mode))
	encoded := C.encode(
		(*C.uchar)(unsafe.Pointer(&samples[0])),
		C.int(0),
		C.int(len(samples)),
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(0),
	)

	return data[:encoded]
}

func (d *codec) decode(data []byte) (samples []byte) {
	samples = make([]byte, len(data)*10)

	C.init(C.int(d.mode))
	decoded := C.decode(
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(0),
		C.int(len(data)),
		(*C.uchar)(unsafe.Pointer(&samples[0])),
		C.int(0),
	)

	return samples[:decoded]
}
