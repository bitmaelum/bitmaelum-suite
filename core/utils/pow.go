package utils

import (
    "bytes"
    "crypto/sha256"
    "encoding/hex"
    "math"
    "math/big"
    "strconv"
)

func intToHex(n uint64) []byte {
    return []byte(strconv.FormatUint(n, 16))
}

func ProofOfWork(bits int, data []byte) string {
    var hashInt big.Int

    // Hash must be less than this
    target := big.NewInt(1)
    target = target.Lsh(target, uint(256 - bits))

    // Count from 0 to MAXINT
    var counter uint64 = 0
    for counter < math.MaxInt64 {
        // SHA256 the data
        hash := sha256.Sum256(bytes.Join([][]byte{
            data,
            intToHex(counter),
        }, []byte{}))
        hashInt.SetBytes(hash[:])

        // Is it less than our target, then we have done our work
        if hashInt.Cmp(target) == -1 {
            break;
        }

        // Higher, so we must do more work. Increase counter and try again
        counter += 1
    }

    sum := sha256.Sum256(intToHex(counter))
    return hex.EncodeToString(sum[:])
}

func ValidateProofOfWork(bits int, data []byte, proof uint64) bool {
    var hashInt big.Int

    hash := sha256.Sum256(bytes.Join([][]byte{
        data,
        intToHex(proof),
    }, []byte{}))
    hashInt.SetBytes(hash[:])

    target := big.NewInt(1)
    target = target.Lsh(target, uint(256 - bits))

    return hashInt.Cmp(target) == -1
}
