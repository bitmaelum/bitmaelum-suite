package core

import (
	"compress/zlib"
	"fmt"
	"io"
)

func Compress(r io.Reader) io.Reader {
	zpr, zpw := io.Pipe()

	//var rbuf = make([]byte, 32 * 1024)

	writer, err := zlib.NewWriterLevel(zpw, zlib.BestCompression)
	if err != nil {
		return nil
	}

	go func() {
		fmt.Printf("Copying r to writer\n")
		_, err := io.Copy(writer, r)
		fmt.Printf("Closing writer\n")
		_ = writer.Close()

		if err != nil {
			fmt.Printf("Closing zpw with err " + err.Error() + "\n")
			_ = zpw.CloseWithError(err)
		} else {
			fmt.Printf("Closing zpw\n")
			_ = zpw.Close()
		}
	}()

	return zpr
}
