package gosshtun

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import "github.com/tinylib/msgp/msgp"

// DecodeMsg implements msgp.Decodable
func (z *AtomicUserMap) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zbai uint32
	zbai, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zbai > 0 {
		zbai--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "U":
			var zcmr uint32
			zcmr, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.U == nil && zcmr > 0 {
				z.U = make(map[string]*User, zcmr)
			} else if len(z.U) > 0 {
				for key, _ := range z.U {
					delete(z.U, key)
				}
			}
			for zcmr > 0 {
				zcmr--
				var zxvk string
				var zbzg *User
				zxvk, err = dc.ReadString()
				if err != nil {
					return
				}
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						return
					}
					zbzg = nil
				} else {
					if zbzg == nil {
						zbzg = new(User)
					}
					err = zbzg.DecodeMsg(dc)
					if err != nil {
						return
					}
				}
				z.U[zxvk] = zbzg
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *AtomicUserMap) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "U"
	err = en.Append(0x81, 0xa1, 0x55)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.U)))
	if err != nil {
		return
	}
	for zxvk, zbzg := range z.U {
		err = en.WriteString(zxvk)
		if err != nil {
			return
		}
		if zbzg == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = zbzg.EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *AtomicUserMap) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "U"
	o = append(o, 0x81, 0xa1, 0x55)
	o = msgp.AppendMapHeader(o, uint32(len(z.U)))
	for zxvk, zbzg := range z.U {
		o = msgp.AppendString(o, zxvk)
		if zbzg == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = zbzg.MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *AtomicUserMap) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zajw uint32
	zajw, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zajw > 0 {
		zajw--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "U":
			var zwht uint32
			zwht, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.U == nil && zwht > 0 {
				z.U = make(map[string]*User, zwht)
			} else if len(z.U) > 0 {
				for key, _ := range z.U {
					delete(z.U, key)
				}
			}
			for zwht > 0 {
				var zxvk string
				var zbzg *User
				zwht--
				zxvk, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					zbzg = nil
				} else {
					if zbzg == nil {
						zbzg = new(User)
					}
					bts, err = zbzg.UnmarshalMsg(bts)
					if err != nil {
						return
					}
				}
				z.U[zxvk] = zbzg
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *AtomicUserMap) Msgsize() (s int) {
	s = 1 + 2 + msgp.MapHeaderSize
	if z.U != nil {
		for zxvk, zbzg := range z.U {
			_ = zbzg
			s += msgp.StringPrefixSize + len(zxvk)
			if zbzg == nil {
				s += msgp.NilSize
			} else {
				s += zbzg.Msgsize()
			}
		}
	}
	return
}
