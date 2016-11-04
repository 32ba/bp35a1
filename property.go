package bp35a1

import (
	"errors"
)

type Property struct {
	EPC uint8
	PDC uint8
	EDT []byte
}

func NewProperty(epc uint8, pdc uint8, edt []byte) (*Property, error) {
	if int(pdc) != len(edt) {
		return nil, errors.New("pdc and length of edt do not match")
	}

	return &Property{
		EPC: epc,
		PDC: pdc,
		EDT: edt,
	}, nil
}

func (p *Property) Bytes() []byte {
	bytes := []byte{p.EPC, p.PDC}
	return append(bytes, p.EDT...)
}
