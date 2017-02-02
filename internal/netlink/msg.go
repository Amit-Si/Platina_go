// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package netlink

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"syscall"

	"unsafe"

	"github.com/platinasystems/go/internal/accumulate"
	"github.com/platinasystems/go/internal/indent"
)

type Byter interface {
	Bytes() []byte
}

type Message interface {
	netlinkMessage()
	MsgType() MsgType
	io.Closer
	Parse([]byte)
	fmt.Stringer
	TxAdd(*Socket)
	io.WriterTo
}

type multiliner interface {
	multiline()
}

type Runer interface {
	Rune() rune
}

func StringOf(wt io.WriterTo) string {
	buf := pool.Bytes.Get().(*bytes.Buffer)
	defer repool(buf)
	wt.WriteTo(buf)
	return buf.String()
}

type Header struct {
	Len      uint32
	Type     MsgType
	Flags    HeaderFlags
	Sequence uint32
	Pid      uint32
}

const SizeofHeader = 16

func (h *Header) String() string {
	return StringOf(h)
}
func (h *Header) WriteTo(w io.Writer) (int64, error) {
	acc := accumulate.New(w)
	defer acc.Fini()
	fmt.Fprintln(acc, "len:", h.Len)
	fmt.Fprintln(acc, "seq:", h.Sequence)
	fmt.Fprintln(acc, "pid:", h.Pid)
	if h.Flags != 0 {
		fmt.Fprintln(acc, "flags:", h.Flags)
	}
	return acc.Tuple()
}
func (h *Header) MsgType() MsgType { return h.Type }

type GenMessage struct {
	Nsid int
	GenMessageActual
}

type GenMessageActual struct {
	Header
	AddressFamily
}

const SizeofGenMessage = SizeofHeader + SizeofGenmsg
const SizeofGenmsg = SizeofAddressFamily

func NewGenMessage() *GenMessage {
	m := pool.GenMessage.Get().(*GenMessage)
	runtime.SetFinalizer(m, (*GenMessage).Close)
	return m
}

func NewGenMessageBytes(b []byte, nsid int) *GenMessage {
	m := NewGenMessage()
	m.Nsid = nsid
	m.Parse(b)
	return m
}

func (m *GenMessage) netlinkMessage() {}
func (m *GenMessage) Close() error {
	runtime.SetFinalizer(m, nil)
	repool(m)
	return nil
}
func (m *GenMessage) Parse(b []byte) {
	m.GenMessageActual = *(*GenMessageActual)(unsafe.Pointer(&b[0]))
}
func (m *GenMessage) String() string {
	return StringOf(m)
}
func (m *GenMessage) TxAdd(s *Socket) {
	b := s.TxAddReq(&m.Header, SizeofGenmsg)
	p := (*GenMessageActual)(unsafe.Pointer(&b[0]))
	p.AddressFamily = m.AddressFamily
}
func (m *GenMessage) WriteTo(w io.Writer) (int64, error) {
	acc := accumulate.New(w)
	defer acc.Fini()
	fmt.Fprint(acc, m.Header.Type, ":\n")
	indent.Increase(acc)
	if m.Nsid != -1 {
		fmt.Fprintln(acc, "nsid:", m.Nsid)
	}
	m.Header.WriteTo(acc)
	fmt.Fprintln(acc, "family:", m.AddressFamily)
	indent.Decrease(acc)
	return acc.Tuple()
}

type NoopMessage struct {
	Nsid int
	NoopMessageActual
}

type NoopMessageActual struct {
	Header
}

const SizeofNoopMessage = SizeofHeader

func NewNoopMessage() *NoopMessage {
	m := pool.NoopMessage.Get().(*NoopMessage)
	runtime.SetFinalizer(m, (*NoopMessage).Close)
	return m
}

func NewNoopMessageBytes(b []byte, nsid int) *NoopMessage {
	m := NewNoopMessage()
	m.Nsid = nsid
	m.Parse(b)
	return m
}

func (m *NoopMessage) netlinkMessage() {}
func (m *NoopMessage) Close() error {
	runtime.SetFinalizer(m, nil)
	repool(m)
	return nil
}
func (m *NoopMessage) Parse(b []byte) {
	m.NoopMessageActual = *(*NoopMessageActual)(unsafe.Pointer(&b[0]))
}
func (m *NoopMessage) String() string {
	return StringOf(m)
}
func (m *NoopMessage) TxAdd(s *Socket) {
	defer m.Close()
	m.Header.Type = NLMSG_NOOP
	s.TxAddReq(&m.Header, 0)
}
func (m *NoopMessage) WriteTo(w io.Writer) (int64, error) {
	acc := accumulate.New(w)
	defer acc.Fini()
	fmt.Fprint(acc, MessageType(m.Header.Type), ":\n")
	indent.Increase(acc)
	if m.Nsid != -1 {
		fmt.Fprintln(acc, "nsid:", m.Nsid)
	}
	m.Header.WriteTo(acc)
	indent.Decrease(acc)
	return acc.Tuple()
}

type DoneMessage struct {
	Nsid int
	DoneMessageActual
}

type DoneMessageActual struct {
	Header
}

const SizeofDoneMessage = SizeofHeader

func NewDoneMessage() *DoneMessage {
	m := pool.DoneMessage.Get().(*DoneMessage)
	runtime.SetFinalizer(m, (*DoneMessage).Close)
	return m
}

func NewDoneMessageBytes(b []byte, nsid int) *DoneMessage {
	m := NewDoneMessage()
	m.Nsid = nsid
	m.Parse(b)
	return m
}

func (m *DoneMessage) netlinkMessage() {}
func (m *DoneMessage) Close() error {
	runtime.SetFinalizer(m, nil)
	repool(m)
	return nil
}
func (m *DoneMessage) String() string {
	return StringOf(m)
}
func (m *DoneMessage) Parse(b []byte) {
	m.DoneMessageActual = *(*DoneMessageActual)(unsafe.Pointer(&b[0]))
}
func (m *DoneMessage) TxAdd(s *Socket) {
	defer m.Close()
	m.Header.Type = NLMSG_NOOP
	s.TxAddReq(&m.Header, 0)
}
func (m *DoneMessage) WriteTo(w io.Writer) (int64, error) {
	acc := accumulate.New(w)
	defer acc.Fini()
	fmt.Fprint(acc, MessageType(m.Header.Type), ":\n")
	indent.Increase(acc)
	if m.Nsid != -1 {
		fmt.Fprintln(acc, "nsid:", m.Nsid)
	}
	m.Header.WriteTo(acc)
	indent.Decrease(acc)
	return acc.Tuple()
}

type ErrorMessage struct {
	Nsid int
	ErrorMessageActual
}

type ErrorMessageActual struct {
	Header
	// Unix errno for error.
	Errno int32
	// Header for message with error.
	Req Header
}

const SizeofErrorMessage = SizeofHeader + 4 + SizeofHeader

func NewErrorMessage() *ErrorMessage {
	m := pool.ErrorMessage.Get().(*ErrorMessage)
	runtime.SetFinalizer(m, (*ErrorMessage).Close)
	return m
}

func NewErrorMessageBytes(b []byte, nsid int) *ErrorMessage {
	m := NewErrorMessage()
	m.Nsid = nsid
	m.Parse(b)
	return m
}

func (m *ErrorMessage) netlinkMessage() {}
func (m *ErrorMessage) Close() error {
	runtime.SetFinalizer(m, nil)
	repool(m)
	return nil
}
func (m *ErrorMessage) Parse(b []byte) {
	m.ErrorMessageActual = *(*ErrorMessageActual)(unsafe.Pointer(&b[0]))
}
func (m *ErrorMessage) String() string {
	return StringOf(m)
}
func (m *ErrorMessage) TxAdd(s *Socket) {
	defer m.Close()
	m.Header.Type = NLMSG_ERROR
	b := s.TxAddReq(&m.Header, 4+SizeofHeader)
	e := (*ErrorMessageActual)(unsafe.Pointer(&b[0]))
	e.Errno = m.Errno
	e.Req = m.Req
}
func (m *ErrorMessage) WriteTo(w io.Writer) (int64, error) {
	acc := accumulate.New(w)
	defer acc.Fini()
	fmt.Fprint(acc, m.Header.Type, ":\n")
	indent.Increase(acc)
	if m.Nsid != -1 {
		fmt.Fprintln(acc, "nsid:", m.Nsid)
	}
	m.Header.WriteTo(acc)
	fmt.Fprintln(acc, "error:", syscall.Errno(-m.Errno))
	fmt.Fprintln(acc, "req:", m.Req.Type)
	indent.Increase(acc)
	m.Req.WriteTo(acc)
	indent.Decrease(acc)
	indent.Decrease(acc)
	return acc.Tuple()
}

type IfInfoMessage struct {
	Nsid int
	IfInfoMessageActual
}

type IfInfoMessageActual struct {
	Header
	IfInfomsg
	Attrs [IFLA_MAX]Attr
}

const SizeofIfInfoMessage = SizeofHeader + SizeofIfInfomsg

type IfInfomsg struct {
	Family uint8
	_      uint8
	Type   uint16
	Index  uint32
	Flags  IfInfoFlags
	Change IfInfoFlags
}

const SizeofIfInfomsg = 16

func NewIfInfoMessage() *IfInfoMessage {
	m := pool.IfInfoMessage.Get().(*IfInfoMessage)
	runtime.SetFinalizer(m, (*IfInfoMessage).Close)
	return m
}

func NewIfInfoMessageBytes(b []byte, nsid int) *IfInfoMessage {
	m := NewIfInfoMessage()
	m.Nsid = nsid
	m.Parse(b)
	return m
}

func (m *IfInfoMessage) netlinkMessage() {}

func (m *IfInfoMessage) Close() error {
	runtime.SetFinalizer(m, nil)
	closeAttrs(m.Attrs[:])
	repool(m)
	return nil
}

func (m *IfInfoMessage) Parse(b []byte) {
	p := (*IfInfoMessageActual)(unsafe.Pointer(&b[0]))
	m.Header = p.Header
	m.IfInfomsg = p.IfInfomsg
	b = b[SizeofIfInfoMessage:]
	for i := 0; i < len(b); {
		n, v, next_i := nextAttr(b, i)
		i = next_i
		switch t := IfInfoAttrKind(n.Kind()); t {
		case IFLA_IFNAME, IFLA_QDISC:
			m.Attrs[t] = StringAttrBytes(v[:len(v)-1])
		case IFLA_MTU, IFLA_LINK, IFLA_MASTER,
			IFLA_WEIGHT,
			IFLA_NET_NS_PID, IFLA_NET_NS_FD, IFLA_LINK_NETNSID,
			IFLA_EXT_MASK, IFLA_PROMISCUITY,
			IFLA_NUM_TX_QUEUES, IFLA_NUM_RX_QUEUES, IFLA_TXQLEN,
			IFLA_GSO_MAX_SEGS, IFLA_GSO_MAX_SIZE,
			IFLA_CARRIER_CHANGES,
			IFLA_GROUP:
			m.Attrs[t] = Uint32AttrBytes(v)
		case IFLA_CARRIER, IFLA_LINKMODE, IFLA_PROTO_DOWN:
			m.Attrs[t] = Uint8Attr(v[0])
		case IFLA_OPERSTATE:
			m.Attrs[t] = IfOperState(v[0])
		case IFLA_STATS:
			m.Attrs[t] = NewLinkStatsBytes(v)
		case IFLA_STATS64:
			m.Attrs[t] = NewLinkStats64Bytes(v)
		case IFLA_AF_SPEC:
			m.Attrs[t] = parse_af_spec(v, n.kind)
		case IFLA_ADDRESS, IFLA_BROADCAST:
			m.Attrs[t] = afAddr(AF_UNSPEC, v)
		case IFLA_MAP:
		default:
			if t < IFLA_MAX {
				m.Attrs[t] = NewHexStringAttrBytes(v)
			} else {
				panic(fmt.Errorf("%#v: unknown IfInfoMessage attr", t))
			}
		}
	}
}

func (m *IfInfoMessage) String() string {
	return StringOf(m)
}

func (m *IfInfoMessage) TxAdd(s *Socket) {
	defer m.Close()
	as := AttrVec(m.Attrs[:])
	b := s.TxAddReq(&m.Header, SizeofIfInfomsg+as.Size())
	i := (*IfInfoMessageActual)(unsafe.Pointer(&b[0]))
	i.IfInfomsg = m.IfInfomsg
	as.Set(b[SizeofIfInfoMessage:])
}

func (m *IfInfoMessage) WriteTo(w io.Writer) (int64, error) {
	acc := accumulate.New(w)
	defer acc.Fini()
	fmt.Fprint(acc, MessageType(m.Header.Type), ":\n")
	indent.Increase(acc)
	if m.Nsid != -1 {
		fmt.Fprintln(acc, "nsid:", m.Nsid)
	}
	m.Header.WriteTo(acc)
	fmt.Fprintln(acc, "index:", m.Index)
	fmt.Fprintln(acc, "family:", AddressFamily(m.Family))
	fmt.Fprintln(acc, "type:", IfInfoAttrKind(m.Header.Type))
	fmt.Fprintln(acc, "ifinfo flags:", m.IfInfomsg.Flags)
	if m.Change != 0 {
		fmt.Fprintln(acc, "changed flags:", IfInfoFlags(m.Change))
	}
	fprintAttrs(acc, ifInfoAttrKindNames, m.Attrs[:])
	indent.Decrease(acc)
	return acc.Tuple()
}

type IfAddrMessage struct {
	Nsid int
	IfAddrMessageActual
}

type IfAddrMessageActual struct {
	Header
	IfAddrmsg
	Attrs [IFA_MAX]Attr
}

const SizeofIfAddrMessage = SizeofHeader + SizeofIfAddrmsg

type IfAddrmsg struct {
	Family    AddressFamily
	Prefixlen uint8
	Flags     uint8
	Scope     uint8
	Index     uint32
}

const SizeofIfAddrmsg = 8

func NewIfAddrMessage() *IfAddrMessage {
	m := pool.IfAddrMessage.Get().(*IfAddrMessage)
	runtime.SetFinalizer(m, (*IfAddrMessage).Close)
	return m
}

func NewIfAddrMessageBytes(b []byte, nsid int) *IfAddrMessage {
	m := NewIfAddrMessage()
	m.Nsid = nsid
	m.Parse(b)
	return m
}

func (m *IfAddrMessage) netlinkMessage() {}

func (m *IfAddrMessage) Close() error {
	runtime.SetFinalizer(m, nil)
	closeAttrs(m.Attrs[:])
	repool(m)
	return nil
}

func (m *IfAddrMessage) Parse(b []byte) {
	p := (*IfAddrMessageActual)(unsafe.Pointer(&b[0]))
	m.Header = p.Header
	m.IfAddrmsg = p.IfAddrmsg
	b = b[SizeofIfAddrMessage:]
	for i := 0; i < len(b); {
		n, v, next_i := nextAttr(b, i)
		i = next_i
		k := IfAddrAttrKind(n.Kind())
		switch k {
		case IFA_LABEL:
			m.Attrs[k] = StringAttrBytes(v[:len(v)-1])
		case IFA_FLAGS:
			m.Attrs[k] = IfAddrFlagAttrBytes(v)
		case IFA_CACHEINFO:
			m.Attrs[k] = NewIfAddrCacheInfoBytes(v)
		case IFA_ADDRESS, IFA_LOCAL, IFA_BROADCAST, IFA_ANYCAST,
			IFA_MULTICAST:
			m.Attrs[k] = afAddr(AddressFamily(m.Family), v)
		default:
			if k < IFA_MAX {
				m.Attrs[k] = NewHexStringAttrBytes(v)
			} else {
				panic(fmt.Errorf("%#v: unknown IfAddrMessage attr", k))
			}
		}
	}
	return
}

func (m *IfAddrMessage) String() string {
	return StringOf(m)
}

func (m *IfAddrMessage) TxAdd(s *Socket) {
	defer m.Close()
	as := AttrVec(m.Attrs[:])
	b := s.TxAddReq(&m.Header, SizeofIfAddrmsg+as.Size())
	i := (*IfAddrMessageActual)(unsafe.Pointer(&b[0]))
	i.IfAddrmsg = m.IfAddrmsg
	as.Set(b[SizeofIfAddrMessage:])
}

func (m *IfAddrMessage) WriteTo(w io.Writer) (int64, error) {
	acc := accumulate.New(w)
	defer acc.Fini()
	fmt.Fprint(acc, MessageType(m.Header.Type), ":\n")
	indent.Increase(acc)
	if m.Nsid != -1 {
		fmt.Fprintln(acc, "nsid:", m.Nsid)
	}
	m.Header.WriteTo(acc)
	fmt.Fprintln(acc, "index:", m.Index)
	fmt.Fprintln(acc, "family:", AddressFamily(m.Family))
	fmt.Fprintln(acc, "prefix:", m.Prefixlen)
	fmt.Fprintln(acc, "ifaddr flags:", m.IfAddrmsg.Flags)
	fmt.Fprintln(acc, "scope:", RtScope(m.Scope))
	fprintAttrs(acc, ifAddrAttrKindNames, m.Attrs[:])
	indent.Decrease(acc)
	return acc.Tuple()
}

type RouteMessage struct {
	Nsid int
	RouteMessageActual
}

type RouteMessageActual struct {
	Header
	Rtmsg
	Attrs [RTA_MAX]Attr
}

const SizeofRouteMessage = SizeofHeader + SizeofRtmsg

type Rtmsg struct {
	Family     AddressFamily
	DstLen     uint8
	SrcLen     uint8
	Tos        uint8
	Table      RouteTableKind
	Protocol   RouteProtocol
	Scope      RtScope
	RouteType  RouteType
	RouteFlags RouteFlags
}

const SizeofRtmsg = 12

func NewRouteMessage() *RouteMessage {
	m := pool.RouteMessage.Get().(*RouteMessage)
	runtime.SetFinalizer(m, (*RouteMessage).Close)
	return m
}

func NewRouteMessageBytes(b []byte, nsid int) *RouteMessage {
	m := NewRouteMessage()
	m.Nsid = nsid
	m.Parse(b)
	return m
}

func (m *RouteMessage) netlinkMessage() {}

func (m *RouteMessage) Close() error {
	runtime.SetFinalizer(m, nil)
	closeAttrs(m.Attrs[:])
	repool(m)
	return nil
}

func (m *RouteMessage) Parse(b []byte) {
	p := (*RouteMessageActual)(unsafe.Pointer(&b[0]))
	m.Header = p.Header
	m.Rtmsg = p.Rtmsg
	b = b[SizeofRouteMessage:]
	for i := 0; i < len(b); {
		n, v, next_i := nextAttr(b, i)
		i = next_i
		k := RouteAttrKind(n.Kind())
		switch k {
		case RTA_DST, RTA_SRC, RTA_PREFSRC, RTA_GATEWAY:
			m.Attrs[k] = afAddr(AddressFamily(m.Family), v)
		case RTA_TABLE, RTA_IIF, RTA_OIF, RTA_PRIORITY, RTA_FLOW:
			m.Attrs[k] = Uint32AttrBytes(v)
		case RTA_ENCAP_TYPE:
			m.Attrs[k] = Uint16AttrBytes(v)
		case RTA_PREF:
			m.Attrs[k] = Uint8Attr(v[0])
		case RTA_CACHEINFO:
			m.Attrs[k] = NewRtaCacheInfoBytes(v)
		default:
			if k < RTA_MAX {
				m.Attrs[k] = NewHexStringAttrBytes(v)
			} else {
				panic(fmt.Errorf("%#v: unknown RouteMessage attr", k))
			}
		}
	}
	return
}

func (m *RouteMessage) String() string {
	return StringOf(m)
}

func (m *RouteMessage) TxAdd(s *Socket) {
	defer m.Close()
	as := AttrVec(m.Attrs[:])
	b := s.TxAddReq(&m.Header, SizeofRtmsg+as.Size())
	i := (*RouteMessageActual)(unsafe.Pointer(&b[0]))
	i.Rtmsg = m.Rtmsg
	as.Set(b[SizeofRouteMessage:])
}

func (m *RouteMessage) WriteTo(w io.Writer) (int64, error) {
	acc := accumulate.New(w)
	defer acc.Fini()
	fmt.Fprint(acc, MessageType(m.Header.Type), ":\n")
	indent.Increase(acc)
	if m.Nsid != -1 {
		fmt.Fprintln(acc, "nsid:", m.Nsid)
	}
	m.Header.WriteTo(acc)
	fmt.Fprintln(acc, "family:", AddressFamily(m.Family))
	fmt.Fprintln(acc, "srclen:", m.SrcLen)
	fmt.Fprintln(acc, "dstlen:", m.DstLen)
	fmt.Fprintln(acc, "tos:", m.Tos)
	fmt.Fprintln(acc, "table:", m.Table)
	fmt.Fprintln(acc, "protocol:", m.Protocol)
	fmt.Fprintln(acc, "scope:", m.Scope)
	fmt.Fprintln(acc, "type:", m.RouteType)
	if m.RouteFlags != 0 {
		fmt.Fprintln(acc, "route flags:", m.RouteFlags)
	}
	fprintAttrs(acc, routeAttrKindNames, m.Attrs[:])
	indent.Decrease(acc)
	return acc.Tuple()
}

type NeighborMessage struct {
	Nsid int
	NeighborMessageActual
}

type NeighborMessageActual struct {
	Header
	Ndmsg
	Attrs [NDA_MAX]Attr
}

const SizeofNeighborMessage = SizeofHeader + SizeofNdmsg

type Ndmsg struct {
	Family AddressFamily
	_      [3]uint8
	Index  uint32
	State  NeighborState
	Flags  uint8
	Type   RouteType
}

const SizeofNdmsg = 12

func NewNeighborMessage() *NeighborMessage {
	m := pool.NeighborMessage.Get().(*NeighborMessage)
	runtime.SetFinalizer(m, (*NeighborMessage).Close)
	return m
}

func NewNeighborMessageBytes(b []byte, nsid int) *NeighborMessage {
	m := NewNeighborMessage()
	m.Nsid = nsid
	m.Parse(b)
	return m
}

func (m *NeighborMessage) netlinkMessage() {}

func (m *NeighborMessage) AttrBytes(kind NeighborAttrKind) []byte {
	return m.Attrs[kind].(Byter).Bytes()
}

func (m *NeighborMessage) Close() error {
	runtime.SetFinalizer(m, nil)
	closeAttrs(m.Attrs[:])
	repool(m)
	return nil
}

func (m *NeighborMessage) Parse(b []byte) {
	p := (*NeighborMessageActual)(unsafe.Pointer(&b[0]))
	m.Header = p.Header
	m.Ndmsg = p.Ndmsg
	b = b[SizeofNeighborMessage:]
	for i := 0; i < len(b); {
		n, v, next_i := nextAttr(b, i)
		i = next_i
		k := NeighborAttrKind(n.Kind())
		switch k {
		case NDA_DST:
			m.Attrs[k] = afAddr(AddressFamily(m.Family), v)
		case NDA_LLADDR:
			m.Attrs[k] = afAddr(AF_UNSPEC, v)
		case NDA_CACHEINFO:
			m.Attrs[k] = NewNdaCacheInfoBytes(v)
		case NDA_PROBES, NDA_VNI, NDA_IFINDEX, NDA_MASTER,
			NDA_LINK_NETNSID:
			m.Attrs[k] = Uint32AttrBytes(v)
		case NDA_VLAN:
			m.Attrs[k] = Uint16AttrBytes(v)
		default:
			if k < NDA_MAX {
				m.Attrs[k] = NewHexStringAttrBytes(v)
			} else {
				panic(fmt.Errorf("%#v: unknown NeighborMessage attr", k))
			}
		}
	}
	return
}

func (m *NeighborMessage) String() string {
	return StringOf(m)
}

func (m *NeighborMessage) TxAdd(s *Socket) {
	defer m.Close()
	as := AttrVec(m.Attrs[:])
	b := s.TxAddReq(&m.Header, SizeofNdmsg+as.Size())
	i := (*NeighborMessageActual)(unsafe.Pointer(&b[0]))
	i.Ndmsg = m.Ndmsg
	as.Set(b[SizeofNeighborMessage:])
}

func (m *NeighborMessage) WriteTo(w io.Writer) (int64, error) {
	acc := accumulate.New(w)
	defer acc.Fini()
	fmt.Fprint(acc, MessageType(m.Header.Type), ":\n")
	indent.Increase(acc)
	if m.Nsid != -1 {
		fmt.Fprintln(acc, "nsid:", m.Nsid)
	}
	m.Header.WriteTo(acc)
	fmt.Fprintln(acc, "index:", m.Index)
	fmt.Fprintln(acc, "address family:", AddressFamily(m.Family))
	fmt.Fprintln(acc, "type:", RouteType(m.Ndmsg.Type))
	fmt.Fprintln(acc, "state:", NeighborState(m.State))
	if m.Ndmsg.Flags != 0 {
		fmt.Fprintln(acc, "neighbor flags:", m.Ndmsg.Flags)
	}
	fprintAttrs(acc, neighborAttrKindNames, m.Attrs[:])
	indent.Decrease(acc)
	return acc.Tuple()
}

type NetnsMessage struct {
	Header
	AddressFamily
	Attrs [NETNSA_MAX]Attr
}

const NetnsPad = 3
const SizeofNetnsmsg = SizeofGenmsg + NetnsPad
const SizeofNetnsMessage = SizeofGenMessage + NetnsPad

func NewNetnsMessage() *NetnsMessage {
	m := pool.NetnsMessage.Get().(*NetnsMessage)
	runtime.SetFinalizer(m, (*NetnsMessage).Close)
	return m
}

func NewNetnsMessageBytes(b []byte) *NetnsMessage {
	m := NewNetnsMessage()
	m.Parse(b)
	return m
}

func (m *NetnsMessage) NSID() int32 {
	return m.Attrs[NETNSA_NSID].(Int32Attr).Int()
}

func (m *NetnsMessage) PID() uint32 {
	return m.Attrs[NETNSA_PID].(Uint32Attr).Uint()
}

func (m *NetnsMessage) FD() uint32 {
	return m.Attrs[NETNSA_FD].(Uint32Attr).Uint()
}

func (m *NetnsMessage) netlinkMessage() {}

func (m *NetnsMessage) AttrBytes(kind NetnsAttrKind) []byte {
	return m.Attrs[kind].(Byter).Bytes()
}

func (m *NetnsMessage) Close() error {
	runtime.SetFinalizer(m, nil)
	closeAttrs(m.Attrs[:])
	repool(m)
	return nil
}

func (m *NetnsMessage) Parse(b []byte) {
	p := (*NetnsMessage)(unsafe.Pointer(&b[0]))
	m.Header = p.Header
	m.AddressFamily = p.AddressFamily
	b = b[SizeofNetnsMessage:]
	m.Attrs[NETNSA_NSID] = Int32Attr(-2)
	m.Attrs[NETNSA_PID] = Uint32Attr(0)
	m.Attrs[NETNSA_FD] = Uint32Attr(^uint32(0))
	for i := 0; i < len(b); {
		n, v, next_i := nextAttr(b, i)
		i = next_i
		k := NetnsAttrKind(n.Kind())
		switch k {
		case NETNSA_NONE:
		case NETNSA_NSID:
			m.Attrs[k] = Int32AttrBytes(v)
		case NETNSA_PID, NETNSA_FD:
			m.Attrs[k] = Uint32AttrBytes(v)
		default:
			panic(fmt.Errorf("%#v: unknown NetnsMessage attr", k))
		}
	}
	return
}

func (m *NetnsMessage) String() string {
	return StringOf(m)
}

func (m *NetnsMessage) TxAdd(s *Socket) {
	defer m.Close()
	as := AttrVec(m.Attrs[:])
	b := s.TxAddReq(&m.Header, SizeofNetnsmsg+as.Size())
	b[SizeofHeader] = byte(m.AddressFamily)
	as.Set(b[SizeofNetnsMessage:])
}

func (m *NetnsMessage) WriteTo(w io.Writer) (int64, error) {
	acc := accumulate.New(w)
	defer acc.Fini()
	fmt.Fprint(acc, m.Header.Type, ":\n")
	indent.Increase(acc)
	m.Header.WriteTo(acc)
	fmt.Fprintln(acc, "family:", m.AddressFamily)
	fprintAttrs(acc, netnsAttrKindNames, m.Attrs[:])
	indent.Decrease(acc)
	return acc.Tuple()
}
