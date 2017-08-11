package sshego

// NOTE: THIS FILE WAS PRODUCED BY THE
// GREENPACK CODE GENERATION TOOL (github.com/glycerine/greenpack)
// DO NOT EDIT

import (
	"github.com/glycerine/greenpack/msgp"
)

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *Filedb) DecodeMsg(dc *msgp.Reader) (err error) {
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
	const maxFields0zgensym_995050f2db4fcb57_1 = 1

	// -- templateDecodeMsg starts here--
	var totalEncodedFields0zgensym_995050f2db4fcb57_1 uint32
	totalEncodedFields0zgensym_995050f2db4fcb57_1, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft0zgensym_995050f2db4fcb57_1 := totalEncodedFields0zgensym_995050f2db4fcb57_1
	missingFieldsLeft0zgensym_995050f2db4fcb57_1 := maxFields0zgensym_995050f2db4fcb57_1 - totalEncodedFields0zgensym_995050f2db4fcb57_1

	var nextMiss0zgensym_995050f2db4fcb57_1 int32 = -1
	var found0zgensym_995050f2db4fcb57_1 [maxFields0zgensym_995050f2db4fcb57_1]bool
	var curField0zgensym_995050f2db4fcb57_1 string

doneWithStruct0zgensym_995050f2db4fcb57_1:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft0zgensym_995050f2db4fcb57_1 > 0 || missingFieldsLeft0zgensym_995050f2db4fcb57_1 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft0zgensym_995050f2db4fcb57_1, missingFieldsLeft0zgensym_995050f2db4fcb57_1, msgp.ShowFound(found0zgensym_995050f2db4fcb57_1[:]), decodeMsgFieldOrder0zgensym_995050f2db4fcb57_1)
		if encodedFieldsLeft0zgensym_995050f2db4fcb57_1 > 0 {
			encodedFieldsLeft0zgensym_995050f2db4fcb57_1--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField0zgensym_995050f2db4fcb57_1 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss0zgensym_995050f2db4fcb57_1 < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss0zgensym_995050f2db4fcb57_1 = 0
			}
			for nextMiss0zgensym_995050f2db4fcb57_1 < maxFields0zgensym_995050f2db4fcb57_1 && (found0zgensym_995050f2db4fcb57_1[nextMiss0zgensym_995050f2db4fcb57_1] || decodeMsgFieldSkip0zgensym_995050f2db4fcb57_1[nextMiss0zgensym_995050f2db4fcb57_1]) {
				nextMiss0zgensym_995050f2db4fcb57_1++
			}
			if nextMiss0zgensym_995050f2db4fcb57_1 == maxFields0zgensym_995050f2db4fcb57_1 {
				// filled all the empty fields!
				break doneWithStruct0zgensym_995050f2db4fcb57_1
			}
			missingFieldsLeft0zgensym_995050f2db4fcb57_1--
			curField0zgensym_995050f2db4fcb57_1 = decodeMsgFieldOrder0zgensym_995050f2db4fcb57_1[nextMiss0zgensym_995050f2db4fcb57_1]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField0zgensym_995050f2db4fcb57_1)
		switch curField0zgensym_995050f2db4fcb57_1 {
		// -- templateDecodeMsg ends here --

		case "HostDb_zid00_ptr":
			found0zgensym_995050f2db4fcb57_1[0] = true
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}

				if z.HostDb != nil {
					dc.PushAlwaysNil()
					err = z.HostDb.DecodeMsg(dc)
					if err != nil {
						return
					}
					dc.PopAlwaysNil()
				}
			} else {
				// not Nil, we have something to read

				if z.HostDb == nil {
					z.HostDb = new(HostDb)
				}
				err = z.HostDb.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	if nextMiss0zgensym_995050f2db4fcb57_1 != -1 {
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

// fields of Filedb
var decodeMsgFieldOrder0zgensym_995050f2db4fcb57_1 = []string{"HostDb_zid00_ptr"}

var decodeMsgFieldSkip0zgensym_995050f2db4fcb57_1 = []bool{false}

// fieldsNotEmpty supports omitempty tags
func (z *Filedb) fieldsNotEmpty(isempty []bool) uint32 {
	if len(isempty) == 0 {
		return 1
	}
	var fieldsInUse uint32 = 1
	isempty[0] = (z.HostDb == nil) // pointer, omitempty
	if isempty[0] {
		fieldsInUse--
	}

	return fieldsInUse
}

// EncodeMsg implements msgp.Encodable
func (z *Filedb) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// honor the omitempty tags
	var empty_zgensym_995050f2db4fcb57_2 [1]bool
	fieldsInUse_zgensym_995050f2db4fcb57_3 := z.fieldsNotEmpty(empty_zgensym_995050f2db4fcb57_2[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zgensym_995050f2db4fcb57_3)
	if err != nil {
		return err
	}

	if !empty_zgensym_995050f2db4fcb57_2[0] {
		// write "HostDb_zid00_ptr"
		err = en.Append(0xb0, 0x48, 0x6f, 0x73, 0x74, 0x44, 0x62, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x30, 0x5f, 0x70, 0x74, 0x72)
		if err != nil {
			return err
		}
		if z.HostDb == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.HostDb.EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}

	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Filedb) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())

	// honor the omitempty tags
	var empty [1]bool
	fieldsInUse := z.fieldsNotEmpty(empty[:])
	o = msgp.AppendMapHeader(o, fieldsInUse)

	if !empty[0] {
		// string "HostDb_zid00_ptr"
		o = append(o, 0xb0, 0x48, 0x6f, 0x73, 0x74, 0x44, 0x62, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x30, 0x5f, 0x70, 0x74, 0x72)
		if z.HostDb == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.HostDb.MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}

	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Filedb) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *Filedb) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	var field []byte
	_ = field
	const maxFields4zgensym_995050f2db4fcb57_5 = 1

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields4zgensym_995050f2db4fcb57_5 uint32
	if !nbs.AlwaysNil {
		totalEncodedFields4zgensym_995050f2db4fcb57_5, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft4zgensym_995050f2db4fcb57_5 := totalEncodedFields4zgensym_995050f2db4fcb57_5
	missingFieldsLeft4zgensym_995050f2db4fcb57_5 := maxFields4zgensym_995050f2db4fcb57_5 - totalEncodedFields4zgensym_995050f2db4fcb57_5

	var nextMiss4zgensym_995050f2db4fcb57_5 int32 = -1
	var found4zgensym_995050f2db4fcb57_5 [maxFields4zgensym_995050f2db4fcb57_5]bool
	var curField4zgensym_995050f2db4fcb57_5 string

doneWithStruct4zgensym_995050f2db4fcb57_5:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft4zgensym_995050f2db4fcb57_5 > 0 || missingFieldsLeft4zgensym_995050f2db4fcb57_5 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft4zgensym_995050f2db4fcb57_5, missingFieldsLeft4zgensym_995050f2db4fcb57_5, msgp.ShowFound(found4zgensym_995050f2db4fcb57_5[:]), unmarshalMsgFieldOrder4zgensym_995050f2db4fcb57_5)
		if encodedFieldsLeft4zgensym_995050f2db4fcb57_5 > 0 {
			encodedFieldsLeft4zgensym_995050f2db4fcb57_5--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField4zgensym_995050f2db4fcb57_5 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss4zgensym_995050f2db4fcb57_5 < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss4zgensym_995050f2db4fcb57_5 = 0
			}
			for nextMiss4zgensym_995050f2db4fcb57_5 < maxFields4zgensym_995050f2db4fcb57_5 && (found4zgensym_995050f2db4fcb57_5[nextMiss4zgensym_995050f2db4fcb57_5] || unmarshalMsgFieldSkip4zgensym_995050f2db4fcb57_5[nextMiss4zgensym_995050f2db4fcb57_5]) {
				nextMiss4zgensym_995050f2db4fcb57_5++
			}
			if nextMiss4zgensym_995050f2db4fcb57_5 == maxFields4zgensym_995050f2db4fcb57_5 {
				// filled all the empty fields!
				break doneWithStruct4zgensym_995050f2db4fcb57_5
			}
			missingFieldsLeft4zgensym_995050f2db4fcb57_5--
			curField4zgensym_995050f2db4fcb57_5 = unmarshalMsgFieldOrder4zgensym_995050f2db4fcb57_5[nextMiss4zgensym_995050f2db4fcb57_5]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField4zgensym_995050f2db4fcb57_5)
		switch curField4zgensym_995050f2db4fcb57_5 {
		// -- templateUnmarshalMsg ends here --

		case "HostDb_zid00_ptr":
			found4zgensym_995050f2db4fcb57_5[0] = true
			if nbs.AlwaysNil {
				if z.HostDb != nil {
					z.HostDb.UnmarshalMsg(msgp.OnlyNilSlice)
				}
			} else {
				// not nbs.AlwaysNil
				if msgp.IsNil(bts) {
					bts = bts[1:]
					if nil != z.HostDb {
						z.HostDb.UnmarshalMsg(msgp.OnlyNilSlice)
					}
				} else {
					// not nbs.AlwaysNil and not IsNil(bts): have something to read

					if z.HostDb == nil {
						z.HostDb = new(HostDb)
					}
					bts, err = z.HostDb.UnmarshalMsg(bts)
					if err != nil {
						return
					}
					if err != nil {
						return
					}
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	if nextMiss4zgensym_995050f2db4fcb57_5 != -1 {
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

// fields of Filedb
var unmarshalMsgFieldOrder4zgensym_995050f2db4fcb57_5 = []string{"HostDb_zid00_ptr"}

var unmarshalMsgFieldSkip4zgensym_995050f2db4fcb57_5 = []bool{false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Filedb) Msgsize() (s int) {
	s = 1 + 17
	if z.HostDb == nil {
		s += msgp.NilSize
	} else {
		s += z.HostDb.Msgsize()
	}
	return
}
