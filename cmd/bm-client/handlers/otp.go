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

package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"os"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

const blockPeriod = 30

// OtpGenerate will generate an OTP valid for otpSite
func OtpGenerate(info *vault.AccountInfo, otpSite *string) {
	// Get
	keys, err := net.LookupTXT("_otp." + *otpSite)
	if err != nil {
		logrus.Fatal(err)
	}

	for _, txt := range keys {
		key, err := bmcrypto.PublicKeyFromString(txt)
		if err == nil {
			if !key.Type.CanKeyExchange() {
				continue
			}

			// Generate secret on the client and compute OTP
			secret, err := bmcrypto.KeyExchange(info.GetActiveKey().PrivKey, *key)
			if err != nil {
				logrus.Fatal(err)
			}

			printOtpLoop(secret, *otpSite)

			return
		}
	}

	logrus.Fatal(fmt.Errorf("The site does not support OTP authentication: " + *otpSite))
}

func printOtpLoop(secret []byte, site string) {
	data := bytes.Join([][]byte{secret, []byte(site)}, nil)
	for {
		otp := computeOTPFromSecret(data, 8)

		v := blockPeriod - ((time.Now().UnixNano() - getBlockTimestamp()) / int64(time.Second))

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"OTP", "Site", "Valid for"})
		table.Append([]string{otp, site, fmt.Sprintf("%d", v)})

		table.Render()

		<-time.After(1 * time.Second)
		fmt.Printf("\033[5A")
	}
}

func getBlockTimestamp() int64 {
	t := time.Now().UnixNano()
	t = (t / (blockPeriod * int64(time.Second)) * (blockPeriod * int64(time.Second)))
	return t
}

func computeOTPFromSecret(secret []byte, length int) string {

	// Put the last nano timestamp int64 into a byte array
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(getBlockTimestamp()))

	// Compute an HMAC with the secret
	mac := hmac.New(sha256.New, secret)
	mac.Write(buf)
	sum := mac.Sum(nil)

	// From https://github.com/pquerna/otp/blob/3006c03e19424e57e998d0faa7afe846b291ca14/hotp/hotp.go#L101
	// "Dynamic truncation" in RFC 4226
	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & 0xf

	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	// Get the modulus to get the length we need
	mod := int32(value % int64(math.Pow10(length)))

	return fmt.Sprintf("%0*d", length, mod)
}
