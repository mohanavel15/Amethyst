package zlib

import (
	"bytes"

	"compress/zlib"
)

func Decode(buf *bytes.Buffer, outbuf []byte) error {
	r, err := zlib.NewReader(buf)
	if err != nil {
		return err
	}

	_, err = r.Read(outbuf)
	if err != nil {
		return err
	}

	_ = r.Close()
	return nil
}

func Encode(buffer []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write(buffer)
	return buf.Bytes(), err
}
