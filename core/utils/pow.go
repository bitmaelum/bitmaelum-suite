package utils

import (
    "bytes"
    "crypto/sha256"
    "math"
    "math/big"
    "strconv"
)

func IntToHex(n int64) []byte {
    return []byte(strconv.FormatInt(n, 16))
}

func ProofOfWork(bits int, data []byte) int64 {
    var hashInt big.Int

    // Hash must be less than this
    target := big.NewInt(1)
    target = target.Lsh(target, uint(256 - bits))

    // Count from 0 to MAXINT
    var counter int64 = 0
    for counter < math.MaxInt64 {
        // SHA256 the data
        hash := sha256.Sum256(bytes.Join([][]byte{
            data,
            IntToHex(counter),
        }, []byte{}))
        hashInt.SetBytes(hash[:])

        // Is it less than our target, then we have done our work
        if hashInt.Cmp(target) == -1 {
            break;
        }

        // Higher, so we must do more work. Increase counter and try again
        counter += 1
    }

    return counter
}

func ValidateProofOfWork(bits int, data []byte, nonce int64) bool {
    var hashInt big.Int

    hash := sha256.Sum256(bytes.Join([][]byte{
        data,
        IntToHex(nonce),
    }, []byte{}))
    hashInt.SetBytes(hash[:])

    target := big.NewInt(1)
    target = target.Lsh(target, uint(256 - bits))

    return hashInt.Cmp(target) == -1
}
