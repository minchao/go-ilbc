#include <math.h>
#include <string.h>
#include "iLBC_define.h"
#include "iLBC_decode.h"
#include "iLBC_encode.h"

static iLBC_Enc_Inst_t g_enc_inst;
static iLBC_Dec_Inst_t g_dec_inst;

void init(int mode) {
    initEncode(&g_enc_inst, mode);
    initDecode(&g_dec_inst, mode, 1);
}

static int _encode(short *samples, unsigned char *data) {
    int i;
	float block[BLOCKL_MAX];

	// Convert to float representaion of voice signal.
    for (i = 0; i < g_enc_inst.blockl; i++) {
        block[i] = samples[i];
    }

    iLBC_encode(data, block, &g_enc_inst);

    return g_enc_inst.no_of_bytes;
}

static int _decode(unsigned char *data, short *samples, int mode) {
    int i;
    float block[BLOCKL_MAX];

    // Validate Mode
    if (mode != 0 && mode != 1) {
        return -1;
    }

    iLBC_decode(block, data, &g_dec_inst, mode);

    // Validate PCM16
    for (i = 0; i < g_dec_inst.blockl; i++) {
        float point;

        point = block[i];
        if (point < MIN_SAMPLE) {
            point = MIN_SAMPLE;
        } else if (point > MAX_SAMPLE) {
            point = MAX_SAMPLE;
        }

        samples[i] = point;
    }

    return g_dec_inst.blockl * 2;
}

int encode(unsigned char *samples, int sampleOffset, int sampleLength,
    unsigned char *data, int dataOffset) {

    int bytes_to_encode;
    int bytes_encoded;
    int truncated;

    samples += sampleOffset;
    data += dataOffset;

    bytes_to_encode = sampleLength;
    bytes_encoded = 0;

    truncated = bytes_to_encode % (g_enc_inst.blockl * 2);
    if(!truncated) {
        bytes_to_encode -= truncated;
	}

	while(bytes_to_encode > 0) {
		int _encoded = _encode((short *)samples, data);
		samples += g_enc_inst.blockl * 2;
		data += _encoded;
		bytes_encoded += _encoded;
		bytes_to_encode -= g_enc_inst.blockl * 2;
	}

	samples -= sampleLength;
	data -= bytes_encoded;

	return bytes_encoded;
}

int decode(unsigned char *data, int dataOffset, int dataLength,
    unsigned char *samples, int sampleOffset) {

	int bytes_to_decode;
	int bytes_decoded;

	samples += sampleOffset;
	data += dataOffset;

	bytes_to_decode = dataLength;
	bytes_decoded = 0;

    while (bytes_to_decode > 0) {
        int _decoded = _decode(data, (short *)samples, 1);
        samples += _decoded;
        data += g_dec_inst.no_of_bytes;
        bytes_decoded += _decoded;
        bytes_to_decode -= g_dec_inst.no_of_bytes;
    }

    // Revert buffer pointers
    samples -= bytes_decoded;
    data -= dataLength;

    return bytes_decoded;
}
