package utils

import (
    "bytes"
    "crypto/sha256"
    "math"
    "math/big"
    "strconv"
)

//
// Proof of work consists of SHA256 hashing data with an additional counter
//
//      counter = 0
//      do
//          counter += 1
//          hash = SHA256(data + counter)
//      until X number of left bits of hash is 0
//
// The idea is that is takes time to calculate a SHA256 from data. Because
// the SHA256 is distributed evenly, there is no way to know how many bits
// on the left are 0. By increasing the counter with 1 every time and hashing
// again, we get many different hashes. As soon as we found a hash with
// X bits 0 on the left, we use the counter value as the proof for the work
// that has been done.
//



// convert a large number to hexidecimal bytes
func intToHex(n uint64) []byte {
    return []byte(strconv.FormatUint(n, 16))
}

// Does a proof-of-work for X number of bits based on the data
// This function can take a long time
func ProofOfWork(bits int, data []byte) uint64 {
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

    return counter
}

// Validate if the proof for the given data has been completed
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
