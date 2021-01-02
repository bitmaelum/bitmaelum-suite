// Copyright (c) 2021 BitMaelum Authors
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
	"bufio"
	"crypto/sha256"
	"encoding/hex"

	"io"
	"os"

	//lint:ignore SA1019 we use RIPE only for checksums, not encryption or hashing and is always used together with other checksum hashes
	"golang.org/x/crypto/ripemd160"
)

// CalculateChecksums calculates a number of hashes for the given reader in one go.
// Taken from http://marcio.io/2015/07/calculating-multiple-file-hashes-in-a-single-pass/
func CalculateChecksums(r io.Reader) (ChecksumList, error) {
	sha256Hash := sha256.New()
	ripemd160Hash := ripemd160.New()

	pageSize := os.Getpagesize()
	reader := bufio.NewReaderSize(r, pageSize)
	multiWriter := io.MultiWriter(ripemd160Hash, sha256Hash)

	_, err := io.Copy(multiWriter, reader)
	if err != nil {
		return ChecksumList{}, err
	}

	ret := make(ChecksumList, 4)
	ret["sha256"] = hex.EncodeToString(sha256Hash.Sum(nil))
	ret["ripemd160"] = hex.EncodeToString(ripemd160Hash.Sum(nil))

	return ret, nil
}
