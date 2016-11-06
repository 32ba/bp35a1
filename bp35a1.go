package bp35a1

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	// "log"
	"strconv"
	"strings"

	"github.com/tarm/serial"
)

type PAN struct {
	Channel     uint8
	ChannelPage uint8
	PanID       uint16
	Addr        string
	LQI         uint8
	PairID      string
}

type PortWrapper struct {
	Port    io.Reader
	scanner *bufio.Scanner
}

func NewPortWrapper(r io.Reader) *PortWrapper {
	return &PortWrapper{
		Port:    r,
		scanner: bufio.NewScanner(r),
	}
}

func (p *PortWrapper) ScanAndText() (bool, string) {
	f := p.scanner.Scan()

	if !f {
		return f, ""
	}

	t := p.scanner.Text()
	// log.Print("<- " + t)

	return f, t
}

func (p *PortWrapper) ReadLine() string {
	p.scanner.Scan()
	t := p.scanner.Text()
	// log.Print("<- " + t)
	return t
}

func (p *PortWrapper) ReadLinesUntilOk() []string {
	var ls []string
	for p.scanner.Scan() {
		t := p.scanner.Text()
		// log.Print("<- " + t)
		ls = append(ls, t)
		if t == "OK" {
			break
		}
	}
	return ls
}

type BP35A1 struct {
	Device string
	Baud   int
	Port   *serial.Port
}

func NewBP35A1(device string, baud int) *BP35A1 {
	return &BP35A1{
		Device: device,
		Baud:   baud,
	}
}

func (b *BP35A1) Connect() error {
	c := &serial.Config{
		Name: b.Device,
		Baud: b.Baud,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	b.Port = s
	return nil
}

func (b *BP35A1) Close() {
	b.Port.Close()
}

func (b *BP35A1) getWrappedScanner() *PortWrapper {
	return NewPortWrapper(b.Port)
}

func (b *BP35A1) write(s string) error {
	// log.Print("-> " + s)
	_, err := b.Port.Write([]byte(s))
	if err != nil {
		return err
	}
	return nil
}

func (b *BP35A1) writeBytes(bs []byte) error {
	// log.Print("-> ", bs)
	// log.Print(fmt.Sprintf("-> %v", bs))
	_, err := b.Port.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func (b *BP35A1) SKINFO() (string, error) {
	err := b.write("SKINFO\r\n")
	if err != nil {
		return "", err
	}

	pw := b.getWrappedScanner()
	ls := pw.ReadLinesUntilOk()
	if len(ls) != 3 {
		return "", errors.New("invalid response")
	}

	if ls[0] != "SKINFO" {
		return "", errors.New("invalid echo back received: " + ls[0])
	}

	return ls[1], nil
}

func (b *BP35A1) SKVER() (string, error) {
	err := b.write("SKVER\r\n")
	if err != nil {
		return "", err
	}

	pw := b.getWrappedScanner()
	ls := pw.ReadLinesUntilOk()
	if len(ls) != 3 {
		return "", errors.New("invalid response")
	}

	if ls[0] != "SKVER" {
		return "", errors.New("invalid echo back received: " + ls[0])
	}

	return ls[1], nil
}

func (b *BP35A1) SKSETPWD(pwd string) error {
	err := b.write("SKSETPWD C " + pwd + "\r\n")
	if err != nil {
		return err
	}

	pw := b.getWrappedScanner()
	// [SKSETPWD C PASSWORD_HERE OK] is returned
	_ = pw.ReadLinesUntilOk()

	return nil
}

func (b *BP35A1) SKSETRBID(rbid string) error {
	err := b.write("SKSETRBID " + rbid + "\r\n")
	if err != nil {
		return err
	}

	pw := b.getWrappedScanner()
	_ = pw.ReadLinesUntilOk()

	return nil
}

func (b *BP35A1) SKSREG(sreg string, val string) error {
	err := b.write("SKSREG " + sreg + " " + val + "\r\n")
	if err != nil {
		return err
	}

	pw := b.getWrappedScanner()
	_ = pw.ReadLinesUntilOk()

	return nil
}

func (b *BP35A1) SKLL64(addr string) (string, error) {
	err := b.write("SKLL64 " + addr + "\r\n")
	if err != nil {
		return "", err
	}

	pw := b.getWrappedScanner()
	_ = pw.ReadLine()

	return pw.ReadLine(), nil
}

func (b *BP35A1) SKSCAN() (*PAN, error) {
	err := b.write("SKSCAN 2 FFFFFFFF 6\r\n")
	if err != nil {
		return nil, err
	}

	pw := b.getWrappedScanner()

	_ = pw.ReadLinesUntilOk()

	for {
		flag, t := pw.ScanAndText()
		if !flag {
			break
		}
		if strings.HasPrefix(t, "EVENT 20") {
			break
		}
	}

	pan := &PAN{}

	for {
		flag, t := pw.ScanAndText()
		if !flag {
			break
		}
		switch {
		case strings.HasPrefix(t, "  Channel:"):
			ui, err := strconv.ParseUint(strings.Split(t, ":")[1], 16, 8)
			if err != nil {
				break
			}
			pan.Channel = uint8(ui)
		case strings.HasPrefix(t, "  Channel Page:"):
			ui, err := strconv.ParseUint(strings.Split(t, ":")[1], 16, 8)
			if err != nil {
				break
			}
			pan.ChannelPage = uint8(ui)
		case strings.HasPrefix(t, "  Pan ID:"):
			ui, err := strconv.ParseUint(strings.Split(t, ":")[1], 16, 16)
			if err != nil {
				break
			}
			pan.PanID = uint16(ui)
		case strings.HasPrefix(t, "  Addr:"):
			pan.Addr = strings.Split(t, ":")[1]
		case strings.HasPrefix(t, "  LQI:"):
			ui, err := strconv.ParseUint(strings.Split(t, ":")[1], 16, 8)
			if err != nil {
				break
			}
			pan.LQI = uint8(ui)
		case strings.HasPrefix(t, "  PairID:"):
			pan.PairID = strings.Split(t, ":")[1]
		}

		if strings.HasPrefix(t, "EVENT 22") {
			break
		}
	}

	return pan, nil
}

func (b *BP35A1) SKJOIN(addr string) (bool, error) {
	err := b.write("SKJOIN " + addr + "\r\n")
	if err != nil {
		return false, err
	}

	pw := b.getWrappedScanner()
	_ = pw.ReadLinesUntilOk()

	for {
		flag, t := pw.ScanAndText()
		if !flag {
			break
		}

		if strings.HasPrefix(t, "EVENT 24") {
			return false, errors.New("failed to join")
		} else if strings.HasPrefix(t, "EVENT 25") {
			return true, nil
		}
	}

	return false, nil
}

func (b *BP35A1) SKSENDTO(handle uint8, ipaddr string, port uint16, sec uint8, data []byte, done <-chan struct{}) (<-chan *Frame, error) {
	cmd := []byte(fmt.Sprintf("SKSENDTO %X %s %.4X %X %.4X ", handle, ipaddr, port, sec, len(data)))
	cmd = append(cmd[:], data[:]...)
	cmd = append(cmd[:], []byte("\r\n")[:]...)
	err := b.writeBytes(cmd)
	if err != nil {
		return nil, err
	}

	pw := b.getWrappedScanner()

	c := make(chan *Frame)

	go func() {
		defer close(c)

		for {
			flag, t := pw.ScanAndText()
			if !flag {
				return
			}

			if strings.HasPrefix(t, "ERXUDP") {
				frame, err := NewFrameFromString(strings.Split(t, " ")[8])
				if err != nil {
					continue
				}

				select {
				case <-done:
					return
				case c <- frame:
				}
			}
		}
	}()

	return c, nil
}
