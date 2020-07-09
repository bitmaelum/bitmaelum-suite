package internal

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"io"
	"os"
)

// CalculateChecksums calculates a number of hashes for the given reader in one go.
// Taken from http://marcio.io/2015/07/calculating-multiple-file-hashes-in-a-single-pass/
func CalculateChecksums(r io.Reader) (message.ChecksumList, error) {
	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()
	sha512Hash := sha512.New()

	pageSize := os.Getpagesize()
	reader := bufio.NewReaderSize(r, pageSize)
	multiWriter := io.MultiWriter(md5Hash, sha1Hash, sha256Hash, sha512Hash)

	_, err := io.Copy(multiWriter, reader)
	if err != nil {
		return message.ChecksumList{}, err
	}

	ret := make(message.ChecksumList, 4)
	ret["md5"] = hex.EncodeToString(md5Hash.Sum(nil))
	ret["sha1"] = hex.EncodeToString(sha1Hash.Sum(nil))
	ret["sha256"] = hex.EncodeToString(sha256Hash.Sum(nil))
	ret["sha512"] = hex.EncodeToString(sha512Hash.Sum(nil))

	return ret, nil
}
