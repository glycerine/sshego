package sshego

// NOTE: THIS FILE WAS PRODUCED BY THE
// GREENPACK CODE GENERATION TOOL (github.com/glycerine/greenpack)
// DO NOT EDIT

import (
	"github.com/glycerine/greenpack/msgp"
)

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *AtomicUserMap) DecodeMsg(dc *msgp.Reader) (err error) {
	var sawTopNil bool
	if dc.IsNil() {
		sawTopNil = true
		err = dc.ReadNil()
		if err != nil {
			return
		}
		dc.PushAlwaysNil()
	}

	var field []byte
	_ = field
	const maxFields2zgensym_3cf2ad9d56d74f32_3 = 1

	// -- templateDecodeMsg starts here--
	var totalEncodedFields2zgensym_3cf2ad9d56d74f32_3 uint32
	totalEncodedFields2zgensym_3cf2ad9d56d74f32_3, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft2zgensym_3cf2ad9d56d74f32_3 := totalEncodedFields2zgensym_3cf2ad9d56d74f32_3
	missingFieldsLeft2zgensym_3cf2ad9d56d74f32_3 := maxFields2zgensym_3cf2ad9d56d74f32_3 - totalEncodedFields2zgensym_3cf2ad9d56d74f32_3

	var nextMiss2zgensym_3cf2ad9d56d74f32_3 int32 = -1
	var found2zgensym_3cf2ad9d56d74f32_3 [maxFields2zgensym_3cf2ad9d56d74f32_3]bool
	var curField2zgensym_3cf2ad9d56d74f32_3 string

doneWithStruct2zgensym_3cf2ad9d56d74f32_3:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft2zgensym_3cf2ad9d56d74f32_3 > 0 || missingFieldsLeft2zgensym_3cf2ad9d56d74f32_3 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft2zgensym_3cf2ad9d56d74f32_3, missingFieldsLeft2zgensym_3cf2ad9d56d74f32_3, msgp.ShowFound(found2zgensym_3cf2ad9d56d74f32_3[:]), decodeMsgFieldOrder2zgensym_3cf2ad9d56d74f32_3)
		if encodedFieldsLeft2zgensym_3cf2ad9d56d74f32_3 > 0 {
			encodedFieldsLeft2zgensym_3cf2ad9d56d74f32_3--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField2zgensym_3cf2ad9d56d74f32_3 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss2zgensym_3cf2ad9d56d74f32_3 < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss2zgensym_3cf2ad9d56d74f32_3 = 0
			}
			for nextMiss2zgensym_3cf2ad9d56d74f32_3 < maxFields2zgensym_3cf2ad9d56d74f32_3 && (found2zgensym_3cf2ad9d56d74f32_3[nextMiss2zgensym_3cf2ad9d56d74f32_3] || decodeMsgFieldSkip2zgensym_3cf2ad9d56d74f32_3[nextMiss2zgensym_3cf2ad9d56d74f32_3]) {
				nextMiss2zgensym_3cf2ad9d56d74f32_3++
			}
			if nextMiss2zgensym_3cf2ad9d56d74f32_3 == maxFields2zgensym_3cf2ad9d56d74f32_3 {
				// filled all the empty fields!
				break doneWithStruct2zgensym_3cf2ad9d56d74f32_3
			}
			missingFieldsLeft2zgensym_3cf2ad9d56d74f32_3--
			curField2zgensym_3cf2ad9d56d74f32_3 = decodeMsgFieldOrder2zgensym_3cf2ad9d56d74f32_3[nextMiss2zgensym_3cf2ad9d56d74f32_3]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField2zgensym_3cf2ad9d56d74f32_3)
		switch curField2zgensym_3cf2ad9d56d74f32_3 {
		// -- templateDecodeMsg ends here --

		case "U__map":
			found2zgensym_3cf2ad9d56d74f32_3[0] = true
			var zgensym_3cf2ad9d56d74f32_4 uint32
			zgensym_3cf2ad9d56d74f32_4, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.U == nil && zgensym_3cf2ad9d56d74f32_4 > 0 {
				z.U = make(map[string]*User, zgensym_3cf2ad9d56d74f32_4)
			} else if len(z.U) > 0 {
				for key, _ := range z.U {
					delete(z.U, key)
				}
			}
			for zgensym_3cf2ad9d56d74f32_4 > 0 {
				zgensym_3cf2ad9d56d74f32_4--
				var zgensym_3cf2ad9d56d74f32_0 string
				var zgensym_3cf2ad9d56d74f32_1 *User
				zgensym_3cf2ad9d56d74f32_0, err = dc.ReadString()
				if err != nil {
					return
				}
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						return
					}

					if zgensym_3cf2ad9d56d74f32_1 != nil {
						dc.PushAlwaysNil()
						err = zgensym_3cf2ad9d56d74f32_1.DecodeMsg(dc)
						if err != nil {
							return
						}
						dc.PopAlwaysNil()
					}
				} else {
					// not Nil, we have something to read

					if zgensym_3cf2ad9d56d74f32_1 == nil {
						zgensym_3cf2ad9d56d74f32_1 = new(User)
					}
					err = zgensym_3cf2ad9d56d74f32_1.DecodeMsg(dc)
					if err != nil {
						return
					}
				}
				z.U[zgensym_3cf2ad9d56d74f32_0] = zgensym_3cf2ad9d56d74f32_1
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	if nextMiss2zgensym_3cf2ad9d56d74f32_3 != -1 {
		dc.PopAlwaysNil()
	}

	if sawTopNil {
		dc.PopAlwaysNil()
	}

	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// fields of AtomicUserMap
var decodeMsgFieldOrder2zgensym_3cf2ad9d56d74f32_3 = []string{"U__map"}

var decodeMsgFieldSkip2zgensym_3cf2ad9d56d74f32_3 = []bool{false}

// fieldsNotEmpty supports omitempty tags
func (z *AtomicUserMap) fieldsNotEmpty(isempty []bool) uint32 {
	if len(isempty) == 0 {
		return 1
	}
	var fieldsInUse uint32 = 1
	isempty[0] = (len(z.U) == 0) // string, omitempty
	if isempty[0] {
		fieldsInUse--
	}

	return fieldsInUse
}

// EncodeMsg implements msgp.Encodable
func (z *AtomicUserMap) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// honor the omitempty tags
	var empty_zgensym_3cf2ad9d56d74f32_5 [1]bool
	fieldsInUse_zgensym_3cf2ad9d56d74f32_6 := z.fieldsNotEmpty(empty_zgensym_3cf2ad9d56d74f32_5[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zgensym_3cf2ad9d56d74f32_6)
	if err != nil {
		return err
	}

	if !empty_zgensym_3cf2ad9d56d74f32_5[0] {
		// write "U__map"
		err = en.Append(0xa6, 0x55, 0x5f, 0x5f, 0x6d, 0x61, 0x70)
		if err != nil {
			return err
		}
		err = en.WriteMapHeader(uint32(len(z.U)))
		if err != nil {
			return
		}
		for zgensym_3cf2ad9d56d74f32_0, zgensym_3cf2ad9d56d74f32_1 := range z.U {
			err = en.WriteString(zgensym_3cf2ad9d56d74f32_0)
			if err != nil {
				return
			}
			if zgensym_3cf2ad9d56d74f32_1 == nil {
				err = en.WriteNil()
				if err != nil {
					return
				}
			} else {
				err = zgensym_3cf2ad9d56d74f32_1.EncodeMsg(en)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

// MarshalMsg implements msgp.Marshaler
func (z *AtomicUserMap) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())

	// honor the omitempty tags
	var empty [1]bool
	fieldsInUse := z.fieldsNotEmpty(empty[:])
	o = msgp.AppendMapHeader(o, fieldsInUse)

	if !empty[0] {
		// string "U__map"
		o = append(o, 0xa6, 0x55, 0x5f, 0x5f, 0x6d, 0x61, 0x70)
		o = msgp.AppendMapHeader(o, uint32(len(z.U)))
		for zgensym_3cf2ad9d56d74f32_0, zgensym_3cf2ad9d56d74f32_1 := range z.U {
			o = msgp.AppendString(o, zgensym_3cf2ad9d56d74f32_0)
			if zgensym_3cf2ad9d56d74f32_1 == nil {
				o = msgp.AppendNil(o)
			} else {
				o, err = zgensym_3cf2ad9d56d74f32_1.MarshalMsg(o)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *AtomicUserMap) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *AtomicUserMap) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	var field []byte
	_ = field
	const maxFields7zgensym_3cf2ad9d56d74f32_8 = 1

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields7zgensym_3cf2ad9d56d74f32_8 uint32
	if !nbs.AlwaysNil {
		totalEncodedFields7zgensym_3cf2ad9d56d74f32_8, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft7zgensym_3cf2ad9d56d74f32_8 := totalEncodedFields7zgensym_3cf2ad9d56d74f32_8
	missingFieldsLeft7zgensym_3cf2ad9d56d74f32_8 := maxFields7zgensym_3cf2ad9d56d74f32_8 - totalEncodedFields7zgensym_3cf2ad9d56d74f32_8

	var nextMiss7zgensym_3cf2ad9d56d74f32_8 int32 = -1
	var found7zgensym_3cf2ad9d56d74f32_8 [maxFields7zgensym_3cf2ad9d56d74f32_8]bool
	var curField7zgensym_3cf2ad9d56d74f32_8 string

doneWithStruct7zgensym_3cf2ad9d56d74f32_8:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft7zgensym_3cf2ad9d56d74f32_8 > 0 || missingFieldsLeft7zgensym_3cf2ad9d56d74f32_8 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft7zgensym_3cf2ad9d56d74f32_8, missingFieldsLeft7zgensym_3cf2ad9d56d74f32_8, msgp.ShowFound(found7zgensym_3cf2ad9d56d74f32_8[:]), unmarshalMsgFieldOrder7zgensym_3cf2ad9d56d74f32_8)
		if encodedFieldsLeft7zgensym_3cf2ad9d56d74f32_8 > 0 {
			encodedFieldsLeft7zgensym_3cf2ad9d56d74f32_8--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField7zgensym_3cf2ad9d56d74f32_8 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss7zgensym_3cf2ad9d56d74f32_8 < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss7zgensym_3cf2ad9d56d74f32_8 = 0
			}
			for nextMiss7zgensym_3cf2ad9d56d74f32_8 < maxFields7zgensym_3cf2ad9d56d74f32_8 && (found7zgensym_3cf2ad9d56d74f32_8[nextMiss7zgensym_3cf2ad9d56d74f32_8] || unmarshalMsgFieldSkip7zgensym_3cf2ad9d56d74f32_8[nextMiss7zgensym_3cf2ad9d56d74f32_8]) {
				nextMiss7zgensym_3cf2ad9d56d74f32_8++
			}
			if nextMiss7zgensym_3cf2ad9d56d74f32_8 == maxFields7zgensym_3cf2ad9d56d74f32_8 {
				// filled all the empty fields!
				break doneWithStruct7zgensym_3cf2ad9d56d74f32_8
			}
			missingFieldsLeft7zgensym_3cf2ad9d56d74f32_8--
			curField7zgensym_3cf2ad9d56d74f32_8 = unmarshalMsgFieldOrder7zgensym_3cf2ad9d56d74f32_8[nextMiss7zgensym_3cf2ad9d56d74f32_8]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField7zgensym_3cf2ad9d56d74f32_8)
		switch curField7zgensym_3cf2ad9d56d74f32_8 {
		// -- templateUnmarshalMsg ends here --

		case "U__map":
			found7zgensym_3cf2ad9d56d74f32_8[0] = true
			if nbs.AlwaysNil {
				if len(z.U) > 0 {
					for key, _ := range z.U {
						delete(z.U, key)
					}
				}

			} else {

				var zgensym_3cf2ad9d56d74f32_9 uint32
				zgensym_3cf2ad9d56d74f32_9, bts, err = nbs.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				if z.U == nil && zgensym_3cf2ad9d56d74f32_9 > 0 {
					z.U = make(map[string]*User, zgensym_3cf2ad9d56d74f32_9)
				} else if len(z.U) > 0 {
					for key, _ := range z.U {
						delete(z.U, key)
					}
				}
				for zgensym_3cf2ad9d56d74f32_9 > 0 {
					var zgensym_3cf2ad9d56d74f32_0 string
					var zgensym_3cf2ad9d56d74f32_1 *User
					zgensym_3cf2ad9d56d74f32_9--
					zgensym_3cf2ad9d56d74f32_0, bts, err = nbs.ReadStringBytes(bts)
					if err != nil {
						return
					}
					if nbs.AlwaysNil {
						if zgensym_3cf2ad9d56d74f32_1 != nil {
							zgensym_3cf2ad9d56d74f32_1.UnmarshalMsg(msgp.OnlyNilSlice)
						}
					} else {
						// not nbs.AlwaysNil
						if msgp.IsNil(bts) {
							bts = bts[1:]
							if nil != zgensym_3cf2ad9d56d74f32_1 {
								zgensym_3cf2ad9d56d74f32_1.UnmarshalMsg(msgp.OnlyNilSlice)
							}
						} else {
							// not nbs.AlwaysNil and not IsNil(bts): have something to read

							if zgensym_3cf2ad9d56d74f32_1 == nil {
								zgensym_3cf2ad9d56d74f32_1 = new(User)
							}
							bts, err = zgensym_3cf2ad9d56d74f32_1.UnmarshalMsg(bts)
							if err != nil {
								return
							}
							if err != nil {
								return
							}
						}
					}
					z.U[zgensym_3cf2ad9d56d74f32_0] = zgensym_3cf2ad9d56d74f32_1
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	if nextMiss7zgensym_3cf2ad9d56d74f32_8 != -1 {
		bts = nbs.PopAlwaysNil()
	}

	if sawTopNil {
		bts = nbs.PopAlwaysNil()
	}
	o = bts
	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// fields of AtomicUserMap
var unmarshalMsgFieldOrder7zgensym_3cf2ad9d56d74f32_8 = []string{"U__map"}

var unmarshalMsgFieldSkip7zgensym_3cf2ad9d56d74f32_8 = []bool{false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *AtomicUserMap) Msgsize() (s int) {
	s = 1 + 7 + msgp.MapHeaderSize
	if z.U != nil {
		for zgensym_3cf2ad9d56d74f32_0, zgensym_3cf2ad9d56d74f32_1 := range z.U {
			_ = zgensym_3cf2ad9d56d74f32_1
			_ = zgensym_3cf2ad9d56d74f32_0
			s += msgp.StringPrefixSize + len(zgensym_3cf2ad9d56d74f32_0)
			if zgensym_3cf2ad9d56d74f32_1 == nil {
				s += msgp.NilSize
			} else {
				s += zgensym_3cf2ad9d56d74f32_1.Msgsize()
			}
		}
	}
	return
}
