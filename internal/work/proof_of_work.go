package work

import (
	"encoding/json"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

type PowRepo struct {
	W proofofwork.ProofOfWork
}

type PowResultType struct {
	Proof uint64
}

func NewPow() (*PowRepo, error) {
	work, err := proofofwork.GenerateWorkData()
	if err != nil {
		return nil, err
	}

	return &PowRepo{
		W: *proofofwork.NewWithoutProof(config.Server.Accounts.ProofOfWork, work),
	}, nil
}

func (p *PowRepo) GetName() string {
	return "pow"
}

func (p *PowRepo) GetWorkOutput() map[string]interface{} {
	return map[string]interface{}{
		"bits": p.W.Bits,
		"data": p.W.Data,
	}
}

func (p *PowRepo) GetWorkProofOutput() map[string]interface{} {
	return map[string]interface{}{
		"proof": p.W.Proof,
	}
}


func (p *PowRepo) Work() {
	p.W.WorkMulticore()
}

func (p *PowRepo) ValidateWork(data []byte) bool {
	res := &PowResultType{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		return false
	}

	p.W.Proof = res.Proof
	return p.W.HasDoneWork() && p.W.IsValid();
}




