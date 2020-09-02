package message

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/ripemd160"
	"io"
	"os"
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
