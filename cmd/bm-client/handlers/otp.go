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

package handlers

import (
	"crypto/aes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

// OtpGenerate will generate an OTP valid for otpServer
func OtpGenerate(info *internal.AccountInfo, otpServer *string) {
	// Get
	recs, err := net.LookupTXT("_bitmaelum." + *otpServer)
	if err != nil {

		logrus.Fatal(err)
	}

	resolver := container.GetResolveService()

	for _, txt := range recs {
		orgHash, err := hash.NewFromHash(strings.ToLower(txt))
		if err != nil {
			continue
		}

		oi, err := resolver.ResolveOrganisation(*orgHash)
		if err == nil {
			// Generate secret on the client and compute OTP
			secret, err := bmcrypto.KeyExchange(info.PrivKey, oi.PublicKey)
			if err != nil {
				logrus.Fatal(err)
			}

			printOtpLoop(secret, *otpServer)

			return
		}
	}

	logrus.Fatal(errors.New("public key not found for " + *otpServer))
}

func printOtpLoop(secret []byte, server string) {

	for {
		otp := computeOTPFromSecret(secret, 8)

		v := 60 - ((time.Now().UnixNano() - getLastMinuteTimestamp()) / int64(time.Second))

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"OTP", "Server", "Valid for"})
		table.Append([]string{otp, server, fmt.Sprintf("%d", v)})

		table.Render()

		<-time.After(1 * time.Second)
	}
}

func getLastMinuteTimestamp() int64 {
	t := time.Now()
	rounded := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	return rounded.UnixNano()
}

func computeOTPFromSecret(secret []byte, otpLength int) string {
	// Create a byte array to store current timestamp to be
	// encrypted (we only need 8 bytes since its a 64bit
	// integer but since AES block size is 16 byte long we
	// need to match it)
	plain := make([]byte, 16)
	binary.LittleEndian.PutUint64(plain, uint64(getLastMinuteTimestamp()))

	// Create a AES block cipher with the secret
	block, _ := aes.NewCipher(secret)

	// And encrypt the byte-array timestamp
	encrypted := make([]byte, len(plain))
	block.Encrypt(encrypted, []byte(plain))

	// Convert the result to integer
	i := binary.LittleEndian.Uint64(encrypted)

	// And get the last otpLength numbers
	otp := fmt.Sprintf("%0*d", otpLength, i)

	return otp[len(otp)-otpLength:]
}
