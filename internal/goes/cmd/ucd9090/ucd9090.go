// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

// Package ucd9090 provides access to the UCD9090 Power Sequencer/Monitor chip
package ucd9090

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/platinasystems/go/internal/goes"
	"github.com/platinasystems/go/internal/redis/publisher"
)

const Name = "ucd9090"

type I2cDev struct {
	Bus      int
	Addr     int
	MuxBus   int
	MuxAddr  int
	MuxValue int
}

var (
	Hook = func() {}
	once sync.Once

	Vdev I2cDev

	VpageByKey map[string]uint8

	loggedFaultCount      uint8
	lastLoggedFaultDetail [12]byte

	first int
)

type cmd struct {
	stop  chan struct{}
	pub   *publisher.Publisher
	last  map[string]float64
	lasts map[string]string
	lastu map[string]uint16
}

func New() *cmd { return new(cmd) }

func (*cmd) Kind() goes.Kind { return goes.Daemon }
func (*cmd) String() string  { return Name }
func (*cmd) Usage() string   { return Name }

func (cmd *cmd) Main(...string) error {
	once.Do(Hook)

	var si syscall.Sysinfo_t
	var err error
	first = 1

	cmd.stop = make(chan struct{})
	cmd.last = make(map[string]float64)
	cmd.lasts = make(map[string]string)
	cmd.lastu = make(map[string]uint16)
	if cmd.pub, err = publisher.New(); err != nil {
		return err
	}

	if err = syscall.Sysinfo(&si); err != nil {
		return err
	}

	//if err = cmd.update(); err != nil {
	//	close(cmd.stop)
	//	return err
	//}
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-cmd.stop:
			return nil
		case <-t.C:
			if Vdev.Addr != 0 {
				if err = cmd.update(); err != nil {
					close(cmd.stop)
					return err
				}
			}
		}
	}
	return nil
}

func (cmd *cmd) Close() error {
	close(cmd.stop)
	return nil
}

func (cmd *cmd) update() error {
	stopped := readStopped()
	if stopped == 1 {
		return nil
	}

	if first == 1 {
		Vdev.ucdInit()
		first = 0
	}

	for k, i := range VpageByKey {
		if strings.Contains(k, "units.V") {
			v, err := Vdev.Vout(i)
			if err != nil {
				return err
			}
			if v != cmd.last[k] {
				cmd.pub.Print(k, ": ", v)
				cmd.last[k] = v
			}
		}
		if strings.Contains(k, "poweroff.events") {
			v, err := Vdev.PowerCycles()
			if err != nil {
				return err
			}
			if (v != "") && (v != cmd.lasts[k]) {
				cmd.pub.Print(k, ": ", v)
				cmd.lasts[k] = v
			}
		}
	}
	return nil
}

func (h *I2cDev) ucdInit() error {
	//FIXME configure UCD run time clock, pending i2c block write
	//now := time.Now()
	//nanos := now.UnixNano()
	//days := nanos / int64(math.Pow(10, 9)) / 60 / 60 / 24
	//millisecs := (nanos - days*60*60*24*int64(math.Pow(10, 9))) / int64(math.Pow(10, 6))
	return nil
}

func (h *I2cDev) Vout(i uint8) (float64, error) {
	if i > 10 {
		panic("Voltage rail subscript out of range\n")
	}
	i--

	r := getRegs()
	r.Page.set(h, i)
	r.VoutMode.get(h)
	r.ReadVout.get(h)
	closeMux(h)
	err := DoI2cRpc()
	if err != nil {
		return 0, err
	}
	n := s[3].D[0] & 0xf
	n--
	n = (n ^ 0xf) & 0xf
	v := uint16(s[5].D[1])<<8 | uint16(s[5].D[0])

	nn := float64(n) * (-1)
	vv := float64(v) * (math.Exp2(nn))
	vv, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", vv), 64)
	return float64(vv), nil
}

/* FIXME fucntions pending i2c block write
func (h *I2cDev) LoggedFaults() error {
	r := getRegs()

	//Print Logged Faults
	r.LoggedFaults.get(h, 13)
	closeMux(h)
	err := DoI2cRpc()
	if err != nil {
		return err
	}
	log.Printf("logged faults: 0x%x", s[1].D)
	return nil
}

func (h *I2cDev) ClearLoggedFaults() error {
	r := getRegs()
	data := []byte{12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	r.LoggedFaults.set(h, data)
	closeMux(h)
	err := DoI2cRpc()
	if err != nil {
		return err
	}
	return nil
}
func (h *I2cDev) ConfigLoggedFaults(i int) error {
	r := getRegs()
	//read misc config
	if i == 0 {
		r.MiscConfig.get(h, 3)
		closeMux(h)
		err := DoI2cRpc()
		if err != nil {
			return err
		}
		log.Printf("misc config: 0x%x", s[1].D)
		return nil
	} else if i == 1 {
		data := []byte{2, 7, 0}
		r.MiscConfig.set(h, data)
		closeMux(h)
		err := DoI2cRpc()
		if err != nil {
			return err
		}
		return nil
	} else if i == 2 {
		data := []byte{2, 0, 0}
		r.MiscConfig.set(h, data)
		closeMux(h)
		err := DoI2cRpc()
		if err != nil {
			return err
		}
		return nil
	} else if i == 3 {
		r.StoreDefaultAll.set(h, 0x2a)
		closeMux(h)
		err := DoI2cRpc()
		if err != nil {
			return err
		}
	}
	return nil
}
*/
func (h *I2cDev) PowerCycles() (string, error) {
	r := getRegs()
	r.LoggedFaultIndex.get(h)
	closeMux(h)
	err := DoI2cRpc()
	if err != nil {
		return "", err
	}

	d := s[1].D[1]

	var milli uint32
	var seconds uint32
	var faultType uint8
	var pwrCycles string

	for i := 0; i < int(d); i++ {
		r.LoggedFaultIndex.set(h, uint16(i)<<8)
		err := DoI2cRpc()
		if err != nil {
			return "", err
		}
		r.LoggedFaultDetail.get(h, 11)
		err = DoI2cRpc()
		if err != nil {
			return "", err
		}

		if i == 0 {
			new := false
			if loggedFaultCount != d {
				loggedFaultCount = d
				copy(lastLoggedFaultDetail[:], s[1].D[0:12])
				new = true
			} else {
				for j := 0; j < 12; j++ {
					if s[1].D[j] != lastLoggedFaultDetail[j] {
						copy(lastLoggedFaultDetail[:], s[1].D[0:12])
						new = true
						break
					}
				}
			}
			if !new {
				return "", nil
			}
		}
		milli = uint32(s[1].D[5]) + uint32(s[1].D[4])<<8 + uint32(s[1].D[3])<<16 + uint32(s[1].D[2])<<24
		seconds = milli / 1000
		timestamp := time.Unix(int64(seconds), 0).Format(time.RFC3339)

		faultType = (s[1].D[6] >> 3) & 0xF

		if !strings.Contains(pwrCycles, timestamp) && (faultType == 0 || faultType == 1) {
			pwrCycles += timestamp + "."
		}
	}
	pwrCycles = strings.Trim(pwrCycles, ".")
	return pwrCycles, nil
}

func (h *I2cDev) LoggedFaultDetail() (string, error) {
	r := getRegs()
	r.LoggedFaultIndex.get(h)
	closeMux(h)
	err := DoI2cRpc()
	if err != nil {
		return "", err
	}

	d := s[1].D[1]

	var milli uint32
	var page uint8
	var seconds uint32
	var faultType uint8
	var paged uint8
	var rail string
	var fault string
	var log string

	for i := 0; i < int(d); i++ {
		r.LoggedFaultIndex.set(h, uint16(i)<<8)
		err := DoI2cRpc()
		if err != nil {
			return "", err
		}
		r.LoggedFaultDetail.get(h, 11)
		err = DoI2cRpc()
		if err != nil {
			return "", err
		}

		if i == 0 {
			new := false
			if loggedFaultCount != d {
				loggedFaultCount = d
				copy(lastLoggedFaultDetail[:], s[1].D[0:12])
				new = true
			} else {
				for j := 0; j < 12; j++ {
					if s[1].D[j] != lastLoggedFaultDetail[j] {
						copy(lastLoggedFaultDetail[:], s[1].D[0:12])
						new = true
						break
					}
				}
			}
			if !new {
				return "", nil
			}
		}
		milli = uint32(s[1].D[5]) + uint32(s[1].D[4])<<8 + uint32(s[1].D[3])<<16 + uint32(s[1].D[2])<<24
		seconds = milli / 1000
		timestamp := time.Unix(int64(seconds), 0).Format(time.RFC3339)

		faultType = (s[1].D[6] >> 3) & 0xF
		paged = s[1].D[6] & 0x80 >> 7
		page = ((s[1].D[7] & 0x80) >> 7) + ((s[1].D[6] & 0x7) << 1)

		if paged == 1 {
			switch page {
			case 0:
				rail = "P5V_SB"
			case 1:
				rail = "P3V8_BMC"
			case 2:
				rail = "P3V3_SB"
			case 3:
				rail = "PERI_3V3"
			case 4:
				rail = "P3V3"
			case 5:
				rail = "VDD_CORE"
			case 6:
				rail = "P1V8"
			case 7:
				rail = "P1V25"
			case 8:
				rail = "P1V2"
			case 9:
				rail = "P1V0"
			default:
				rail = "n/a"
			}
			switch faultType {
			case 0:
				fault = "VOUT_OV"
			case 1:
				fault = "VOUT_UV"
			case 2:
				fault = "TON_MAX"
			case 3:
				fault = "IOUT_OC"
			case 4:
				fault = "IOUT_UC"
			case 5:
				fault = "TEMPERATURE_OT"
			case 6:
				fault = "SEQUENCE ON TIMEOUT"
			case 7:
				fault = "SEQUENCE OFF TIMEOUT"
			default:
				fault = "unknown"
			}
		} else {
			rail = "n/a"
			switch faultType {
			case 1:
				fault = "SYSTEM WATCHDOG TIMEOUT"
			case 2:
				fault = "RESEQUENCE ERROR"
			case 3:
				fault = "WATCHDOG TIMEOUT"
			case 8:
				fault = "FAN FAULT"
			case 9:
				fault = "GPI FAULT"
			default:
				fault = "unknown"
			}

		}
		log += timestamp + "." + rail + "." + fault + "\n"
	}
	return log, nil
}
