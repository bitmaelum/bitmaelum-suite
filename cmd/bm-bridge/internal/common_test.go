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
	"testing"

	"github.com/stretchr/testify/assert"
)

const complexMessage = `MIME-Version: 1.0
Date: Wed, 10 Feb 2021 11:17:16 +0100
Message-ID: <CAFdzZ3+s7OVdbawjQbq6dfig=wWng0n9kz1_cNpKS-x=XBtssQ@mail.gmail.com>
Subject: test mime
From: Antonio Calatrava <antoniocalatrava@bitmaelum.network>
To: Joshua Thijssen <jaytaph@bitmaelum.network>
Content-Type: multipart/mixed; boundary="00000000000015d01905baf8b5fa"

--00000000000015d01905baf8b5fa
Content-Type: multipart/alternative; boundary="00000000000015d01705baf8b5f8"

--00000000000015d01705baf8b5f8
Content-Type: text/plain; charset="UTF-8"

This is a test
BitMaelum <http://bitmaelum.com>

--00000000000015d01705baf8b5f8
Content-Type: text/html; charset="UTF-8"

<div dir="ltr">This is a test<div><a href="http://bitmaelum.com">BitMaelum</a><br></div></div>

--00000000000015d01705baf8b5f8--
--00000000000015d01905baf8b5fa
Content-Type: text/plain; charset="US-ASCII"; name="bitmaelum_passwords.txt"
Content-Disposition: attachment; filename="bitmaelum_passwords.txt"
Content-Transfer-Encoding: base64
X-Attachment-Id: f_kl24vyj50
Content-ID: <f_kl24vyj50>

aGVyZSBhcmUgdGhlIG1hc3RlciBwYXNzd29yZHMgZm9yIHB3bmluZyBiaXRtYWVsdW0gbmV0d29y
azoKCnlkaXNqMzluZGxhc2RhCm9ydWJrczgybjNrajJhCnVram5zZjhocnRtYXNuCmFhbmpzdWlh
c2t1bmtzCnJrc2llbmY5OGRzOG4zCmVramRzdXlza2puc2RmCmFkZmtkc2prbmRzZnNmCmhsc2Zp
dWVuZnVpc2RqCmFmaXVqYnNka3Vmc2RmCmNsc2lmOHNqZHNmaWtzCmtkZml1c2RzZmlvc3VzCmVu
a3I5NDNqbmJzZGtqCnJsa3NpZmprc2Q5b2tqCm5mb3Nka3NodWlkc2hqCm9zZGZpczgzamZrc2Rq
Cnc5ZjgzbmZpYWJkc2tmCgpwbGVhc2UgcmVhZCBmaXJzdCBsZXR0ZXIgb2YgZWFjaCAicGFzc3dv
cmQiCg==
--00000000000015d01905baf8b5fa--`

const simpleMessage = `MIME-Version: 1.0
Date: Wed, 10 Feb 2021 11:17:16 +0100
Message-ID: <CAFdzZ3+s7OVdbawjQbq6dfig=wWng0n9kz1_cNpKS-x=XBtssQ@mail.gmail.com>
Subject: Simple text/plain message
From: Antonio Calatrava <antoniocalatrava@bitmaelum.network>
To: Joshua Thijssen <jaytaph@bitmaelum.network>
Content-Type: text/plain

Heya!, testing smtp gateway

Antonio.`

func TestEncodeDecodeComplexMessage(t *testing.T) {
	m, err := DecodeFromMime(complexMessage)
	assert.Equal(t, nil, err)
	assert.Equal(t, m.From.Name, "Antonio Calatrava")
	assert.Equal(t, m.To[0].Address, "jaytaph@bitmaelum.network")
	assert.Equal(t, m.Subject, "test mime")
	assert.Equal(t, len(m.Blocks), 3)
	assert.Equal(t, len(m.Attachments), 1)

	a, err := m.EncodeToMime()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(a), 1670)
}

func TestEncodeDecodeSimpleMessage(t *testing.T) {
	m, err := DecodeFromMime(simpleMessage)
	assert.Equal(t, nil, err)
	assert.Equal(t, m.From.Name, "Antonio Calatrava")
	assert.Equal(t, m.To[0].Address, "jaytaph@bitmaelum.network")
	assert.Equal(t, m.Subject, "Simple text/plain message")
	assert.Equal(t, len(m.Blocks), 2)
	assert.Equal(t, len(m.Attachments), 0)

	a, err := m.EncodeToMime()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(a), 388)
}
