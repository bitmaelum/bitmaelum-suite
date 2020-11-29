// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package message

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

// ZlibDecompress will return a reader that automatically decompresses the stream
func ZlibDecompress(r io.Reader) (io.Reader, error) {
	return zlib.NewReader(r)
}
