package checksum

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "encoding/hex"
    "github.com/jaytaph/mailv2/core/message"
    "hash/crc32"
)

// @TODO: Multi hash through http://marcio.io/2015/07/calculating-multiple-file-hashes-in-a-single-pass/

func Sha1(data []byte) message.ChecksumType {
    hasher := sha1.New()
    _, _  = hasher.Write([]byte(data))

    return message.ChecksumType{
        Hash: "sha1",
        Value: hex.EncodeToString(hasher.Sum(nil)),
    }
}

func Sha256(data []byte) message.ChecksumType {
    hasher := sha256.New()
    _, _  = hasher.Write([]byte(data))

    return message.ChecksumType{
        Hash: "sha256",
        Value: hex.EncodeToString(hasher.Sum(nil)),
    }
}

func Crc32(data []byte) message.ChecksumType {
    hasher := crc32.NewIEEE()
    _, _  = hasher.Write([]byte(data))

    return message.ChecksumType{
        Hash: "crc32",
        Value: hex.EncodeToString(hasher.Sum(nil)),
    }
}

func Md5(data []byte) message.ChecksumType {
    hasher := md5.New()
    _, _ = hasher.Write([]byte(data))

    return message.ChecksumType{
        Hash: "md5",
        Value: hex.EncodeToString(hasher.Sum(nil)),
    }
}
