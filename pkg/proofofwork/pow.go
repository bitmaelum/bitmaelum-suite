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

package proofofwork

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"runtime"
	"strconv"
	"strings"
)

var errIncorrectFormat = errors.New("incorrect proof-of-work format")

// Override for testing purposes
var randomReader = rand.Reader

/*
 * Proof of work consists of SHA256 hashing data with an additional counter
 *
 *      counter = 0
 *      do
 *          counter += 1
 *          hash = SHA256(data + counter)
 *      until X number of left bits of hash is 0[^1]
 *
 * The idea is that is takes time to calculate a SHA256 from data. Because
 * the SHA256 is distributed evenly, there is no way to know how many bits
 * on the left are 0. By increasing the counter with 1 every time and hashing
 * again, we get many different hashes. As soon as we found a hash with
 * X bits 0 on the left, we use the counter value as the proof for the work
 * that has been done.
 *
 * ^1 Even though we are checking against bits to be 0, the actual check is if
 *    the hash is lower than the target hash.
 */

// ProofOfWork represents a proof-of-work which either can be completed or not
type ProofOfWork struct {
	Bits  int    `json:"bits"`
	Data  string `json:"data"`
	Proof uint64 `json:"proof,omitempty"`
}

// New generates a new ProofOfWork structure.
func New(bits int, data string, proof uint64) *ProofOfWork {
	pow := &ProofOfWork{
		Bits:  bits,
		Data:  data,
		Proof: proof,
	}

	return pow
}

// NewWithoutProof returns a new proof-of-work without actual proof (needs to be worked on)
func NewWithoutProof(bits int, data string) *ProofOfWork {
	return New(bits, data, 0)
}

// NewFromString generates a proof-of-work based on the given string
func NewFromString(s string) (*ProofOfWork, error) {
	parts := strings.SplitN(s, "$", 3)
	if len(parts) != 3 {
		return nil, errIncorrectFormat
	}

	bits, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	proof, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}

	return &ProofOfWork{
		Bits:  bits,
		Data:  string(data),
		Proof: uint64(proof),
	}, nil
}

// GenerateWorkData generates random work
func GenerateWorkData() (string, error) {
	data := make([]byte, 32)
	_, err := io.ReadFull(randomReader, data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// HasDoneWork returns true if this instance already has done proof-of-work
func (pow *ProofOfWork) HasDoneWork() bool {
	return pow.Proof > 0
}

// Return a representation of proof-of-work
func (pow *ProofOfWork) String() string {
	return fmt.Sprintf("%d$%s$%d", pow.Bits, base64.StdEncoding.EncodeToString([]byte(pow.Data)), pow.Proof)
}

// WorkMulticore will work on all cores possible
func (pow *ProofOfWork) WorkMulticore() {
	// Use the maximum number of cores for work to be the most efficient.
	pow.Work(maxCores())
}

// Work actually does the proof-of-work on the given amount of cores
func (pow *ProofOfWork) Work(cores int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	found := make(chan uint64)
	for i := 0; i < cores; i++ {
		go pow.workOnCore(ctx, uint64(i), uint64(cores), found)
	}
	pow.Proof = <-found
}

// workOnCore does hashing on a core, starting at start and increasing with step each time until
// the context is cancelled (completed), or we found the proof, in which case we send it to the
// channel.
func (pow *ProofOfWork) workOnCore(ctx context.Context, start, step uint64, found chan uint64) {
	var hashInt big.Int

	// Hash must be less than this
	target := big.NewInt(1)
	target = target.Lsh(target, uint(256-pow.Bits))

	var counter uint64 = start
	for counter < math.MaxInt64 {
		select {
		case <-ctx.Done():
			// Break when context is cancelled (other go thread has found our proof)
			return
		default:
			// 1st round of SHA256
			hash := sha256.Sum256(bytes.Join([][]byte{
				[]byte(pow.Data),
				intToHex(counter),
			}, []byte{}))

			// 2nd round of SHA256
			hash = sha256.Sum256(hash[:])
			hashInt.SetBytes(hash[:])

			// Is it less than our target, then we have done our work
			if hashInt.Cmp(target) == -1 {
				found <- counter
				return
			}
		}

		counter += step
	}
}

// MarshalJSON marshals a pow into bytes
func (pow *ProofOfWork) MarshalJSON() ([]byte, error) {
	return json.Marshal(pow.String())
}

// UnmarshalJSON unmarshals bytes into a pow
func (pow *ProofOfWork) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	powCopy, err := NewFromString(s)
	if err != nil {
		return err
	}

	// This seems wrong, but copy() doesn't work?
	pow.Bits = powCopy.Bits
	pow.Data = powCopy.Data
	pow.Proof = powCopy.Proof

	return err
}

// IsValid returns true when the given work can be validated against the proof
func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int

	// 1st round
	hash := sha256.Sum256(bytes.Join([][]byte{
		[]byte(pow.Data),
		intToHex(pow.Proof),
	}, []byte{}))

	// 2nd round of SHA256
	hash = sha256.Sum256(hash[:])
	hashInt.SetBytes(hash[:])

	target := big.NewInt(1)
	target = target.Lsh(target, uint(256-pow.Bits))

	return hashInt.Cmp(target) == -1
}

// intToHex converts a large number to hexadecimal bytes
func intToHex(n uint64) []byte {
	return []byte(strconv.FormatUint(n, 16))
}

// maxCores returns the number of cores in the current system
func maxCores() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
