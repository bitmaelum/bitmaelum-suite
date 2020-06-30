package core

import (
	"compress/zlib"
	"io"
)

// ZlibCompress compresses a stream through zlib compression
func ZlibCompress(r io.Reader) io.Reader {
	zpr, zpw := io.Pipe()

	writer, err := zlib.NewWriterLevel(zpw, zlib.BestCompression)
	if err != nil {
		return nil
	}

	go func() {
		_, err := io.Copy(writer, r)
		_ = writer.Close()

		if err != nil {
			_ = zpw.CloseWithError(err)
		} else {
			_ = zpw.Close()
		}
	}()

	return zpr
}
