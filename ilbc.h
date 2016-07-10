#ifndef ILBC_H
#define ILBC_H

void init(int mode);

int encode(unsigned char *samples, int sampleOffset, int sampleLength,
	unsigned char *data, int dataOffset);

int decode(unsigned char *data, int dataOffset, int dataLength,
	unsigned char *samples, int sampleOffset);

#endif