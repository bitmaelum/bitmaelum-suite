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

package work

import (
	"encoding/json"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// PowRepo is the repository for proof of work
type PowRepo struct {
	W proofofwork.ProofOfWork
}

// PowResultType is the result that gets returned when done the work
type PowResultType struct {
	Proof uint64
}

// NewPow will return a new proof-of-work repository filled
func NewPow() (*PowRepo, error) {
	work, err := proofofwork.GenerateWorkData()
	if err != nil {
		return nil, err
	}

	return &PowRepo{
		W: *proofofwork.NewWithoutProof(config.Server.Accounts.ProofOfWork, work),
	}, nil
}

// GetName will return the name of the work type
func (p *PowRepo) GetName() string {
	return "pow"
}

// GetWorkOutput will return a list of data that will be returned in the ticket
func (p *PowRepo) GetWorkOutput() map[string]interface{} {
	return map[string]interface{}{
		"bits": p.W.Bits,
		"data": p.W.Data,
	}
}

// GetWorkProofOutput will return the proof / work done that will be send back to the server
func (p *PowRepo) GetWorkProofOutput() map[string]interface{} {
	return map[string]interface{}{
		"proof": p.W.Proof,
	}
}

// Work will actually do the work
func (p *PowRepo) Work() {
	p.W.WorkMulticore()
}

// ValidateWork will validate the work data
func (p *PowRepo) ValidateWork(data []byte) bool {
	res := &PowResultType{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		return false
	}

	p.W.Proof = res.Proof
	return p.W.HasDoneWork() && p.W.IsValid()
}
