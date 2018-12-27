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
	CFB BlockCipherMode = iota
	CTR
	OFB
)

const bufInitCap = 256

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
	case CFB:
		stream = cipher.NewCFBEncrypter(block, iv)
	case CTR:
		stream = cipher.NewCTR(block, iv)
	case OFB:
		stream = cipher.NewOFB(block, iv)
	default:
		return wt, errors.New("unhandled block cipher mode")
	}
	return &streamEncrypter{
		underlying: wt,
		stream:     stream,
		iv:         iv,
		buf:        make([]byte, 0, bufInitCap),
	}, nil
}

type streamEncrypter struct {
	underlying io.WriteCloser
	stream     cipher.Stream
	iv         []byte
	buf        []byte
}

func (enc *streamEncrypter) Close() error {
	return enc.underlying.Close()
}

func (enc *streamEncrypter) Write(bs []byte) (int, error) {
	var count int
	if len(enc.iv) > 0 {
		n, err := enc.underlying.Write(enc.iv)
		count += n
		enc.iv = enc.iv[n:]
		if err != nil {
			return count, err
		}
	}

	size := len(bs)
	for cap(enc.buf) < size {
		enc.buf = make([]byte, 0, cap(enc.buf)<<1)
	}
	buf := enc.buf[:size]
	enc.stream.XORKeyStream(buf, bs)
	n, err := enc.underlying.Write(buf)
	count += n

	return count, err
}
