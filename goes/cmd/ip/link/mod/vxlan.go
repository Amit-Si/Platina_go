// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package mod

import (
	"fmt"
	"net"
	"strings"
	"syscall"

	"github.com/platinasystems/go/goes/cmd/ip/internal/rtnl"
)

// ip link COMMAND type vxlan id VNI
//	[ { group | remote } ADDR ]
//	[ local ADDR ]
//	[ dev PHYS_DEV ]
//	[ ttl TTL ]
//	[ tos TOS ]
//	[ flowlabel LABEL ]
//	[ ageing SECONDS ]
//	[ maxaddress NUMBER ]
//	[ dstport PORT ]
//	[ srcport LOW:HIGH ]
//	[ [no-]learning ]
//	[ [no-]proxy ]
//	[ [no-]rsc ]
//	[ [no-]l2miss ]
//	[ [no-]l3miss ]
//	[ [no-]udpcsum ]
//	[ [no-]udp6zerocsumtx ]
//	[ [no-]udp6zerocsumrx ]
//	[ [no-]remcsumtx ]
//	[ [no-]external ]
//	[ gbp ]
//	[ gpe ]
func (c *Command) parseTypeVxlan() error {
	var gaddr, laddr, raddr net.IP
	var s string
	var u8 uint8
	var u16 uint16
	var u32 uint32
	var err error
	c.args = c.opt.Parms.More(c.args, []string{"id", "vni"})
	s = c.opt.Parms.ByName["id"]
	if len(s) == 0 {
		return fmt.Errorf("id: missing")
	}
	if _, err = fmt.Sscan(s, &u32); err != nil {
		return fmt.Errorf("vni: %q %v", s, err)
	} else if u32 >= 1<<24 {
		return fmt.Errorf("vni: %q %v", s, syscall.ERANGE)
	}
	c.tinfo = append(c.tinfo,
		rtnl.Attr{rtnl.IFLA_VXLAN_ID, rtnl.Uint32Attr(u32)})
	c.args = c.opt.Parms.More(c.args, "group", "local", "remote")
	for _, x := range []struct {
		name string
		p    *net.IP
	}{
		{"group", &gaddr},
		{"local", &laddr},
		{"remote", &raddr},
	} {
		s = c.opt.Parms.ByName[x.name]
		if len(s) == 0 {
			continue
		}
		*x.p = net.ParseIP(s)
		if *x.p == nil {
			return fmt.Errorf("%s: %q invalid", x.name, s)
		}
	}
	for _, addr := range []net.IP{gaddr, raddr} {
		if addr != nil {
			if ip4 := addr.To4(); ip4 != nil {
				c.tinfo = append(c.tinfo,
					rtnl.Attr{rtnl.IFLA_VXLAN_GROUP,
						rtnl.BytesAttr(ip4)})
			} else {
				c.tinfo = append(c.tinfo,
					rtnl.Attr{rtnl.IFLA_VXLAN_GROUP6,
						rtnl.BytesAttr(addr.To16())})
			}
			break
		}
	}
	if ip4 := laddr.To4(); ip4 != nil {
		c.tinfo = append(c.tinfo,
			rtnl.Attr{rtnl.IFLA_VXLAN_LOCAL, rtnl.BytesAttr(ip4)})
	} else {
		c.tinfo = append(c.tinfo,
			rtnl.Attr{rtnl.IFLA_VXLAN_LOCAL6,
				rtnl.BytesAttr(laddr.To16())})
	}
	c.args = c.opt.Parms.More(c.args, "dev")
	if s = c.opt.Parms.ByName["dev"]; len(s) > 0 {
		if dev, found := c.ifindexByName[s]; !found {
			return fmt.Errorf("dev: %q not found", s)
		} else {
			c.tinfo = append(c.tinfo,
				rtnl.Attr{rtnl.IFLA_VXLAN_LINK,
					rtnl.Uint32Attr(dev)})
		}
	}
	for _, x := range []struct {
		names []string
		t     uint16
	}{
		{[]string{"ttl", "hoplimit"}, rtnl.IFLA_VXLAN_TTL},
		{[]string{"tos", "dsfield"}, rtnl.IFLA_VXLAN_TOS},
	} {
		c.args = c.opt.Parms.More(c.args, x.names)
		s = c.opt.Parms.ByName[x.names[0]]
		if len(s) == 0 || s == "inherit" {
			continue
		}
		if _, err = fmt.Sscan(s, &u8); err != nil {
			return fmt.Errorf("%s: %q %v", x.names[0], s, err)
		}
		c.tinfo = append(c.tinfo, rtnl.Attr{x.t, rtnl.Uint8Attr(u8)})
	}
	c.args = c.opt.Parms.More(c.args, "flowlabel")
	if s = c.opt.Parms.ByName["flowlabel"]; len(s) > 0 {
		var u32 uint32
		if _, err = fmt.Sscan(s, &u32); err != nil {
			return fmt.Errorf("flowlabel: %q %v", s, err)
		}
		c.tinfo = append(c.tinfo, rtnl.Attr{rtnl.IFLA_VXLAN_LABEL,
			rtnl.Be32Attr(u32)})
	}
	for _, x := range []struct {
		names []string
		t     uint16
	}{
		{[]string{"ageing"}, rtnl.IFLA_VXLAN_AGEING},
		{[]string{"maxaddress"}, rtnl.IFLA_VXLAN_AGEING},
	} {
		var u32 uint32
		c.args = c.opt.Parms.More(c.args, x.names)
		s = c.opt.Parms.ByName[x.names[0]]
		if len(s) == 0 {
			continue
		}
		if _, err = fmt.Sscan(s, &u32); err != nil {
			return fmt.Errorf("%s: %q %v", x.names[0], s, err)
		}
		c.tinfo = append(c.tinfo, rtnl.Attr{x.t, rtnl.Uint32Attr(u32)})
	}
	c.args = c.opt.Parms.More(c.args, "dstport")
	if s = c.opt.Parms.ByName["dstport"]; len(s) > 0 {
		if _, err = fmt.Sscan(s, &u16); err != nil {
			return fmt.Errorf("dstport: %q %v", s, err)
		}
		c.tinfo = append(c.tinfo, rtnl.Attr{rtnl.IFLA_VXLAN_PORT,
			rtnl.Be16Attr(u16)})
	}
	c.args = c.opt.Parms.More(c.args, []string{"srcport", "port"})
	if s = c.opt.Parms.ByName["srcport"]; len(s) > 0 {
		var pr rtnl.IflaVxlanPortRange
		colon := strings.Index(s, ":")
		if colon < 1 {
			return fmt.Errorf("srcport: %q invalid", s)
		}
		if _, err = fmt.Sscan(s[:colon], &pr.Low); err != nil {
			return fmt.Errorf("srcport low: %q %v", s, err)
		}
		if _, err = fmt.Sscan(s[colon+1:], &pr.High); err != nil {
			return fmt.Errorf("srcport high: %q %v", s, err)
		}
		c.tinfo = append(c.tinfo, rtnl.Attr{rtnl.IFLA_VXLAN_PORT_RANGE,
			pr})
	}
	for _, x := range []struct {
		set   []string
		unset []string
		t     uint16
	}{
		{
			[]string{"learning", "+learning"},
			[]string{"no-learning", "-learning"},
			rtnl.IFLA_VXLAN_LEARNING,
		},
		{
			[]string{"proxy", "+proxy"},
			[]string{"no-proxy", "-proxy"},
			rtnl.IFLA_VXLAN_PROXY,
		},
		{
			[]string{"rsc", "+rsc"},
			[]string{"no-rsc", "-rsc"},
			rtnl.IFLA_VXLAN_RSC,
		},
		{
			[]string{"l2miss", "+l2miss"},
			[]string{"no-l2miss", "-l2miss"},
			rtnl.IFLA_VXLAN_L2MISS,
		},
		{
			[]string{"l3miss", "+l3miss"},
			[]string{"no-l3miss", "-l3miss"},
			rtnl.IFLA_VXLAN_L2MISS,
		},
		{
			[]string{"udpcsum", "+udpcsum"},
			[]string{"no-udpcsum", "-udpcsum"},
			rtnl.IFLA_VXLAN_UDP_CSUM,
		},
		{
			[]string{"udp6zerocsumtx", "+udp6zerocsumtx"},
			[]string{"no-udp6zerocsumtx", "-udp6zerocsumtx"},
			rtnl.IFLA_VXLAN_UDP_ZERO_CSUM6_TX,
		},
		{
			[]string{"udp6zerocsumrx", "+udp6zerocsumrx"},
			[]string{"no-udp6zerocsumrx", "-udp6zerocsumrx"},
			rtnl.IFLA_VXLAN_UDP_ZERO_CSUM6_RX,
		},
		{
			[]string{"remcsumtx", "+remcsumtx"},
			[]string{"no-remcsumtx", "-remcsumtx"},
			rtnl.IFLA_VXLAN_REMCSUM_TX,
		},
		{
			[]string{"remcsumrx", "+remcsumrx"},
			[]string{"no-remcsumrx", "-remcsumrx"},
			rtnl.IFLA_VXLAN_REMCSUM_TX,
		},
	} {
		c.args = c.opt.Flags.More(c.args, x.set, x.unset)
		if c.opt.Flags.ByName[x.set[0]] {
			c.tinfo = append(c.tinfo, rtnl.Attr{x.t,
				rtnl.Uint8Attr(1)})
		} else if c.opt.Flags.ByName[x.set[0]] {
			c.tinfo = append(c.tinfo, rtnl.Attr{x.t,
				rtnl.Uint8Attr(0)})
		}
	}
	c.args = c.opt.Flags.More(c.args,
		[]string{"external", "+external"},
		[]string{"no-external", "-external"},
	)
	if c.opt.Flags.ByName["external"] {
		c.tinfo = append(c.tinfo,
			rtnl.Attr{rtnl.IFLA_VXLAN_COLLECT_METADATA,
				rtnl.Uint8Attr(1)})
		c.tinfo = append(c.tinfo, rtnl.Attr{rtnl.IFLA_VXLAN_LEARNING,
			rtnl.Uint8Attr(0)})
	} else if c.opt.Flags.ByName["no-external"] {
		c.tinfo = append(c.tinfo,
			rtnl.Attr{rtnl.IFLA_VXLAN_COLLECT_METADATA,
				rtnl.Uint8Attr(0)})
	}
	for _, x := range []struct {
		set []string
		t   uint16
	}{
		{
			[]string{"gbp", "+gbp"},
			rtnl.IFLA_VXLAN_GBP,
		},
		{
			[]string{"gbp", "+gbp"},
			rtnl.IFLA_VXLAN_GPE,
		},
	} {
		c.args = c.opt.Flags.More(c.args, x.set)
		if c.opt.Flags.ByName[x.set[0]] {
			c.tinfo = append(c.tinfo,
				rtnl.Attr{x.t, rtnl.NilAttr{}})
		}
	}
	return nil
}
