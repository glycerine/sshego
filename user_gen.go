package sshego

// NOTE: THIS FILE WAS PRODUCED BY THE
// GREENPACK CODE GENERATION TOOL (github.com/glycerine/greenpack)
// DO NOT EDIT

import (
	"github.com/glycerine/greenpack/msgp"
)

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *HostDb) DecodeMsg(dc *msgp.Reader) (err error) {
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
	const maxFields0zgensym_189e87a53e58dbf2_1 = 2

	// -- templateDecodeMsg starts here--
	var totalEncodedFields0zgensym_189e87a53e58dbf2_1 uint32
	totalEncodedFields0zgensym_189e87a53e58dbf2_1, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft0zgensym_189e87a53e58dbf2_1 := totalEncodedFields0zgensym_189e87a53e58dbf2_1
	missingFieldsLeft0zgensym_189e87a53e58dbf2_1 := maxFields0zgensym_189e87a53e58dbf2_1 - totalEncodedFields0zgensym_189e87a53e58dbf2_1

	var nextMiss0zgensym_189e87a53e58dbf2_1 int32 = -1
	var found0zgensym_189e87a53e58dbf2_1 [maxFields0zgensym_189e87a53e58dbf2_1]bool
	var curField0zgensym_189e87a53e58dbf2_1 string

doneWithStruct0zgensym_189e87a53e58dbf2_1:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft0zgensym_189e87a53e58dbf2_1 > 0 || missingFieldsLeft0zgensym_189e87a53e58dbf2_1 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft0zgensym_189e87a53e58dbf2_1, missingFieldsLeft0zgensym_189e87a53e58dbf2_1, msgp.ShowFound(found0zgensym_189e87a53e58dbf2_1[:]), decodeMsgFieldOrder0zgensym_189e87a53e58dbf2_1)
		if encodedFieldsLeft0zgensym_189e87a53e58dbf2_1 > 0 {
			encodedFieldsLeft0zgensym_189e87a53e58dbf2_1--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField0zgensym_189e87a53e58dbf2_1 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss0zgensym_189e87a53e58dbf2_1 < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss0zgensym_189e87a53e58dbf2_1 = 0
			}
			for nextMiss0zgensym_189e87a53e58dbf2_1 < maxFields0zgensym_189e87a53e58dbf2_1 && (found0zgensym_189e87a53e58dbf2_1[nextMiss0zgensym_189e87a53e58dbf2_1] || decodeMsgFieldSkip0zgensym_189e87a53e58dbf2_1[nextMiss0zgensym_189e87a53e58dbf2_1]) {
				nextMiss0zgensym_189e87a53e58dbf2_1++
			}
			if nextMiss0zgensym_189e87a53e58dbf2_1 == maxFields0zgensym_189e87a53e58dbf2_1 {
				// filled all the empty fields!
				break doneWithStruct0zgensym_189e87a53e58dbf2_1
			}
			missingFieldsLeft0zgensym_189e87a53e58dbf2_1--
			curField0zgensym_189e87a53e58dbf2_1 = decodeMsgFieldOrder0zgensym_189e87a53e58dbf2_1[nextMiss0zgensym_189e87a53e58dbf2_1]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField0zgensym_189e87a53e58dbf2_1)
		switch curField0zgensym_189e87a53e58dbf2_1 {
		// -- templateDecodeMsg ends here --

		case "Users__ptr":
			found0zgensym_189e87a53e58dbf2_1[0] = true
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}

				if z.Users != nil {
					dc.PushAlwaysNil()
					err = z.Users.DecodeMsg(dc)
					if err != nil {
						return
					}
					dc.PopAlwaysNil()
				}
			} else {
				// not Nil, we have something to read

				if z.Users == nil {
					z.Users = new(AtomicUserMap)
				}
				err = z.Users.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "HostPrivateKeyPath__str":
			found0zgensym_189e87a53e58dbf2_1[1] = true
			z.HostPrivateKeyPath, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	if nextMiss0zgensym_189e87a53e58dbf2_1 != -1 {
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

// fields of HostDb
var decodeMsgFieldOrder0zgensym_189e87a53e58dbf2_1 = []string{"Users__ptr", "HostPrivateKeyPath__str"}

var decodeMsgFieldSkip0zgensym_189e87a53e58dbf2_1 = []bool{false, false}

// fieldsNotEmpty supports omitempty tags
func (z *HostDb) fieldsNotEmpty(isempty []bool) uint32 {
	if len(isempty) == 0 {
		return 2
	}
	var fieldsInUse uint32 = 2
	isempty[0] = (z.Users == nil) // pointer, omitempty
	if isempty[0] {
		fieldsInUse--
	}
	isempty[1] = (len(z.HostPrivateKeyPath) == 0) // string, omitempty
	if isempty[1] {
		fieldsInUse--
	}

	return fieldsInUse
}

// EncodeMsg implements msgp.Encodable
func (z *HostDb) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// honor the omitempty tags
	var empty_zgensym_189e87a53e58dbf2_2 [2]bool
	fieldsInUse_zgensym_189e87a53e58dbf2_3 := z.fieldsNotEmpty(empty_zgensym_189e87a53e58dbf2_2[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zgensym_189e87a53e58dbf2_3)
	if err != nil {
		return err
	}

	if !empty_zgensym_189e87a53e58dbf2_2[0] {
		// write "Users__ptr"
		err = en.Append(0xaa, 0x55, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x5f, 0x70, 0x74, 0x72)
		if err != nil {
			return err
		}
		if z.Users == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.Users.EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_2[1] {
		// write "HostPrivateKeyPath__str"
		err = en.Append(0xb7, 0x48, 0x6f, 0x73, 0x74, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.HostPrivateKeyPath)
		if err != nil {
			return
		}
	}

	return
}

// MarshalMsg implements msgp.Marshaler
func (z *HostDb) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())

	// honor the omitempty tags
	var empty [2]bool
	fieldsInUse := z.fieldsNotEmpty(empty[:])
	o = msgp.AppendMapHeader(o, fieldsInUse)

	if !empty[0] {
		// string "Users__ptr"
		o = append(o, 0xaa, 0x55, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x5f, 0x70, 0x74, 0x72)
		if z.Users == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.Users.MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}

	if !empty[1] {
		// string "HostPrivateKeyPath__str"
		o = append(o, 0xb7, 0x48, 0x6f, 0x73, 0x74, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.HostPrivateKeyPath)
	}

	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *HostDb) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *HostDb) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	var field []byte
	_ = field
	const maxFields4zgensym_189e87a53e58dbf2_5 = 2

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields4zgensym_189e87a53e58dbf2_5 uint32
	if !nbs.AlwaysNil {
		totalEncodedFields4zgensym_189e87a53e58dbf2_5, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft4zgensym_189e87a53e58dbf2_5 := totalEncodedFields4zgensym_189e87a53e58dbf2_5
	missingFieldsLeft4zgensym_189e87a53e58dbf2_5 := maxFields4zgensym_189e87a53e58dbf2_5 - totalEncodedFields4zgensym_189e87a53e58dbf2_5

	var nextMiss4zgensym_189e87a53e58dbf2_5 int32 = -1
	var found4zgensym_189e87a53e58dbf2_5 [maxFields4zgensym_189e87a53e58dbf2_5]bool
	var curField4zgensym_189e87a53e58dbf2_5 string

doneWithStruct4zgensym_189e87a53e58dbf2_5:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft4zgensym_189e87a53e58dbf2_5 > 0 || missingFieldsLeft4zgensym_189e87a53e58dbf2_5 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft4zgensym_189e87a53e58dbf2_5, missingFieldsLeft4zgensym_189e87a53e58dbf2_5, msgp.ShowFound(found4zgensym_189e87a53e58dbf2_5[:]), unmarshalMsgFieldOrder4zgensym_189e87a53e58dbf2_5)
		if encodedFieldsLeft4zgensym_189e87a53e58dbf2_5 > 0 {
			encodedFieldsLeft4zgensym_189e87a53e58dbf2_5--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField4zgensym_189e87a53e58dbf2_5 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss4zgensym_189e87a53e58dbf2_5 < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss4zgensym_189e87a53e58dbf2_5 = 0
			}
			for nextMiss4zgensym_189e87a53e58dbf2_5 < maxFields4zgensym_189e87a53e58dbf2_5 && (found4zgensym_189e87a53e58dbf2_5[nextMiss4zgensym_189e87a53e58dbf2_5] || unmarshalMsgFieldSkip4zgensym_189e87a53e58dbf2_5[nextMiss4zgensym_189e87a53e58dbf2_5]) {
				nextMiss4zgensym_189e87a53e58dbf2_5++
			}
			if nextMiss4zgensym_189e87a53e58dbf2_5 == maxFields4zgensym_189e87a53e58dbf2_5 {
				// filled all the empty fields!
				break doneWithStruct4zgensym_189e87a53e58dbf2_5
			}
			missingFieldsLeft4zgensym_189e87a53e58dbf2_5--
			curField4zgensym_189e87a53e58dbf2_5 = unmarshalMsgFieldOrder4zgensym_189e87a53e58dbf2_5[nextMiss4zgensym_189e87a53e58dbf2_5]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField4zgensym_189e87a53e58dbf2_5)
		switch curField4zgensym_189e87a53e58dbf2_5 {
		// -- templateUnmarshalMsg ends here --

		case "Users__ptr":
			found4zgensym_189e87a53e58dbf2_5[0] = true
			if nbs.AlwaysNil {
				if z.Users != nil {
					z.Users.UnmarshalMsg(msgp.OnlyNilSlice)
				}
			} else {
				// not nbs.AlwaysNil
				if msgp.IsNil(bts) {
					bts = bts[1:]
					if nil != z.Users {
						z.Users.UnmarshalMsg(msgp.OnlyNilSlice)
					}
				} else {
					// not nbs.AlwaysNil and not IsNil(bts): have something to read

					if z.Users == nil {
						z.Users = new(AtomicUserMap)
					}
					bts, err = z.Users.UnmarshalMsg(bts)
					if err != nil {
						return
					}
					if err != nil {
						return
					}
				}
			}
		case "HostPrivateKeyPath__str":
			found4zgensym_189e87a53e58dbf2_5[1] = true
			z.HostPrivateKeyPath, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	if nextMiss4zgensym_189e87a53e58dbf2_5 != -1 {
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

// fields of HostDb
var unmarshalMsgFieldOrder4zgensym_189e87a53e58dbf2_5 = []string{"Users__ptr", "HostPrivateKeyPath__str"}

var unmarshalMsgFieldSkip4zgensym_189e87a53e58dbf2_5 = []bool{false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *HostDb) Msgsize() (s int) {
	s = 1 + 11
	if z.Users == nil {
		s += msgp.NilSize
	} else {
		s += z.Users.Msgsize()
	}
	s += 24 + msgp.StringPrefixSize + len(z.HostPrivateKeyPath)
	return
}

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *LoginRecord) DecodeMsg(dc *msgp.Reader) (err error) {
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
	const maxFields6zgensym_189e87a53e58dbf2_7 = 5

	// -- templateDecodeMsg starts here--
	var totalEncodedFields6zgensym_189e87a53e58dbf2_7 uint32
	totalEncodedFields6zgensym_189e87a53e58dbf2_7, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft6zgensym_189e87a53e58dbf2_7 := totalEncodedFields6zgensym_189e87a53e58dbf2_7
	missingFieldsLeft6zgensym_189e87a53e58dbf2_7 := maxFields6zgensym_189e87a53e58dbf2_7 - totalEncodedFields6zgensym_189e87a53e58dbf2_7

	var nextMiss6zgensym_189e87a53e58dbf2_7 int32 = -1
	var found6zgensym_189e87a53e58dbf2_7 [maxFields6zgensym_189e87a53e58dbf2_7]bool
	var curField6zgensym_189e87a53e58dbf2_7 string

doneWithStruct6zgensym_189e87a53e58dbf2_7:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft6zgensym_189e87a53e58dbf2_7 > 0 || missingFieldsLeft6zgensym_189e87a53e58dbf2_7 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft6zgensym_189e87a53e58dbf2_7, missingFieldsLeft6zgensym_189e87a53e58dbf2_7, msgp.ShowFound(found6zgensym_189e87a53e58dbf2_7[:]), decodeMsgFieldOrder6zgensym_189e87a53e58dbf2_7)
		if encodedFieldsLeft6zgensym_189e87a53e58dbf2_7 > 0 {
			encodedFieldsLeft6zgensym_189e87a53e58dbf2_7--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField6zgensym_189e87a53e58dbf2_7 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss6zgensym_189e87a53e58dbf2_7 < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss6zgensym_189e87a53e58dbf2_7 = 0
			}
			for nextMiss6zgensym_189e87a53e58dbf2_7 < maxFields6zgensym_189e87a53e58dbf2_7 && (found6zgensym_189e87a53e58dbf2_7[nextMiss6zgensym_189e87a53e58dbf2_7] || decodeMsgFieldSkip6zgensym_189e87a53e58dbf2_7[nextMiss6zgensym_189e87a53e58dbf2_7]) {
				nextMiss6zgensym_189e87a53e58dbf2_7++
			}
			if nextMiss6zgensym_189e87a53e58dbf2_7 == maxFields6zgensym_189e87a53e58dbf2_7 {
				// filled all the empty fields!
				break doneWithStruct6zgensym_189e87a53e58dbf2_7
			}
			missingFieldsLeft6zgensym_189e87a53e58dbf2_7--
			curField6zgensym_189e87a53e58dbf2_7 = decodeMsgFieldOrder6zgensym_189e87a53e58dbf2_7[nextMiss6zgensym_189e87a53e58dbf2_7]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField6zgensym_189e87a53e58dbf2_7)
		switch curField6zgensym_189e87a53e58dbf2_7 {
		// -- templateDecodeMsg ends here --

		case "FirstTm__tim":
			found6zgensym_189e87a53e58dbf2_7[0] = true
			z.FirstTm, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "LastTm__tim":
			found6zgensym_189e87a53e58dbf2_7[1] = true
			z.LastTm, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "SeenCount__i64":
			found6zgensym_189e87a53e58dbf2_7[2] = true
			z.SeenCount, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "AcceptedCount__i64":
			found6zgensym_189e87a53e58dbf2_7[3] = true
			z.AcceptedCount, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "PubFinger__str":
			found6zgensym_189e87a53e58dbf2_7[4] = true
			z.PubFinger, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	if nextMiss6zgensym_189e87a53e58dbf2_7 != -1 {
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

// fields of LoginRecord
var decodeMsgFieldOrder6zgensym_189e87a53e58dbf2_7 = []string{"FirstTm__tim", "LastTm__tim", "SeenCount__i64", "AcceptedCount__i64", "PubFinger__str"}

var decodeMsgFieldSkip6zgensym_189e87a53e58dbf2_7 = []bool{false, false, false, false, false}

// fieldsNotEmpty supports omitempty tags
func (z *LoginRecord) fieldsNotEmpty(isempty []bool) uint32 {
	if len(isempty) == 0 {
		return 5
	}
	var fieldsInUse uint32 = 5
	isempty[0] = (z.FirstTm.IsZero()) // time.Time, omitempty
	if isempty[0] {
		fieldsInUse--
	}
	isempty[1] = (z.LastTm.IsZero()) // time.Time, omitempty
	if isempty[1] {
		fieldsInUse--
	}
	isempty[2] = (z.SeenCount == 0) // number, omitempty
	if isempty[2] {
		fieldsInUse--
	}
	isempty[3] = (z.AcceptedCount == 0) // number, omitempty
	if isempty[3] {
		fieldsInUse--
	}
	isempty[4] = (len(z.PubFinger) == 0) // string, omitempty
	if isempty[4] {
		fieldsInUse--
	}

	return fieldsInUse
}

// EncodeMsg implements msgp.Encodable
func (z *LoginRecord) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// honor the omitempty tags
	var empty_zgensym_189e87a53e58dbf2_8 [5]bool
	fieldsInUse_zgensym_189e87a53e58dbf2_9 := z.fieldsNotEmpty(empty_zgensym_189e87a53e58dbf2_8[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zgensym_189e87a53e58dbf2_9)
	if err != nil {
		return err
	}

	if !empty_zgensym_189e87a53e58dbf2_8[0] {
		// write "FirstTm__tim"
		err = en.Append(0xac, 0x46, 0x69, 0x72, 0x73, 0x74, 0x54, 0x6d, 0x5f, 0x5f, 0x74, 0x69, 0x6d)
		if err != nil {
			return err
		}
		err = en.WriteTime(z.FirstTm)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_8[1] {
		// write "LastTm__tim"
		err = en.Append(0xab, 0x4c, 0x61, 0x73, 0x74, 0x54, 0x6d, 0x5f, 0x5f, 0x74, 0x69, 0x6d)
		if err != nil {
			return err
		}
		err = en.WriteTime(z.LastTm)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_8[2] {
		// write "SeenCount__i64"
		err = en.Append(0xae, 0x53, 0x65, 0x65, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x5f, 0x69, 0x36, 0x34)
		if err != nil {
			return err
		}
		err = en.WriteInt64(z.SeenCount)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_8[3] {
		// write "AcceptedCount__i64"
		err = en.Append(0xb2, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x5f, 0x69, 0x36, 0x34)
		if err != nil {
			return err
		}
		err = en.WriteInt64(z.AcceptedCount)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_8[4] {
		// write "PubFinger__str"
		err = en.Append(0xae, 0x50, 0x75, 0x62, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.PubFinger)
		if err != nil {
			return
		}
	}

	return
}

// MarshalMsg implements msgp.Marshaler
func (z *LoginRecord) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())

	// honor the omitempty tags
	var empty [5]bool
	fieldsInUse := z.fieldsNotEmpty(empty[:])
	o = msgp.AppendMapHeader(o, fieldsInUse)

	if !empty[0] {
		// string "FirstTm__tim"
		o = append(o, 0xac, 0x46, 0x69, 0x72, 0x73, 0x74, 0x54, 0x6d, 0x5f, 0x5f, 0x74, 0x69, 0x6d)
		o = msgp.AppendTime(o, z.FirstTm)
	}

	if !empty[1] {
		// string "LastTm__tim"
		o = append(o, 0xab, 0x4c, 0x61, 0x73, 0x74, 0x54, 0x6d, 0x5f, 0x5f, 0x74, 0x69, 0x6d)
		o = msgp.AppendTime(o, z.LastTm)
	}

	if !empty[2] {
		// string "SeenCount__i64"
		o = append(o, 0xae, 0x53, 0x65, 0x65, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x5f, 0x69, 0x36, 0x34)
		o = msgp.AppendInt64(o, z.SeenCount)
	}

	if !empty[3] {
		// string "AcceptedCount__i64"
		o = append(o, 0xb2, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x5f, 0x69, 0x36, 0x34)
		o = msgp.AppendInt64(o, z.AcceptedCount)
	}

	if !empty[4] {
		// string "PubFinger__str"
		o = append(o, 0xae, 0x50, 0x75, 0x62, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.PubFinger)
	}

	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *LoginRecord) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *LoginRecord) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	var field []byte
	_ = field
	const maxFields10zgensym_189e87a53e58dbf2_11 = 5

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields10zgensym_189e87a53e58dbf2_11 uint32
	if !nbs.AlwaysNil {
		totalEncodedFields10zgensym_189e87a53e58dbf2_11, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft10zgensym_189e87a53e58dbf2_11 := totalEncodedFields10zgensym_189e87a53e58dbf2_11
	missingFieldsLeft10zgensym_189e87a53e58dbf2_11 := maxFields10zgensym_189e87a53e58dbf2_11 - totalEncodedFields10zgensym_189e87a53e58dbf2_11

	var nextMiss10zgensym_189e87a53e58dbf2_11 int32 = -1
	var found10zgensym_189e87a53e58dbf2_11 [maxFields10zgensym_189e87a53e58dbf2_11]bool
	var curField10zgensym_189e87a53e58dbf2_11 string

doneWithStruct10zgensym_189e87a53e58dbf2_11:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft10zgensym_189e87a53e58dbf2_11 > 0 || missingFieldsLeft10zgensym_189e87a53e58dbf2_11 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft10zgensym_189e87a53e58dbf2_11, missingFieldsLeft10zgensym_189e87a53e58dbf2_11, msgp.ShowFound(found10zgensym_189e87a53e58dbf2_11[:]), unmarshalMsgFieldOrder10zgensym_189e87a53e58dbf2_11)
		if encodedFieldsLeft10zgensym_189e87a53e58dbf2_11 > 0 {
			encodedFieldsLeft10zgensym_189e87a53e58dbf2_11--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField10zgensym_189e87a53e58dbf2_11 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss10zgensym_189e87a53e58dbf2_11 < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss10zgensym_189e87a53e58dbf2_11 = 0
			}
			for nextMiss10zgensym_189e87a53e58dbf2_11 < maxFields10zgensym_189e87a53e58dbf2_11 && (found10zgensym_189e87a53e58dbf2_11[nextMiss10zgensym_189e87a53e58dbf2_11] || unmarshalMsgFieldSkip10zgensym_189e87a53e58dbf2_11[nextMiss10zgensym_189e87a53e58dbf2_11]) {
				nextMiss10zgensym_189e87a53e58dbf2_11++
			}
			if nextMiss10zgensym_189e87a53e58dbf2_11 == maxFields10zgensym_189e87a53e58dbf2_11 {
				// filled all the empty fields!
				break doneWithStruct10zgensym_189e87a53e58dbf2_11
			}
			missingFieldsLeft10zgensym_189e87a53e58dbf2_11--
			curField10zgensym_189e87a53e58dbf2_11 = unmarshalMsgFieldOrder10zgensym_189e87a53e58dbf2_11[nextMiss10zgensym_189e87a53e58dbf2_11]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField10zgensym_189e87a53e58dbf2_11)
		switch curField10zgensym_189e87a53e58dbf2_11 {
		// -- templateUnmarshalMsg ends here --

		case "FirstTm__tim":
			found10zgensym_189e87a53e58dbf2_11[0] = true
			z.FirstTm, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "LastTm__tim":
			found10zgensym_189e87a53e58dbf2_11[1] = true
			z.LastTm, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "SeenCount__i64":
			found10zgensym_189e87a53e58dbf2_11[2] = true
			z.SeenCount, bts, err = nbs.ReadInt64Bytes(bts)

			if err != nil {
				return
			}
		case "AcceptedCount__i64":
			found10zgensym_189e87a53e58dbf2_11[3] = true
			z.AcceptedCount, bts, err = nbs.ReadInt64Bytes(bts)

			if err != nil {
				return
			}
		case "PubFinger__str":
			found10zgensym_189e87a53e58dbf2_11[4] = true
			z.PubFinger, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	if nextMiss10zgensym_189e87a53e58dbf2_11 != -1 {
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

// fields of LoginRecord
var unmarshalMsgFieldOrder10zgensym_189e87a53e58dbf2_11 = []string{"FirstTm__tim", "LastTm__tim", "SeenCount__i64", "AcceptedCount__i64", "PubFinger__str"}

var unmarshalMsgFieldSkip10zgensym_189e87a53e58dbf2_11 = []bool{false, false, false, false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *LoginRecord) Msgsize() (s int) {
	s = 1 + 13 + msgp.TimeSize + 12 + msgp.TimeSize + 15 + msgp.Int64Size + 19 + msgp.Int64Size + 15 + msgp.StringPrefixSize + len(z.PubFinger)
	return
}

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *User) DecodeMsg(dc *msgp.Reader) (err error) {
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
	const maxFields15zgensym_189e87a53e58dbf2_16 = 17

	// -- templateDecodeMsg starts here--
	var totalEncodedFields15zgensym_189e87a53e58dbf2_16 uint32
	totalEncodedFields15zgensym_189e87a53e58dbf2_16, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft15zgensym_189e87a53e58dbf2_16 := totalEncodedFields15zgensym_189e87a53e58dbf2_16
	missingFieldsLeft15zgensym_189e87a53e58dbf2_16 := maxFields15zgensym_189e87a53e58dbf2_16 - totalEncodedFields15zgensym_189e87a53e58dbf2_16

	var nextMiss15zgensym_189e87a53e58dbf2_16 int32 = -1
	var found15zgensym_189e87a53e58dbf2_16 [maxFields15zgensym_189e87a53e58dbf2_16]bool
	var curField15zgensym_189e87a53e58dbf2_16 string

doneWithStruct15zgensym_189e87a53e58dbf2_16:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft15zgensym_189e87a53e58dbf2_16 > 0 || missingFieldsLeft15zgensym_189e87a53e58dbf2_16 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft15zgensym_189e87a53e58dbf2_16, missingFieldsLeft15zgensym_189e87a53e58dbf2_16, msgp.ShowFound(found15zgensym_189e87a53e58dbf2_16[:]), decodeMsgFieldOrder15zgensym_189e87a53e58dbf2_16)
		if encodedFieldsLeft15zgensym_189e87a53e58dbf2_16 > 0 {
			encodedFieldsLeft15zgensym_189e87a53e58dbf2_16--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField15zgensym_189e87a53e58dbf2_16 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss15zgensym_189e87a53e58dbf2_16 < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss15zgensym_189e87a53e58dbf2_16 = 0
			}
			for nextMiss15zgensym_189e87a53e58dbf2_16 < maxFields15zgensym_189e87a53e58dbf2_16 && (found15zgensym_189e87a53e58dbf2_16[nextMiss15zgensym_189e87a53e58dbf2_16] || decodeMsgFieldSkip15zgensym_189e87a53e58dbf2_16[nextMiss15zgensym_189e87a53e58dbf2_16]) {
				nextMiss15zgensym_189e87a53e58dbf2_16++
			}
			if nextMiss15zgensym_189e87a53e58dbf2_16 == maxFields15zgensym_189e87a53e58dbf2_16 {
				// filled all the empty fields!
				break doneWithStruct15zgensym_189e87a53e58dbf2_16
			}
			missingFieldsLeft15zgensym_189e87a53e58dbf2_16--
			curField15zgensym_189e87a53e58dbf2_16 = decodeMsgFieldOrder15zgensym_189e87a53e58dbf2_16[nextMiss15zgensym_189e87a53e58dbf2_16]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField15zgensym_189e87a53e58dbf2_16)
		switch curField15zgensym_189e87a53e58dbf2_16 {
		// -- templateDecodeMsg ends here --

		case "MyEmail__str":
			found15zgensym_189e87a53e58dbf2_16[0] = true
			z.MyEmail, err = dc.ReadString()
			if err != nil {
				return
			}
		case "MyFullname__str":
			found15zgensym_189e87a53e58dbf2_16[1] = true
			z.MyFullname, err = dc.ReadString()
			if err != nil {
				return
			}
		case "MyLogin__str":
			found15zgensym_189e87a53e58dbf2_16[2] = true
			z.MyLogin, err = dc.ReadString()
			if err != nil {
				return
			}
		case "PublicKeyPath__str":
			found15zgensym_189e87a53e58dbf2_16[3] = true
			z.PublicKeyPath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "PrivateKeyPath__str":
			found15zgensym_189e87a53e58dbf2_16[4] = true
			z.PrivateKeyPath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "TOTPpath__str":
			found15zgensym_189e87a53e58dbf2_16[5] = true
			z.TOTPpath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "QrPath__str":
			found15zgensym_189e87a53e58dbf2_16[6] = true
			z.QrPath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Issuer__str":
			found15zgensym_189e87a53e58dbf2_16[7] = true
			z.Issuer, err = dc.ReadString()
			if err != nil {
				return
			}
		case "SeenPubKey__map":
			found15zgensym_189e87a53e58dbf2_16[8] = true
			var zgensym_189e87a53e58dbf2_17 uint32
			zgensym_189e87a53e58dbf2_17, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.SeenPubKey == nil && zgensym_189e87a53e58dbf2_17 > 0 {
				z.SeenPubKey = make(map[string]LoginRecord, zgensym_189e87a53e58dbf2_17)
			} else if len(z.SeenPubKey) > 0 {
				for key, _ := range z.SeenPubKey {
					delete(z.SeenPubKey, key)
				}
			}
			for zgensym_189e87a53e58dbf2_17 > 0 {
				zgensym_189e87a53e58dbf2_17--
				var zgensym_189e87a53e58dbf2_12 string
				var zgensym_189e87a53e58dbf2_13 LoginRecord
				zgensym_189e87a53e58dbf2_12, err = dc.ReadString()
				if err != nil {
					return
				}
				err = zgensym_189e87a53e58dbf2_13.DecodeMsg(dc)
				if err != nil {
					return
				}
				z.SeenPubKey[zgensym_189e87a53e58dbf2_12] = zgensym_189e87a53e58dbf2_13
			}
		case "ScryptedPassword__bin":
			found15zgensym_189e87a53e58dbf2_16[9] = true
			z.ScryptedPassword, err = dc.ReadBytes(z.ScryptedPassword)
			if err != nil {
				return
			}
		case "ClearPw__str":
			found15zgensym_189e87a53e58dbf2_16[10] = true
			z.ClearPw, err = dc.ReadString()
			if err != nil {
				return
			}
		case "TOTPorig__str":
			found15zgensym_189e87a53e58dbf2_16[11] = true
			z.TOTPorig, err = dc.ReadString()
			if err != nil {
				return
			}
		case "FirstLoginTime__tim":
			found15zgensym_189e87a53e58dbf2_16[12] = true
			z.FirstLoginTime, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "LastLoginTime__tim":
			found15zgensym_189e87a53e58dbf2_16[13] = true
			z.LastLoginTime, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "LastLoginAddr__str":
			found15zgensym_189e87a53e58dbf2_16[14] = true
			z.LastLoginAddr, err = dc.ReadString()
			if err != nil {
				return
			}
		case "IPwhitelist__slc":
			found15zgensym_189e87a53e58dbf2_16[15] = true
			var zgensym_189e87a53e58dbf2_18 uint32
			zgensym_189e87a53e58dbf2_18, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.IPwhitelist) >= int(zgensym_189e87a53e58dbf2_18) {
				z.IPwhitelist = (z.IPwhitelist)[:zgensym_189e87a53e58dbf2_18]
			} else {
				z.IPwhitelist = make([]string, zgensym_189e87a53e58dbf2_18)
			}
			for zgensym_189e87a53e58dbf2_14 := range z.IPwhitelist {
				z.IPwhitelist[zgensym_189e87a53e58dbf2_14], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "DisabledAcct__boo":
			found15zgensym_189e87a53e58dbf2_16[16] = true
			z.DisabledAcct, err = dc.ReadBool()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	if nextMiss15zgensym_189e87a53e58dbf2_16 != -1 {
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

// fields of User
var decodeMsgFieldOrder15zgensym_189e87a53e58dbf2_16 = []string{"MyEmail__str", "MyFullname__str", "MyLogin__str", "PublicKeyPath__str", "PrivateKeyPath__str", "TOTPpath__str", "QrPath__str", "Issuer__str", "SeenPubKey__map", "ScryptedPassword__bin", "ClearPw__str", "TOTPorig__str", "FirstLoginTime__tim", "LastLoginTime__tim", "LastLoginAddr__str", "IPwhitelist__slc", "DisabledAcct__boo"}

var decodeMsgFieldSkip15zgensym_189e87a53e58dbf2_16 = []bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false}

// fieldsNotEmpty supports omitempty tags
func (z *User) fieldsNotEmpty(isempty []bool) uint32 {
	if len(isempty) == 0 {
		return 17
	}
	var fieldsInUse uint32 = 17
	isempty[0] = (len(z.MyEmail) == 0) // string, omitempty
	if isempty[0] {
		fieldsInUse--
	}
	isempty[1] = (len(z.MyFullname) == 0) // string, omitempty
	if isempty[1] {
		fieldsInUse--
	}
	isempty[2] = (len(z.MyLogin) == 0) // string, omitempty
	if isempty[2] {
		fieldsInUse--
	}
	isempty[3] = (len(z.PublicKeyPath) == 0) // string, omitempty
	if isempty[3] {
		fieldsInUse--
	}
	isempty[4] = (len(z.PrivateKeyPath) == 0) // string, omitempty
	if isempty[4] {
		fieldsInUse--
	}
	isempty[5] = (len(z.TOTPpath) == 0) // string, omitempty
	if isempty[5] {
		fieldsInUse--
	}
	isempty[6] = (len(z.QrPath) == 0) // string, omitempty
	if isempty[6] {
		fieldsInUse--
	}
	isempty[7] = (len(z.Issuer) == 0) // string, omitempty
	if isempty[7] {
		fieldsInUse--
	}
	isempty[8] = (len(z.SeenPubKey) == 0) // string, omitempty
	if isempty[8] {
		fieldsInUse--
	}
	isempty[9] = (len(z.ScryptedPassword) == 0) // string, omitempty
	if isempty[9] {
		fieldsInUse--
	}
	isempty[10] = (len(z.ClearPw) == 0) // string, omitempty
	if isempty[10] {
		fieldsInUse--
	}
	isempty[11] = (len(z.TOTPorig) == 0) // string, omitempty
	if isempty[11] {
		fieldsInUse--
	}
	isempty[12] = (z.FirstLoginTime.IsZero()) // time.Time, omitempty
	if isempty[12] {
		fieldsInUse--
	}
	isempty[13] = (z.LastLoginTime.IsZero()) // time.Time, omitempty
	if isempty[13] {
		fieldsInUse--
	}
	isempty[14] = (len(z.LastLoginAddr) == 0) // string, omitempty
	if isempty[14] {
		fieldsInUse--
	}
	isempty[15] = (len(z.IPwhitelist) == 0) // string, omitempty
	if isempty[15] {
		fieldsInUse--
	}
	isempty[16] = (!z.DisabledAcct) // bool, omitempty
	if isempty[16] {
		fieldsInUse--
	}

	return fieldsInUse
}

// EncodeMsg implements msgp.Encodable
func (z *User) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// honor the omitempty tags
	var empty_zgensym_189e87a53e58dbf2_19 [17]bool
	fieldsInUse_zgensym_189e87a53e58dbf2_20 := z.fieldsNotEmpty(empty_zgensym_189e87a53e58dbf2_19[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zgensym_189e87a53e58dbf2_20)
	if err != nil {
		return err
	}

	if !empty_zgensym_189e87a53e58dbf2_19[0] {
		// write "MyEmail__str"
		err = en.Append(0xac, 0x4d, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.MyEmail)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[1] {
		// write "MyFullname__str"
		err = en.Append(0xaf, 0x4d, 0x79, 0x46, 0x75, 0x6c, 0x6c, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.MyFullname)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[2] {
		// write "MyLogin__str"
		err = en.Append(0xac, 0x4d, 0x79, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.MyLogin)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[3] {
		// write "PublicKeyPath__str"
		err = en.Append(0xb2, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.PublicKeyPath)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[4] {
		// write "PrivateKeyPath__str"
		err = en.Append(0xb3, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.PrivateKeyPath)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[5] {
		// write "TOTPpath__str"
		err = en.Append(0xad, 0x54, 0x4f, 0x54, 0x50, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.TOTPpath)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[6] {
		// write "QrPath__str"
		err = en.Append(0xab, 0x51, 0x72, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.QrPath)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[7] {
		// write "Issuer__str"
		err = en.Append(0xab, 0x49, 0x73, 0x73, 0x75, 0x65, 0x72, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.Issuer)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[8] {
		// write "SeenPubKey__map"
		err = en.Append(0xaf, 0x53, 0x65, 0x65, 0x6e, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x5f, 0x5f, 0x6d, 0x61, 0x70)
		if err != nil {
			return err
		}
		err = en.WriteMapHeader(uint32(len(z.SeenPubKey)))
		if err != nil {
			return
		}
		for zgensym_189e87a53e58dbf2_12, zgensym_189e87a53e58dbf2_13 := range z.SeenPubKey {
			err = en.WriteString(zgensym_189e87a53e58dbf2_12)
			if err != nil {
				return
			}
			err = zgensym_189e87a53e58dbf2_13.EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[9] {
		// write "ScryptedPassword__bin"
		err = en.Append(0xb5, 0x53, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x5f, 0x62, 0x69, 0x6e)
		if err != nil {
			return err
		}
		err = en.WriteBytes(z.ScryptedPassword)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[10] {
		// write "ClearPw__str"
		err = en.Append(0xac, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x50, 0x77, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.ClearPw)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[11] {
		// write "TOTPorig__str"
		err = en.Append(0xad, 0x54, 0x4f, 0x54, 0x50, 0x6f, 0x72, 0x69, 0x67, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.TOTPorig)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[12] {
		// write "FirstLoginTime__tim"
		err = en.Append(0xb3, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x5f, 0x5f, 0x74, 0x69, 0x6d)
		if err != nil {
			return err
		}
		err = en.WriteTime(z.FirstLoginTime)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[13] {
		// write "LastLoginTime__tim"
		err = en.Append(0xb2, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x5f, 0x5f, 0x74, 0x69, 0x6d)
		if err != nil {
			return err
		}
		err = en.WriteTime(z.LastLoginTime)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[14] {
		// write "LastLoginAddr__str"
		err = en.Append(0xb2, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.LastLoginAddr)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[15] {
		// write "IPwhitelist__slc"
		err = en.Append(0xb0, 0x49, 0x50, 0x77, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x5f, 0x73, 0x6c, 0x63)
		if err != nil {
			return err
		}
		err = en.WriteArrayHeader(uint32(len(z.IPwhitelist)))
		if err != nil {
			return
		}
		for zgensym_189e87a53e58dbf2_14 := range z.IPwhitelist {
			err = en.WriteString(z.IPwhitelist[zgensym_189e87a53e58dbf2_14])
			if err != nil {
				return
			}
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_19[16] {
		// write "DisabledAcct__boo"
		err = en.Append(0xb1, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x41, 0x63, 0x63, 0x74, 0x5f, 0x5f, 0x62, 0x6f, 0x6f)
		if err != nil {
			return err
		}
		err = en.WriteBool(z.DisabledAcct)
		if err != nil {
			return
		}
	}

	return
}

// MarshalMsg implements msgp.Marshaler
func (z *User) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())

	// honor the omitempty tags
	var empty [17]bool
	fieldsInUse := z.fieldsNotEmpty(empty[:])
	o = msgp.AppendMapHeader(o, fieldsInUse)

	if !empty[0] {
		// string "MyEmail__str"
		o = append(o, 0xac, 0x4d, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.MyEmail)
	}

	if !empty[1] {
		// string "MyFullname__str"
		o = append(o, 0xaf, 0x4d, 0x79, 0x46, 0x75, 0x6c, 0x6c, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.MyFullname)
	}

	if !empty[2] {
		// string "MyLogin__str"
		o = append(o, 0xac, 0x4d, 0x79, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.MyLogin)
	}

	if !empty[3] {
		// string "PublicKeyPath__str"
		o = append(o, 0xb2, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.PublicKeyPath)
	}

	if !empty[4] {
		// string "PrivateKeyPath__str"
		o = append(o, 0xb3, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.PrivateKeyPath)
	}

	if !empty[5] {
		// string "TOTPpath__str"
		o = append(o, 0xad, 0x54, 0x4f, 0x54, 0x50, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.TOTPpath)
	}

	if !empty[6] {
		// string "QrPath__str"
		o = append(o, 0xab, 0x51, 0x72, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.QrPath)
	}

	if !empty[7] {
		// string "Issuer__str"
		o = append(o, 0xab, 0x49, 0x73, 0x73, 0x75, 0x65, 0x72, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.Issuer)
	}

	if !empty[8] {
		// string "SeenPubKey__map"
		o = append(o, 0xaf, 0x53, 0x65, 0x65, 0x6e, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x5f, 0x5f, 0x6d, 0x61, 0x70)
		o = msgp.AppendMapHeader(o, uint32(len(z.SeenPubKey)))
		for zgensym_189e87a53e58dbf2_12, zgensym_189e87a53e58dbf2_13 := range z.SeenPubKey {
			o = msgp.AppendString(o, zgensym_189e87a53e58dbf2_12)
			o, err = zgensym_189e87a53e58dbf2_13.MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}

	if !empty[9] {
		// string "ScryptedPassword__bin"
		o = append(o, 0xb5, 0x53, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x5f, 0x62, 0x69, 0x6e)
		o = msgp.AppendBytes(o, z.ScryptedPassword)
	}

	if !empty[10] {
		// string "ClearPw__str"
		o = append(o, 0xac, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x50, 0x77, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.ClearPw)
	}

	if !empty[11] {
		// string "TOTPorig__str"
		o = append(o, 0xad, 0x54, 0x4f, 0x54, 0x50, 0x6f, 0x72, 0x69, 0x67, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.TOTPorig)
	}

	if !empty[12] {
		// string "FirstLoginTime__tim"
		o = append(o, 0xb3, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x5f, 0x5f, 0x74, 0x69, 0x6d)
		o = msgp.AppendTime(o, z.FirstLoginTime)
	}

	if !empty[13] {
		// string "LastLoginTime__tim"
		o = append(o, 0xb2, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x5f, 0x5f, 0x74, 0x69, 0x6d)
		o = msgp.AppendTime(o, z.LastLoginTime)
	}

	if !empty[14] {
		// string "LastLoginAddr__str"
		o = append(o, 0xb2, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.LastLoginAddr)
	}

	if !empty[15] {
		// string "IPwhitelist__slc"
		o = append(o, 0xb0, 0x49, 0x50, 0x77, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x5f, 0x73, 0x6c, 0x63)
		o = msgp.AppendArrayHeader(o, uint32(len(z.IPwhitelist)))
		for zgensym_189e87a53e58dbf2_14 := range z.IPwhitelist {
			o = msgp.AppendString(o, z.IPwhitelist[zgensym_189e87a53e58dbf2_14])
		}
	}

	if !empty[16] {
		// string "DisabledAcct__boo"
		o = append(o, 0xb1, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x41, 0x63, 0x63, 0x74, 0x5f, 0x5f, 0x62, 0x6f, 0x6f)
		o = msgp.AppendBool(o, z.DisabledAcct)
	}

	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *User) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *User) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	var field []byte
	_ = field
	const maxFields21zgensym_189e87a53e58dbf2_22 = 17

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields21zgensym_189e87a53e58dbf2_22 uint32
	if !nbs.AlwaysNil {
		totalEncodedFields21zgensym_189e87a53e58dbf2_22, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft21zgensym_189e87a53e58dbf2_22 := totalEncodedFields21zgensym_189e87a53e58dbf2_22
	missingFieldsLeft21zgensym_189e87a53e58dbf2_22 := maxFields21zgensym_189e87a53e58dbf2_22 - totalEncodedFields21zgensym_189e87a53e58dbf2_22

	var nextMiss21zgensym_189e87a53e58dbf2_22 int32 = -1
	var found21zgensym_189e87a53e58dbf2_22 [maxFields21zgensym_189e87a53e58dbf2_22]bool
	var curField21zgensym_189e87a53e58dbf2_22 string

doneWithStruct21zgensym_189e87a53e58dbf2_22:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft21zgensym_189e87a53e58dbf2_22 > 0 || missingFieldsLeft21zgensym_189e87a53e58dbf2_22 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft21zgensym_189e87a53e58dbf2_22, missingFieldsLeft21zgensym_189e87a53e58dbf2_22, msgp.ShowFound(found21zgensym_189e87a53e58dbf2_22[:]), unmarshalMsgFieldOrder21zgensym_189e87a53e58dbf2_22)
		if encodedFieldsLeft21zgensym_189e87a53e58dbf2_22 > 0 {
			encodedFieldsLeft21zgensym_189e87a53e58dbf2_22--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField21zgensym_189e87a53e58dbf2_22 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss21zgensym_189e87a53e58dbf2_22 < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss21zgensym_189e87a53e58dbf2_22 = 0
			}
			for nextMiss21zgensym_189e87a53e58dbf2_22 < maxFields21zgensym_189e87a53e58dbf2_22 && (found21zgensym_189e87a53e58dbf2_22[nextMiss21zgensym_189e87a53e58dbf2_22] || unmarshalMsgFieldSkip21zgensym_189e87a53e58dbf2_22[nextMiss21zgensym_189e87a53e58dbf2_22]) {
				nextMiss21zgensym_189e87a53e58dbf2_22++
			}
			if nextMiss21zgensym_189e87a53e58dbf2_22 == maxFields21zgensym_189e87a53e58dbf2_22 {
				// filled all the empty fields!
				break doneWithStruct21zgensym_189e87a53e58dbf2_22
			}
			missingFieldsLeft21zgensym_189e87a53e58dbf2_22--
			curField21zgensym_189e87a53e58dbf2_22 = unmarshalMsgFieldOrder21zgensym_189e87a53e58dbf2_22[nextMiss21zgensym_189e87a53e58dbf2_22]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField21zgensym_189e87a53e58dbf2_22)
		switch curField21zgensym_189e87a53e58dbf2_22 {
		// -- templateUnmarshalMsg ends here --

		case "MyEmail__str":
			found21zgensym_189e87a53e58dbf2_22[0] = true
			z.MyEmail, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "MyFullname__str":
			found21zgensym_189e87a53e58dbf2_22[1] = true
			z.MyFullname, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "MyLogin__str":
			found21zgensym_189e87a53e58dbf2_22[2] = true
			z.MyLogin, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "PublicKeyPath__str":
			found21zgensym_189e87a53e58dbf2_22[3] = true
			z.PublicKeyPath, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "PrivateKeyPath__str":
			found21zgensym_189e87a53e58dbf2_22[4] = true
			z.PrivateKeyPath, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "TOTPpath__str":
			found21zgensym_189e87a53e58dbf2_22[5] = true
			z.TOTPpath, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "QrPath__str":
			found21zgensym_189e87a53e58dbf2_22[6] = true
			z.QrPath, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Issuer__str":
			found21zgensym_189e87a53e58dbf2_22[7] = true
			z.Issuer, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "SeenPubKey__map":
			found21zgensym_189e87a53e58dbf2_22[8] = true
			if nbs.AlwaysNil {
				if len(z.SeenPubKey) > 0 {
					for key, _ := range z.SeenPubKey {
						delete(z.SeenPubKey, key)
					}
				}

			} else {

				var zgensym_189e87a53e58dbf2_23 uint32
				zgensym_189e87a53e58dbf2_23, bts, err = nbs.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				if z.SeenPubKey == nil && zgensym_189e87a53e58dbf2_23 > 0 {
					z.SeenPubKey = make(map[string]LoginRecord, zgensym_189e87a53e58dbf2_23)
				} else if len(z.SeenPubKey) > 0 {
					for key, _ := range z.SeenPubKey {
						delete(z.SeenPubKey, key)
					}
				}
				for zgensym_189e87a53e58dbf2_23 > 0 {
					var zgensym_189e87a53e58dbf2_12 string
					var zgensym_189e87a53e58dbf2_13 LoginRecord
					zgensym_189e87a53e58dbf2_23--
					zgensym_189e87a53e58dbf2_12, bts, err = nbs.ReadStringBytes(bts)
					if err != nil {
						return
					}
					bts, err = zgensym_189e87a53e58dbf2_13.UnmarshalMsg(bts)
					if err != nil {
						return
					}
					if err != nil {
						return
					}
					z.SeenPubKey[zgensym_189e87a53e58dbf2_12] = zgensym_189e87a53e58dbf2_13
				}
			}
		case "ScryptedPassword__bin":
			found21zgensym_189e87a53e58dbf2_22[9] = true
			if nbs.AlwaysNil || msgp.IsNil(bts) {
				if !nbs.AlwaysNil {
					bts = bts[1:]
				}
				z.ScryptedPassword = z.ScryptedPassword[:0]
			} else {
				z.ScryptedPassword, bts, err = nbs.ReadBytesBytes(bts, z.ScryptedPassword)

				if err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		case "ClearPw__str":
			found21zgensym_189e87a53e58dbf2_22[10] = true
			z.ClearPw, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "TOTPorig__str":
			found21zgensym_189e87a53e58dbf2_22[11] = true
			z.TOTPorig, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "FirstLoginTime__tim":
			found21zgensym_189e87a53e58dbf2_22[12] = true
			z.FirstLoginTime, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "LastLoginTime__tim":
			found21zgensym_189e87a53e58dbf2_22[13] = true
			z.LastLoginTime, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "LastLoginAddr__str":
			found21zgensym_189e87a53e58dbf2_22[14] = true
			z.LastLoginAddr, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "IPwhitelist__slc":
			found21zgensym_189e87a53e58dbf2_22[15] = true
			if nbs.AlwaysNil {
				(z.IPwhitelist) = (z.IPwhitelist)[:0]
			} else {

				var zgensym_189e87a53e58dbf2_24 uint32
				zgensym_189e87a53e58dbf2_24, bts, err = nbs.ReadArrayHeaderBytes(bts)
				if err != nil {
					return
				}
				if cap(z.IPwhitelist) >= int(zgensym_189e87a53e58dbf2_24) {
					z.IPwhitelist = (z.IPwhitelist)[:zgensym_189e87a53e58dbf2_24]
				} else {
					z.IPwhitelist = make([]string, zgensym_189e87a53e58dbf2_24)
				}
				for zgensym_189e87a53e58dbf2_14 := range z.IPwhitelist {
					z.IPwhitelist[zgensym_189e87a53e58dbf2_14], bts, err = nbs.ReadStringBytes(bts)

					if err != nil {
						return
					}
				}
			}
		case "DisabledAcct__boo":
			found21zgensym_189e87a53e58dbf2_22[16] = true
			z.DisabledAcct, bts, err = nbs.ReadBoolBytes(bts)

			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	if nextMiss21zgensym_189e87a53e58dbf2_22 != -1 {
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

// fields of User
var unmarshalMsgFieldOrder21zgensym_189e87a53e58dbf2_22 = []string{"MyEmail__str", "MyFullname__str", "MyLogin__str", "PublicKeyPath__str", "PrivateKeyPath__str", "TOTPpath__str", "QrPath__str", "Issuer__str", "SeenPubKey__map", "ScryptedPassword__bin", "ClearPw__str", "TOTPorig__str", "FirstLoginTime__tim", "LastLoginTime__tim", "LastLoginAddr__str", "IPwhitelist__slc", "DisabledAcct__boo"}

var unmarshalMsgFieldSkip21zgensym_189e87a53e58dbf2_22 = []bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *User) Msgsize() (s int) {
	s = 3 + 13 + msgp.StringPrefixSize + len(z.MyEmail) + 16 + msgp.StringPrefixSize + len(z.MyFullname) + 13 + msgp.StringPrefixSize + len(z.MyLogin) + 19 + msgp.StringPrefixSize + len(z.PublicKeyPath) + 20 + msgp.StringPrefixSize + len(z.PrivateKeyPath) + 14 + msgp.StringPrefixSize + len(z.TOTPpath) + 12 + msgp.StringPrefixSize + len(z.QrPath) + 12 + msgp.StringPrefixSize + len(z.Issuer) + 16 + msgp.MapHeaderSize
	if z.SeenPubKey != nil {
		for zgensym_189e87a53e58dbf2_12, zgensym_189e87a53e58dbf2_13 := range z.SeenPubKey {
			_ = zgensym_189e87a53e58dbf2_13
			_ = zgensym_189e87a53e58dbf2_12
			s += msgp.StringPrefixSize + len(zgensym_189e87a53e58dbf2_12) + zgensym_189e87a53e58dbf2_13.Msgsize()
		}
	}
	s += 22 + msgp.BytesPrefixSize + len(z.ScryptedPassword) + 13 + msgp.StringPrefixSize + len(z.ClearPw) + 14 + msgp.StringPrefixSize + len(z.TOTPorig) + 20 + msgp.TimeSize + 19 + msgp.TimeSize + 19 + msgp.StringPrefixSize + len(z.LastLoginAddr) + 17 + msgp.ArrayHeaderSize
	for zgensym_189e87a53e58dbf2_14 := range z.IPwhitelist {
		s += msgp.StringPrefixSize + len(z.IPwhitelist[zgensym_189e87a53e58dbf2_14])
	}
	s += 18 + msgp.BoolSize
	return
}
