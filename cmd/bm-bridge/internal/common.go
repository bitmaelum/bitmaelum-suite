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

package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"strings"
	"time"

	"net/mail"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	errAccountNotFound = errors.New("account not found")
	errRoutingNotFound = errors.New("cannot find routing ID for this account")
)

const dateLayout = "Mon, 02 Jan 2006 15:04:05 -0700"

// DefaultDomain is the default (dummy?) domain to be used to translate
// the BitMaelum address to/from
const DefaultDomain = "@bitmaelum.network"

// GatewayAddress is the gateway address to send external email messages to
const GatewayAddress = "mailgateway!"

const dummyBoundary = "stop-using-email-and-start-using-bitmaelum"

// MimeMessage contains the struct to encode or decode a mime message
// the attachments are in the format "filename" -> base64 data
type MimeMessage struct {
	ID          string
	From        *mail.Address
	To          []*mail.Address
	Subject     string
	Date        time.Time
	Blocks      []string
	Attachments map[string][]byte
}

// GetClientAndInfo will get AccountInfo and API from account in the Vault
func GetClientAndInfo(v *vault.Vault, acc string) (*vault.AccountInfo, *api.API, error) {
	info, err := vault.GetAccount(v, acc)
	if err != nil {
		return nil, nil, errAccountNotFound
	}

	resolver := container.Instance.GetResolveService()
	routingInfo, err := resolver.ResolveRouting(info.RoutingID)
	if err != nil {
		logrus.Error(err)
		return nil, nil, errRoutingNotFound
	}

	client, err := api.NewAuthenticated(*info.Address, info.GetActiveKey().PrivKey, routingInfo.Routing, nil)
	if err != nil {
		return nil, nil, err
	}

	return info, client, nil
}

// AddrToEmail will translate an address string to a valid email using DefaultDomain
func AddrToEmail(address string) string {
	address = strings.Replace(address, "@", "_", -1)
	address = strings.Replace(address, "!", "", -1)

	return address + config.Bridge.Server.SMTP.Domain
}

// EmailToAddr will translate a (mocked?) domain on the DefaultDomain to an address string
func EmailToAddr(email string) string {
	address := strings.Replace(email, config.Bridge.Server.SMTP.Domain, "!", -1)
	address = strings.Replace(address, "_", "@", -1)
	return address
}

// EncodeToMime will encode a MimeMessage struct to a MIME message
func (msg *MimeMessage) EncodeToMime() ([]byte, error) {
	if len(msg.Blocks) < 1 {
		return nil, errors.New("not a valid MimeMessage")
	}

	// If there is a mimeparts block we will use that to
	// recreate the MIME message
	if block, err := getBlock("mimeparts", msg.Blocks); err == nil {
		m, err := encodeWithMimeBlock(block, msg)
		if err != nil {
			return nil, err
		}

		return []byte(generateHeader(msg) + m), nil
	}

	// There is no mimeparts block, convert a standard Bitmaelum message to
	// a MIME valid message
	var mimeMsg string
	if len(msg.Blocks) > 1 || len(msg.Attachments) > 0 {
		mimeMsg = "Content-Type: multipart/mixed; boundary=\"" + dummyBoundary + "\"\r\n\r\n"
		for _, block := range msg.Blocks {
			parts := strings.SplitN(block, ",", 2)
			if parts[0] == "default" {
				parts[0] = "text/plain"
			}
			mimeMsg = mimeMsg + "--" + dummyBoundary + "\r\n"
			if strings.HasPrefix(parts[0], "text") {
				mimeMsg = mimeMsg + "Content-Type: " + parts[0] + "\r\n\r\n" + parts[1] + "\r\n\r\n"
			} else {
				mimeMsg = mimeMsg + "Content-Type: application/octet-stream; name=\"" + parts[0] + ".dat\"\r\n"
				mimeMsg = mimeMsg + "Content-Disposition: attachment; filename=\"" + parts[0] + ".dat\"\r\n"
				mimeMsg = mimeMsg + "Content-Transfer-Encoding: base64\r\n\r\n" + string(internal.Encode([]byte(parts[1]))) + "\r\n"
			}
		}
		for name, attachment := range msg.Attachments {
			mimeMsg = mimeMsg + "--" + dummyBoundary + "\r\n"
			mimeMsg = mimeMsg + "Content-Type: application/octet-stream; name=\"" + name + "\"\r\n"
			mimeMsg = mimeMsg + "Content-Disposition: attachment; filename=\"" + name + "\"\r\n"
			mimeMsg = mimeMsg + "Content-Transfer-Encoding: base64\r\n\r\n" + string(attachment) + "\r\n"
		}

		mimeMsg = mimeMsg + "--" + dummyBoundary + "--"
	} else {
		parts := strings.SplitN(msg.Blocks[0], ",", 2)
		if parts[0] == "default" {
			parts[0] = "text/plain"
		}

		mimeMsg = "Content-Type: " + parts[0] + "\r\n\r\n"
		mimeMsg = mimeMsg + strings.SplitN(msg.Blocks[0], ",", 2)[1] + "\r\n"
	}

	return []byte(generateHeader(msg) + mimeMsg), nil
}

func generateHeader(msg *MimeMessage) string {
	header := "MIME-Version: 1.0\r\n"
	header = header + "Message-Id: " + msg.ID + "\r\n"
	if msg.From != nil {
		header = header + "From: \"" + msg.From.Name + "\" <" + msg.From.Address + ">\r\n"
	}
	if msg.To != nil && len(msg.To) > 0 {
		header = header + "To: \"" + msg.To[0].Name + "\" <" + msg.To[0].Address + ">\r\n"
	}
	header = header + "Subject: " + msg.Subject + "\r\n"
	header = header + "Date: " + msg.Date.Format(dateLayout) + "\r\n"
	return header
}

func encodeWithMimeBlock(mimeParts string, msg *MimeMessage) (string, error) {
	var mimeMsg string
	parts := strings.Split(mimeParts, "\r\n\r\n")
	for _, part := range parts {
		mimeMsg = mimeMsg + part + "\r\n\r\n"

		// Check if it's a multipart or inline/attachment
		isAttachment := false
		var fileName string
		var blockID string
		lines := strings.Split(part, "\r\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Content-Disposition") {
				disposition, params, _ := mime.ParseMediaType(strings.Split(line, ":")[1])
				if strings.HasPrefix(disposition, "attachment") || strings.HasPrefix(disposition, "inline") {
					isAttachment = true
					fileName = params["filename"]
					break
				}
			}

			if strings.HasPrefix(line, "X-Bitmaelum-Block") {
				blockID, _, _ = mime.ParseMediaType(strings.Split(line, ":")[1])
			}
		}

		if isAttachment {
			mimeMsg = mimeMsg + string(msg.Attachments[fileName]) + "\r\n"
			continue
		}

		if block, err := getBlock(blockID, msg.Blocks); err == nil {
			mimeMsg = mimeMsg + block + "\r\n"
		}
	}

	return mimeMsg, nil
}

// DecodeFromMime will decode a mime message and return a MimeMessage struct
func DecodeFromMime(m string) (*MimeMessage, error) {
	// Make sure it's using CRLF to be fully MIME compliant
	m = strings.Replace(m, "\r\n", "\n", -1)
	m = strings.Replace(m, "\n", "\r\n", -1)

	mDecodedMsg := &MimeMessage{}

	msg, err := mail.ReadMessage(bytes.NewBufferString(m))
	if err != nil {
		return mDecodedMsg, err
	}

	mDecodedMsg.To, _ = (&mail.AddressParser{}).ParseList(msg.Header.Get("To"))
	mDecodedMsg.From, _ = (&mail.AddressParser{}).Parse(msg.Header.Get("From"))
	mDecodedMsg.Subject = msg.Header.Get("Subject")
	mDecodedMsg.Date, _ = msg.Header.Date()
	mDecodedMsg.ID = msg.Header.Get("Message-Id")

	var mimeParts string
	var hasDefaultBlock bool
	mDecodedMsg.Attachments = make(map[string][]byte)

	for name, values := range msg.Header {
		if !strings.HasPrefix(name, "Content") {
			continue
		}
		for _, val := range values {
			mimeParts = fmt.Sprintf("%s%s: %s\r\n", mimeParts, name, decodeRFC2047(val))
		}
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mimeParts, mDecodedMsg.Blocks, mDecodedMsg.Attachments, err = readParts(mimeParts, mDecodedMsg.Blocks, mDecodedMsg.Attachments, msg.Body, params["boundary"], 0)
	} else {
		mimeParts = fmt.Sprintf("%sX-Bitmaelum-Block: default", mimeParts)
		body, err := ioutil.ReadAll(msg.Body)
		if err != nil {
			return mDecodedMsg, err
		}

		hasDefaultBlock = true
		mDecodedMsg.Blocks = append(mDecodedMsg.Blocks, "default,"+string(body))
	}

	// If there is no default block then set the first "text/plain" as default
	if !hasDefaultBlock {
		if err := renameBlock("text/plain", "default", &mDecodedMsg.Blocks); err == nil {
			hasDefaultBlock = true
			mimeParts = strings.Replace(mimeParts, "text/plain", "default", 1)
		}
	}

	// If no "text/plain" is found then set the first "text/html" as default
	if !hasDefaultBlock {
		if err := renameBlock("text/html", "default", &mDecodedMsg.Blocks); err == nil {
			hasDefaultBlock = true
			mimeParts = strings.Replace(mimeParts, "text/html", "default", 1)
		}
	}

	// It no "text/html" and not "text/plain" are found then set the first block as default
	if !hasDefaultBlock {
		parts := strings.SplitN(mDecodedMsg.Blocks[0], ",", 2)
		mimeParts = strings.Replace(mimeParts, parts[0], "default", 1)
		mDecodedMsg.Blocks[0] = "default," + parts[1]
	}

	mDecodedMsg.Blocks = append(mDecodedMsg.Blocks, "mimeparts,"+mimeParts)

	return mDecodedMsg, err
}

func decodeRFC2047(s string) string {
	// GO 1.5 does not decode headers, but this may change in future releases...
	decoded, err := (&mime.WordDecoder{}).DecodeHeader(s)
	if err != nil || len(decoded) == 0 {
		return s
	}
	return decoded
}

func readParts(mimeParts string, blocks []string, attachments map[string][]byte, body io.Reader, boundary string, idx int) (string, []string, map[string][]byte, error) {
	mr := multipart.NewReader(body, boundary)
	for {
		msg, err := mr.NextPart()
		if err == io.EOF {
			mimeParts = fmt.Sprintf("%s\r\n--"+boundary+"--", mimeParts)
			return mimeParts, blocks, attachments, nil
		}
		if err != nil {
			return mimeParts, blocks, attachments, err
		}

		mimeParts = fmt.Sprintf("%s\r\n--"+boundary+"\r\n", mimeParts)
		// decode any Q-encoded values
		for name, values := range msg.Header {
			for _, val := range values {
				mimeParts = fmt.Sprintf("%s%s: %s\r\n", mimeParts, name, decodeRFC2047(val))
			}
		}

		mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
		if err != nil {
			return mimeParts, blocks, attachments, err
		}
		if strings.HasPrefix(mediaType, "multipart/") {
			mimeParts, blocks, attachments, err = readParts(mimeParts, blocks, attachments, msg, params["boundary"], idx+1)
			if err != nil {
				return mimeParts, blocks, attachments, err
			}
			continue
		}

		body, err := ioutil.ReadAll(msg)
		if err != nil {
			return mimeParts, blocks, attachments, err
		}

		disposition, data, err := mime.ParseMediaType(msg.Header.Get("Content-Disposition"))
		if err == nil && (disposition == "attachment" || disposition == "inline") {
			attachments[data["filename"]] = body
		} else {
			blockID, _ := uuid.NewRandom()
			blocks = append(blocks, blockID.String()+","+string(body))
			mimeParts = fmt.Sprintf("%sX-Bitmaelum-Block: %s\r\n", mimeParts, blockID.String())
		}

	}

}

func getBlock(id string, blocks []string) (string, error) {
	if id == "" {
		return "", errors.New("block id not specified")
	}

	for _, block := range blocks {
		parts := strings.SplitN(block, ",", 2)
		if parts[0] == id {
			return parts[1], nil
		}
	}

	return "", errors.New("block not found")
}

func renameBlock(from string, to string, blocks *[]string) error {
	for idx, block := range *blocks {

		if strings.HasPrefix(block, from) {
			parts := strings.SplitN(block, ",", 2)
			(*blocks)[idx] = to + "," + parts[1]
			return nil
		}

	}

	return errors.New("block not found")
}
