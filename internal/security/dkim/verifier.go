package dkim

import (
	"bytes"

	"github.com/emersion/go-msgauth/dkim"
)

type VerificationResult struct {
	Valid       bool
	Domain      string
	Selector    string
	Error       error
	HeaderField string
}

type Verifier struct{}

func NewVerifier() *Verifier {
	return &Verifier{}
}

func (v *Verifier) Verify(message []byte) ([]*VerificationResult, error) {
	r := bytes.NewReader(message)
	verifications, err := dkim.Verify(r)
	if err != nil {
		return []*VerificationResult{{Valid: false, Error: err}}, nil
	}

	var results []*VerificationResult
	for _, v := range verifications {
		results = append(results, &VerificationResult{
			Valid:       v.Err == nil,
			Domain:      v.Domain,
			Selector:    v.Identifier,
			Error:       v.Err,
			HeaderField: v.HeaderKeys[0],
		})
	}

	return results, nil
}
