package core

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/bitmaelum/bitmaelum-server/internal/message"
	"hash/crc32"
)

// @TODO: Multi hash through http://marcio.io/2015/07/calculating-multiple-file-hashes-in-a-single-pass/

// Sha1 Return SHA1 checksum structure from data
func Sha1(data []byte) message.Checksum {
	hasher := sha1.New()
	_, _ = hasher.Write([]byte(data))

	return message.Checksum{
		Hash:  "sha1",
		Value: hex.EncodeToString(hasher.Sum(nil)),
	}
}

// Sha256 Return SHA256 checksum structure from data
func Sha256(data []byte) message.Checksum {
	hasher := sha256.New()
	_, _ = hasher.Write([]byte(data))

	return message.Checksum{
		Hash:  "sha256",
		Value: hex.EncodeToString(hasher.Sum(nil)),
	}
}

// Crc32 Return CRC32 checksum structure from data
func Crc32(data []byte) message.Checksum {
	hasher := crc32.NewIEEE()
	_, _ = hasher.Write([]byte(data))

	return message.Checksum{
		Hash:  "crc32",
		Value: hex.EncodeToString(hasher.Sum(nil)),
	}
}

// Md5 Return MD5 checksum structure from data
func Md5(data []byte) message.Checksum {
	hasher := md5.New()
	_, _ = hasher.Write([]byte(data))

	return message.Checksum{
		Hash:  "md5",
		Value: hex.EncodeToString(hasher.Sum(nil)),
	}
}
