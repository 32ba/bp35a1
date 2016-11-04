package bp35a1

import (
	"bytes"
	"testing"
)

func TestNewProperty(t *testing.T) {
	epc := uint8(0xE7)
	pdc := uint8(0x04)
	edt := []byte{0x01, 0x23, 0x34, 0x56}
	p, err := NewProperty(epc, pdc, edt)
	if err != nil {
		t.Errorf("Unexpected error returned: %v", err)
	}
	if p.EPC != epc {
		t.Errorf("EPC: %v, expected %v", p.EPC, epc)
	}
	if p.PDC != pdc {
		t.Errorf("PDC: %v, expected %v", p.PDC, pdc)
	}
	if !bytes.Equal(p.EDT, edt) {
		t.Errorf("EDT: %v, expected %v", p.EDT, edt)
	}
}

func TestNewPropertyValidates(t *testing.T) {
	epc := uint8(0xE7)
	pdc := uint8(0x04)
	edt := []byte{0x01, 0x23, 0x34}
	p, err := NewProperty(epc, pdc, edt)
	if p != nil || err == nil {
		t.Errorf("Error expected")
	}
}

func TestNewPropertyBytes(t *testing.T) {
	p, _ := NewProperty(0xE7, 0x04, []byte{0x01, 0x23, 0x45, 0x67})
	b := p.Bytes()
	expected := []byte{0xE7, 0x04, 0x01, 0x23, 0x45, 0x67}
	if !bytes.Equal(b, expected) {
		t.Errorf("Bytes: %v, expected: %v", b, expected)
	}
}
