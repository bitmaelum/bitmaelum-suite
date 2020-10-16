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

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordWrap(t *testing.T) {
	s := "NEWS PLAY ELSE CABLE UNLOCK SUSPECT TOAST MIXTURE SCARE POTTERY MONTH ESSAY IMMUNE BURGER STING RATE PLANET TOWN"

	assert.Equal(t, "NEWS PLAY ELSE CABLE UNLOCK SUSPECT\nTOAST MIXTURE SCARE POTTERY MONTH ESSAY\nIMMUNE BURGER STING RATE PLANET TOWN", WordWrap(s, 40))
	assert.Equal(t, "NEWS PLAY ELSE CABLE UNLOCK SUSPECT TOAST MIXTURE SCARE\nPOTTERY MONTH ESSAY IMMUNE BURGER STING RATE PLANET TOWN", WordWrap(s, 60))
	assert.Equal(t, "NEWS PLAY\nELSE CABLE\nUNLOCK\nSUSPECT\nTOAST\nMIXTURE\nSCARE\nPOTTERY\nMONTH\nESSAY\nIMMUNE\nBURGER\nSTING RATE\nPLANET\nTOWN", WordWrap(s, 10))
	assert.Equal(t, "\nNEWS\nPLAY\nELSE\nCABLE\nUNLOCK\nSUSPECT\nTOAST\nMIXTURE\nSCARE\nPOTTERY\nMONTH\nESSAY\nIMMUNE\nBURGER\nSTING\nRATE\nPLANET\nTOWN", WordWrap(s, 3))
}
