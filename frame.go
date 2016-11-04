package bp35a1

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strconv"
)

type Frame struct {
	EHD1       uint8
	EHD2       uint8
	TID        uint16
	SEOJ       uint32 // Use only 24 bits
	DEOJ       uint32 // Use only 24 bits
	ESV        uint8
	OPC        uint8
	Properties []*Property
}

func NewFrame(ehd1, ehd2 uint8, tid uint16, seoj, deoj uint32, esv, opc uint8, props []*Property) *Frame {
	return &Frame{
		EHD1:       ehd1,
		EHD2:       ehd2,
		TID:        tid,
		SEOJ:       seoj,
		DEOJ:       deoj,
		ESV:        esv,
		OPC:        opc,
		Properties: props,
	}
}

func NewFrameFromString(s string) (*Frame, error) {
	ehd1, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return nil, err
	}

	ehd2, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return nil, err
	}

	tid, err := strconv.ParseUint(s[4:8], 16, 16)
	if err != nil {
		return nil, err
	}

	seoj, err := strconv.ParseUint(s[8:14], 16, 24)
	if err != nil {
		return nil, err
	}

	deoj, err := strconv.ParseUint(s[14:20], 16, 24)
	if err != nil {
		return nil, err
	}

	esv, err := strconv.ParseUint(s[20:22], 16, 8)
	if err != nil {
		return nil, err
	}

	opc, err := strconv.ParseUint(s[22:24], 16, 8)
	if err != nil {
		return nil, err
	}

	pd := s[24:]
	var props []*Property

	for i := 0; i < int(opc); i++ {
		epc, err := strconv.ParseUint(pd[0:2], 16, 8)
		if err != nil {
			return nil, err
		}

		pdc, err := strconv.ParseUint(pd[2:4], 16, 8)
		if err != nil {
			return nil, err
		}

		edt, err := hex.DecodeString(pd[4 : 4+pdc*2])
		if err != nil {
			return nil, err
		}

		prop, err := NewProperty(uint8(epc), uint8(pdc), edt)
		if err != nil {
			return nil, err
		}

		props = append(props, prop)
	}

	return &Frame{
		EHD1:       uint8(ehd1),
		EHD2:       uint8(ehd2),
		TID:        uint16(tid),
		SEOJ:       uint32(seoj),
		DEOJ:       uint32(deoj),
		ESV:        uint8(esv),
		OPC:        uint8(opc),
		Properties: props,
	}, nil
}

func (f *Frame) Bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, f.EHD1)
	_ = binary.Write(buf, binary.BigEndian, f.EHD2)
	_ = binary.Write(buf, binary.BigEndian, f.TID)
	_ = binary.Write(buf, binary.BigEndian, uint8((f.SEOJ>>16)&0xff))
	_ = binary.Write(buf, binary.BigEndian, uint8((f.SEOJ>>8)&0xff))
	_ = binary.Write(buf, binary.BigEndian, uint8((f.SEOJ)&0xff))
	_ = binary.Write(buf, binary.BigEndian, uint8((f.DEOJ>>16)&0xff))
	_ = binary.Write(buf, binary.BigEndian, uint8((f.DEOJ>>8)&0xff))
	_ = binary.Write(buf, binary.BigEndian, uint8((f.DEOJ)&0xff))
	_ = binary.Write(buf, binary.BigEndian, f.ESV)
	_ = binary.Write(buf, binary.BigEndian, f.OPC)
	for _, p := range f.Properties {
		_, _ = buf.Write(p.Bytes())
	}
	return buf.Bytes()
}
