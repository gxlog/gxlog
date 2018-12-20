package file

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

type BlockCipherMode int

const (
	ModeCFB BlockCipherMode = iota
	ModeCTR
	ModeOFB
)

const cBufInitCap = 256

func newAESWriter(wt io.WriteCloser, key string, mode BlockCipherMode) (io.WriteCloser, error) {
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return wt, err
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return wt, err
	}
	iv := make([]byte, aes.BlockSize) // initialization vector
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return wt, err
	}
	return newStreamEncrypter(wt, block, iv, mode)
}

func newStreamEncrypter(wt io.WriteCloser, block cipher.Block, iv []byte,
	mode BlockCipherMode) (io.WriteCloser, error) {
	var stream cipher.Stream
	switch mode {
	case ModeCFB:
		stream = cipher.NewCFBEncrypter(block, iv)
	case ModeCTR:
		stream = cipher.NewCTR(block, iv)
	case ModeOFB:
		stream = cipher.NewOFB(block, iv)
	default:
		return wt, errors.New("unhandled block cipher mode")
	}
	return &streamEncrypter{
		underlying: wt,
		stream:     stream,
		iv:         iv,
		buf:        make([]byte, 0, cBufInitCap),
	}, nil
}

type streamEncrypter struct {
	underlying io.WriteCloser
	stream     cipher.Stream
	iv         []byte
	buf        []byte
}

func (this *streamEncrypter) Close() error {
	return this.underlying.Close()
}

func (this *streamEncrypter) Write(bs []byte) (int, error) {
	var count int
	if len(this.iv) > 0 {
		n, err := this.underlying.Write(this.iv)
		count += n
		this.iv = this.iv[n:]
		if err != nil {
			return count, err
		}
	}

	size := len(bs)
	for cap(this.buf) < size {
		this.buf = make([]byte, 0, cap(this.buf)<<1)
	}
	buf := this.buf[:size]
	this.stream.XORKeyStream(buf, bs)
	n, err := this.underlying.Write(buf)
	count += n

	return count, err
}
