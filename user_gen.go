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
	const maxFields0zgensym_189e87a53e58dbf2_1 = 3

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

		case "UserHomePrefix__str":
			found0zgensym_189e87a53e58dbf2_1[0] = true
			z.UserHomePrefix, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Persist__rct":
			found0zgensym_189e87a53e58dbf2_1[2] = true
			const maxFields2zgensym_189e87a53e58dbf2_3 = 2

			// -- templateDecodeMsg starts here--
			var totalEncodedFields2zgensym_189e87a53e58dbf2_3 uint32
			totalEncodedFields2zgensym_189e87a53e58dbf2_3, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			encodedFieldsLeft2zgensym_189e87a53e58dbf2_3 := totalEncodedFields2zgensym_189e87a53e58dbf2_3
			missingFieldsLeft2zgensym_189e87a53e58dbf2_3 := maxFields2zgensym_189e87a53e58dbf2_3 - totalEncodedFields2zgensym_189e87a53e58dbf2_3

			var nextMiss2zgensym_189e87a53e58dbf2_3 int32 = -1
			var found2zgensym_189e87a53e58dbf2_3 [maxFields2zgensym_189e87a53e58dbf2_3]bool
			var curField2zgensym_189e87a53e58dbf2_3 string

		doneWithStruct2zgensym_189e87a53e58dbf2_3:
			// First fill all the encoded fields, then
			// treat the remaining, missing fields, as Nil.
			for encodedFieldsLeft2zgensym_189e87a53e58dbf2_3 > 0 || missingFieldsLeft2zgensym_189e87a53e58dbf2_3 > 0 {
				//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft2zgensym_189e87a53e58dbf2_3, missingFieldsLeft2zgensym_189e87a53e58dbf2_3, msgp.ShowFound(found2zgensym_189e87a53e58dbf2_3[:]), decodeMsgFieldOrder2zgensym_189e87a53e58dbf2_3)
				if encodedFieldsLeft2zgensym_189e87a53e58dbf2_3 > 0 {
					encodedFieldsLeft2zgensym_189e87a53e58dbf2_3--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					curField2zgensym_189e87a53e58dbf2_3 = msgp.UnsafeString(field)
				} else {
					//missing fields need handling
					if nextMiss2zgensym_189e87a53e58dbf2_3 < 0 {
						// tell the reader to only give us Nils
						// until further notice.
						dc.PushAlwaysNil()
						nextMiss2zgensym_189e87a53e58dbf2_3 = 0
					}
					for nextMiss2zgensym_189e87a53e58dbf2_3 < maxFields2zgensym_189e87a53e58dbf2_3 && (found2zgensym_189e87a53e58dbf2_3[nextMiss2zgensym_189e87a53e58dbf2_3] || decodeMsgFieldSkip2zgensym_189e87a53e58dbf2_3[nextMiss2zgensym_189e87a53e58dbf2_3]) {
						nextMiss2zgensym_189e87a53e58dbf2_3++
					}
					if nextMiss2zgensym_189e87a53e58dbf2_3 == maxFields2zgensym_189e87a53e58dbf2_3 {
						// filled all the empty fields!
						break doneWithStruct2zgensym_189e87a53e58dbf2_3
					}
					missingFieldsLeft2zgensym_189e87a53e58dbf2_3--
					curField2zgensym_189e87a53e58dbf2_3 = decodeMsgFieldOrder2zgensym_189e87a53e58dbf2_3[nextMiss2zgensym_189e87a53e58dbf2_3]
				}
				//fmt.Printf("switching on curField: '%v'\n", curField2zgensym_189e87a53e58dbf2_3)
				switch curField2zgensym_189e87a53e58dbf2_3 {
				// -- templateDecodeMsg ends here --

				case "Users_zid00_ptr":
					found2zgensym_189e87a53e58dbf2_3[0] = true
					if dc.IsNil() {
						err = dc.ReadNil()
						if err != nil {
							return
						}

						if z.Persist.Users != nil {
							dc.PushAlwaysNil()
							err = z.Persist.Users.DecodeMsg(dc)
							if err != nil {
								return
							}
							dc.PopAlwaysNil()
						}
					} else {
						// not Nil, we have something to read

						if z.Persist.Users == nil {
							z.Persist.Users = new(AtomicUserMap)
						}
						err = z.Persist.Users.DecodeMsg(dc)
						if err != nil {
							return
						}
					}
				case "HostPrivateKeyPath_zid01_str":
					found2zgensym_189e87a53e58dbf2_3[1] = true
					z.Persist.HostPrivateKeyPath, err = dc.ReadString()
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
			if nextMiss2zgensym_189e87a53e58dbf2_3 != -1 {
				dc.PopAlwaysNil()
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
var decodeMsgFieldOrder0zgensym_189e87a53e58dbf2_1 = []string{"UserHomePrefix__str", "", "Persist__rct"}

var decodeMsgFieldSkip0zgensym_189e87a53e58dbf2_1 = []bool{false, true, false}

// fields of HostDbPersist
var decodeMsgFieldOrder2zgensym_189e87a53e58dbf2_3 = []string{"Users_zid00_ptr", "HostPrivateKeyPath_zid01_str"}

var decodeMsgFieldSkip2zgensym_189e87a53e58dbf2_3 = []bool{false, false}

// fieldsNotEmpty supports omitempty tags
func (z *HostDb) fieldsNotEmpty(isempty []bool) uint32 {
	if len(isempty) == 0 {
		return 2
	}
	var fieldsInUse uint32 = 2
	isempty[0] = (len(z.UserHomePrefix) == 0) // string, omitempty
	if isempty[0] {
		fieldsInUse--
	}
	isempty[2] = false // struct values are never empty
	if isempty[2] {
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
	var empty_zgensym_189e87a53e58dbf2_4 [3]bool
	fieldsInUse_zgensym_189e87a53e58dbf2_5 := z.fieldsNotEmpty(empty_zgensym_189e87a53e58dbf2_4[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zgensym_189e87a53e58dbf2_5)
	if err != nil {
		return err
	}

	if !empty_zgensym_189e87a53e58dbf2_4[0] {
		// write "UserHomePrefix__str"
		err = en.Append(0xb3, 0x55, 0x73, 0x65, 0x72, 0x48, 0x6f, 0x6d, 0x65, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.UserHomePrefix)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_4[2] {
		// write "Persist__rct"
		err = en.Append(0xac, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x5f, 0x5f, 0x72, 0x63, 0x74)
		if err != nil {
			return err
		}

		// honor the omitempty tags
		var empty_zgensym_189e87a53e58dbf2_6 [2]bool
		fieldsInUse_zgensym_189e87a53e58dbf2_7 := z.Persist.fieldsNotEmpty(empty_zgensym_189e87a53e58dbf2_6[:])

		// map header
		err = en.WriteMapHeader(fieldsInUse_zgensym_189e87a53e58dbf2_7)
		if err != nil {
			return err
		}

		if !empty_zgensym_189e87a53e58dbf2_6[0] {
			// write "Users_zid00_ptr"
			err = en.Append(0xaf, 0x55, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x30, 0x5f, 0x70, 0x74, 0x72)
			if err != nil {
				return err
			}
			if z.Persist.Users == nil {
				err = en.WriteNil()
				if err != nil {
					return
				}
			} else {
				err = z.Persist.Users.EncodeMsg(en)
				if err != nil {
					return
				}
			}
		}

		if !empty_zgensym_189e87a53e58dbf2_6[1] {
			// write "HostPrivateKeyPath_zid01_str"
			err = en.Append(0xbc, 0x48, 0x6f, 0x73, 0x74, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x31, 0x5f, 0x73, 0x74, 0x72)
			if err != nil {
				return err
			}
			err = en.WriteString(z.Persist.HostPrivateKeyPath)
			if err != nil {
				return
			}
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
	var empty [3]bool
	fieldsInUse := z.fieldsNotEmpty(empty[:])
	o = msgp.AppendMapHeader(o, fieldsInUse)

	if !empty[0] {
		// string "UserHomePrefix__str"
		o = append(o, 0xb3, 0x55, 0x73, 0x65, 0x72, 0x48, 0x6f, 0x6d, 0x65, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.UserHomePrefix)
	}

	if !empty[2] {
		// string "Persist__rct"
		o = append(o, 0xac, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x5f, 0x5f, 0x72, 0x63, 0x74)

		// honor the omitempty tags
		var empty [2]bool
		fieldsInUse := z.Persist.fieldsNotEmpty(empty[:])
		o = msgp.AppendMapHeader(o, fieldsInUse)

		if !empty[0] {
			// string "Users_zid00_ptr"
			o = append(o, 0xaf, 0x55, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x30, 0x5f, 0x70, 0x74, 0x72)
			if z.Persist.Users == nil {
				o = msgp.AppendNil(o)
			} else {
				o, err = z.Persist.Users.MarshalMsg(o)
				if err != nil {
					return
				}
			}
		}

		if !empty[1] {
			// string "HostPrivateKeyPath_zid01_str"
			o = append(o, 0xbc, 0x48, 0x6f, 0x73, 0x74, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x31, 0x5f, 0x73, 0x74, 0x72)
			o = msgp.AppendString(o, z.Persist.HostPrivateKeyPath)
		}

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
	const maxFields8zgensym_189e87a53e58dbf2_9 = 3

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields8zgensym_189e87a53e58dbf2_9 uint32
	if !nbs.AlwaysNil {
		totalEncodedFields8zgensym_189e87a53e58dbf2_9, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft8zgensym_189e87a53e58dbf2_9 := totalEncodedFields8zgensym_189e87a53e58dbf2_9
	missingFieldsLeft8zgensym_189e87a53e58dbf2_9 := maxFields8zgensym_189e87a53e58dbf2_9 - totalEncodedFields8zgensym_189e87a53e58dbf2_9

	var nextMiss8zgensym_189e87a53e58dbf2_9 int32 = -1
	var found8zgensym_189e87a53e58dbf2_9 [maxFields8zgensym_189e87a53e58dbf2_9]bool
	var curField8zgensym_189e87a53e58dbf2_9 string

doneWithStruct8zgensym_189e87a53e58dbf2_9:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft8zgensym_189e87a53e58dbf2_9 > 0 || missingFieldsLeft8zgensym_189e87a53e58dbf2_9 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft8zgensym_189e87a53e58dbf2_9, missingFieldsLeft8zgensym_189e87a53e58dbf2_9, msgp.ShowFound(found8zgensym_189e87a53e58dbf2_9[:]), unmarshalMsgFieldOrder8zgensym_189e87a53e58dbf2_9)
		if encodedFieldsLeft8zgensym_189e87a53e58dbf2_9 > 0 {
			encodedFieldsLeft8zgensym_189e87a53e58dbf2_9--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField8zgensym_189e87a53e58dbf2_9 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss8zgensym_189e87a53e58dbf2_9 < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss8zgensym_189e87a53e58dbf2_9 = 0
			}
			for nextMiss8zgensym_189e87a53e58dbf2_9 < maxFields8zgensym_189e87a53e58dbf2_9 && (found8zgensym_189e87a53e58dbf2_9[nextMiss8zgensym_189e87a53e58dbf2_9] || unmarshalMsgFieldSkip8zgensym_189e87a53e58dbf2_9[nextMiss8zgensym_189e87a53e58dbf2_9]) {
				nextMiss8zgensym_189e87a53e58dbf2_9++
			}
			if nextMiss8zgensym_189e87a53e58dbf2_9 == maxFields8zgensym_189e87a53e58dbf2_9 {
				// filled all the empty fields!
				break doneWithStruct8zgensym_189e87a53e58dbf2_9
			}
			missingFieldsLeft8zgensym_189e87a53e58dbf2_9--
			curField8zgensym_189e87a53e58dbf2_9 = unmarshalMsgFieldOrder8zgensym_189e87a53e58dbf2_9[nextMiss8zgensym_189e87a53e58dbf2_9]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField8zgensym_189e87a53e58dbf2_9)
		switch curField8zgensym_189e87a53e58dbf2_9 {
		// -- templateUnmarshalMsg ends here --

		case "UserHomePrefix__str":
			found8zgensym_189e87a53e58dbf2_9[0] = true
			z.UserHomePrefix, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Persist__rct":
			found8zgensym_189e87a53e58dbf2_9[2] = true
			const maxFields10zgensym_189e87a53e58dbf2_11 = 2

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

				case "Users_zid00_ptr":
					found10zgensym_189e87a53e58dbf2_11[0] = true
					if nbs.AlwaysNil {
						if z.Persist.Users != nil {
							z.Persist.Users.UnmarshalMsg(msgp.OnlyNilSlice)
						}
					} else {
						// not nbs.AlwaysNil
						if msgp.IsNil(bts) {
							bts = bts[1:]
							if nil != z.Persist.Users {
								z.Persist.Users.UnmarshalMsg(msgp.OnlyNilSlice)
							}
						} else {
							// not nbs.AlwaysNil and not IsNil(bts): have something to read

							if z.Persist.Users == nil {
								z.Persist.Users = new(AtomicUserMap)
							}
							bts, err = z.Persist.Users.UnmarshalMsg(bts)
							if err != nil {
								return
							}
							if err != nil {
								return
							}
						}
					}
				case "HostPrivateKeyPath_zid01_str":
					found10zgensym_189e87a53e58dbf2_11[1] = true
					z.Persist.HostPrivateKeyPath, bts, err = nbs.ReadStringBytes(bts)

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

		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	if nextMiss8zgensym_189e87a53e58dbf2_9 != -1 {
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
var unmarshalMsgFieldOrder8zgensym_189e87a53e58dbf2_9 = []string{"UserHomePrefix__str", "", "Persist__rct"}

var unmarshalMsgFieldSkip8zgensym_189e87a53e58dbf2_9 = []bool{false, true, false}

// fields of HostDbPersist
var unmarshalMsgFieldOrder10zgensym_189e87a53e58dbf2_11 = []string{"Users_zid00_ptr", "HostPrivateKeyPath_zid01_str"}

var unmarshalMsgFieldSkip10zgensym_189e87a53e58dbf2_11 = []bool{false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *HostDb) Msgsize() (s int) {
	s = 1 + 20 + msgp.StringPrefixSize + len(z.UserHomePrefix) + 13 + 1 + 16
	if z.Persist.Users == nil {
		s += msgp.NilSize
	} else {
		s += z.Persist.Users.Msgsize()
	}
	s += 29 + msgp.StringPrefixSize + len(z.Persist.HostPrivateKeyPath)
	return
}

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *HostDbPersist) DecodeMsg(dc *msgp.Reader) (err error) {
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
	const maxFields12zgensym_189e87a53e58dbf2_13 = 2

	// -- templateDecodeMsg starts here--
	var totalEncodedFields12zgensym_189e87a53e58dbf2_13 uint32
	totalEncodedFields12zgensym_189e87a53e58dbf2_13, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft12zgensym_189e87a53e58dbf2_13 := totalEncodedFields12zgensym_189e87a53e58dbf2_13
	missingFieldsLeft12zgensym_189e87a53e58dbf2_13 := maxFields12zgensym_189e87a53e58dbf2_13 - totalEncodedFields12zgensym_189e87a53e58dbf2_13

	var nextMiss12zgensym_189e87a53e58dbf2_13 int32 = -1
	var found12zgensym_189e87a53e58dbf2_13 [maxFields12zgensym_189e87a53e58dbf2_13]bool
	var curField12zgensym_189e87a53e58dbf2_13 string

doneWithStruct12zgensym_189e87a53e58dbf2_13:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft12zgensym_189e87a53e58dbf2_13 > 0 || missingFieldsLeft12zgensym_189e87a53e58dbf2_13 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft12zgensym_189e87a53e58dbf2_13, missingFieldsLeft12zgensym_189e87a53e58dbf2_13, msgp.ShowFound(found12zgensym_189e87a53e58dbf2_13[:]), decodeMsgFieldOrder12zgensym_189e87a53e58dbf2_13)
		if encodedFieldsLeft12zgensym_189e87a53e58dbf2_13 > 0 {
			encodedFieldsLeft12zgensym_189e87a53e58dbf2_13--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField12zgensym_189e87a53e58dbf2_13 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss12zgensym_189e87a53e58dbf2_13 < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss12zgensym_189e87a53e58dbf2_13 = 0
			}
			for nextMiss12zgensym_189e87a53e58dbf2_13 < maxFields12zgensym_189e87a53e58dbf2_13 && (found12zgensym_189e87a53e58dbf2_13[nextMiss12zgensym_189e87a53e58dbf2_13] || decodeMsgFieldSkip12zgensym_189e87a53e58dbf2_13[nextMiss12zgensym_189e87a53e58dbf2_13]) {
				nextMiss12zgensym_189e87a53e58dbf2_13++
			}
			if nextMiss12zgensym_189e87a53e58dbf2_13 == maxFields12zgensym_189e87a53e58dbf2_13 {
				// filled all the empty fields!
				break doneWithStruct12zgensym_189e87a53e58dbf2_13
			}
			missingFieldsLeft12zgensym_189e87a53e58dbf2_13--
			curField12zgensym_189e87a53e58dbf2_13 = decodeMsgFieldOrder12zgensym_189e87a53e58dbf2_13[nextMiss12zgensym_189e87a53e58dbf2_13]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField12zgensym_189e87a53e58dbf2_13)
		switch curField12zgensym_189e87a53e58dbf2_13 {
		// -- templateDecodeMsg ends here --

		case "Users_zid00_ptr":
			found12zgensym_189e87a53e58dbf2_13[0] = true
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
		case "HostPrivateKeyPath_zid01_str":
			found12zgensym_189e87a53e58dbf2_13[1] = true
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
	if nextMiss12zgensym_189e87a53e58dbf2_13 != -1 {
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

// fields of HostDbPersist
var decodeMsgFieldOrder12zgensym_189e87a53e58dbf2_13 = []string{"Users_zid00_ptr", "HostPrivateKeyPath_zid01_str"}

var decodeMsgFieldSkip12zgensym_189e87a53e58dbf2_13 = []bool{false, false}

// fieldsNotEmpty supports omitempty tags
func (z *HostDbPersist) fieldsNotEmpty(isempty []bool) uint32 {
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
func (z *HostDbPersist) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// honor the omitempty tags
	var empty_zgensym_189e87a53e58dbf2_14 [2]bool
	fieldsInUse_zgensym_189e87a53e58dbf2_15 := z.fieldsNotEmpty(empty_zgensym_189e87a53e58dbf2_14[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zgensym_189e87a53e58dbf2_15)
	if err != nil {
		return err
	}

	if !empty_zgensym_189e87a53e58dbf2_14[0] {
		// write "Users_zid00_ptr"
		err = en.Append(0xaf, 0x55, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x30, 0x5f, 0x70, 0x74, 0x72)
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

	if !empty_zgensym_189e87a53e58dbf2_14[1] {
		// write "HostPrivateKeyPath_zid01_str"
		err = en.Append(0xbc, 0x48, 0x6f, 0x73, 0x74, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x31, 0x5f, 0x73, 0x74, 0x72)
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
func (z *HostDbPersist) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())

	// honor the omitempty tags
	var empty [2]bool
	fieldsInUse := z.fieldsNotEmpty(empty[:])
	o = msgp.AppendMapHeader(o, fieldsInUse)

	if !empty[0] {
		// string "Users_zid00_ptr"
		o = append(o, 0xaf, 0x55, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x30, 0x5f, 0x70, 0x74, 0x72)
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
		// string "HostPrivateKeyPath_zid01_str"
		o = append(o, 0xbc, 0x48, 0x6f, 0x73, 0x74, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x31, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.HostPrivateKeyPath)
	}

	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *HostDbPersist) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *HostDbPersist) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	var field []byte
	_ = field
	const maxFields16zgensym_189e87a53e58dbf2_17 = 2

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields16zgensym_189e87a53e58dbf2_17 uint32
	if !nbs.AlwaysNil {
		totalEncodedFields16zgensym_189e87a53e58dbf2_17, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft16zgensym_189e87a53e58dbf2_17 := totalEncodedFields16zgensym_189e87a53e58dbf2_17
	missingFieldsLeft16zgensym_189e87a53e58dbf2_17 := maxFields16zgensym_189e87a53e58dbf2_17 - totalEncodedFields16zgensym_189e87a53e58dbf2_17

	var nextMiss16zgensym_189e87a53e58dbf2_17 int32 = -1
	var found16zgensym_189e87a53e58dbf2_17 [maxFields16zgensym_189e87a53e58dbf2_17]bool
	var curField16zgensym_189e87a53e58dbf2_17 string

doneWithStruct16zgensym_189e87a53e58dbf2_17:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft16zgensym_189e87a53e58dbf2_17 > 0 || missingFieldsLeft16zgensym_189e87a53e58dbf2_17 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft16zgensym_189e87a53e58dbf2_17, missingFieldsLeft16zgensym_189e87a53e58dbf2_17, msgp.ShowFound(found16zgensym_189e87a53e58dbf2_17[:]), unmarshalMsgFieldOrder16zgensym_189e87a53e58dbf2_17)
		if encodedFieldsLeft16zgensym_189e87a53e58dbf2_17 > 0 {
			encodedFieldsLeft16zgensym_189e87a53e58dbf2_17--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField16zgensym_189e87a53e58dbf2_17 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss16zgensym_189e87a53e58dbf2_17 < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss16zgensym_189e87a53e58dbf2_17 = 0
			}
			for nextMiss16zgensym_189e87a53e58dbf2_17 < maxFields16zgensym_189e87a53e58dbf2_17 && (found16zgensym_189e87a53e58dbf2_17[nextMiss16zgensym_189e87a53e58dbf2_17] || unmarshalMsgFieldSkip16zgensym_189e87a53e58dbf2_17[nextMiss16zgensym_189e87a53e58dbf2_17]) {
				nextMiss16zgensym_189e87a53e58dbf2_17++
			}
			if nextMiss16zgensym_189e87a53e58dbf2_17 == maxFields16zgensym_189e87a53e58dbf2_17 {
				// filled all the empty fields!
				break doneWithStruct16zgensym_189e87a53e58dbf2_17
			}
			missingFieldsLeft16zgensym_189e87a53e58dbf2_17--
			curField16zgensym_189e87a53e58dbf2_17 = unmarshalMsgFieldOrder16zgensym_189e87a53e58dbf2_17[nextMiss16zgensym_189e87a53e58dbf2_17]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField16zgensym_189e87a53e58dbf2_17)
		switch curField16zgensym_189e87a53e58dbf2_17 {
		// -- templateUnmarshalMsg ends here --

		case "Users_zid00_ptr":
			found16zgensym_189e87a53e58dbf2_17[0] = true
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
		case "HostPrivateKeyPath_zid01_str":
			found16zgensym_189e87a53e58dbf2_17[1] = true
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
	if nextMiss16zgensym_189e87a53e58dbf2_17 != -1 {
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

// fields of HostDbPersist
var unmarshalMsgFieldOrder16zgensym_189e87a53e58dbf2_17 = []string{"Users_zid00_ptr", "HostPrivateKeyPath_zid01_str"}

var unmarshalMsgFieldSkip16zgensym_189e87a53e58dbf2_17 = []bool{false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *HostDbPersist) Msgsize() (s int) {
	s = 1 + 16
	if z.Users == nil {
		s += msgp.NilSize
	} else {
		s += z.Users.Msgsize()
	}
	s += 29 + msgp.StringPrefixSize + len(z.HostPrivateKeyPath)
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
	const maxFields18zgensym_189e87a53e58dbf2_19 = 5

	// -- templateDecodeMsg starts here--
	var totalEncodedFields18zgensym_189e87a53e58dbf2_19 uint32
	totalEncodedFields18zgensym_189e87a53e58dbf2_19, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft18zgensym_189e87a53e58dbf2_19 := totalEncodedFields18zgensym_189e87a53e58dbf2_19
	missingFieldsLeft18zgensym_189e87a53e58dbf2_19 := maxFields18zgensym_189e87a53e58dbf2_19 - totalEncodedFields18zgensym_189e87a53e58dbf2_19

	var nextMiss18zgensym_189e87a53e58dbf2_19 int32 = -1
	var found18zgensym_189e87a53e58dbf2_19 [maxFields18zgensym_189e87a53e58dbf2_19]bool
	var curField18zgensym_189e87a53e58dbf2_19 string

doneWithStruct18zgensym_189e87a53e58dbf2_19:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft18zgensym_189e87a53e58dbf2_19 > 0 || missingFieldsLeft18zgensym_189e87a53e58dbf2_19 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft18zgensym_189e87a53e58dbf2_19, missingFieldsLeft18zgensym_189e87a53e58dbf2_19, msgp.ShowFound(found18zgensym_189e87a53e58dbf2_19[:]), decodeMsgFieldOrder18zgensym_189e87a53e58dbf2_19)
		if encodedFieldsLeft18zgensym_189e87a53e58dbf2_19 > 0 {
			encodedFieldsLeft18zgensym_189e87a53e58dbf2_19--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField18zgensym_189e87a53e58dbf2_19 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss18zgensym_189e87a53e58dbf2_19 < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss18zgensym_189e87a53e58dbf2_19 = 0
			}
			for nextMiss18zgensym_189e87a53e58dbf2_19 < maxFields18zgensym_189e87a53e58dbf2_19 && (found18zgensym_189e87a53e58dbf2_19[nextMiss18zgensym_189e87a53e58dbf2_19] || decodeMsgFieldSkip18zgensym_189e87a53e58dbf2_19[nextMiss18zgensym_189e87a53e58dbf2_19]) {
				nextMiss18zgensym_189e87a53e58dbf2_19++
			}
			if nextMiss18zgensym_189e87a53e58dbf2_19 == maxFields18zgensym_189e87a53e58dbf2_19 {
				// filled all the empty fields!
				break doneWithStruct18zgensym_189e87a53e58dbf2_19
			}
			missingFieldsLeft18zgensym_189e87a53e58dbf2_19--
			curField18zgensym_189e87a53e58dbf2_19 = decodeMsgFieldOrder18zgensym_189e87a53e58dbf2_19[nextMiss18zgensym_189e87a53e58dbf2_19]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField18zgensym_189e87a53e58dbf2_19)
		switch curField18zgensym_189e87a53e58dbf2_19 {
		// -- templateDecodeMsg ends here --

		case "FirstTm__tim":
			found18zgensym_189e87a53e58dbf2_19[0] = true
			z.FirstTm, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "LastTm__tim":
			found18zgensym_189e87a53e58dbf2_19[1] = true
			z.LastTm, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "SeenCount__i64":
			found18zgensym_189e87a53e58dbf2_19[2] = true
			z.SeenCount, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "AcceptedCount__i64":
			found18zgensym_189e87a53e58dbf2_19[3] = true
			z.AcceptedCount, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "PubFinger__str":
			found18zgensym_189e87a53e58dbf2_19[4] = true
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
	if nextMiss18zgensym_189e87a53e58dbf2_19 != -1 {
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
var decodeMsgFieldOrder18zgensym_189e87a53e58dbf2_19 = []string{"FirstTm__tim", "LastTm__tim", "SeenCount__i64", "AcceptedCount__i64", "PubFinger__str"}

var decodeMsgFieldSkip18zgensym_189e87a53e58dbf2_19 = []bool{false, false, false, false, false}

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
	var empty_zgensym_189e87a53e58dbf2_20 [5]bool
	fieldsInUse_zgensym_189e87a53e58dbf2_21 := z.fieldsNotEmpty(empty_zgensym_189e87a53e58dbf2_20[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zgensym_189e87a53e58dbf2_21)
	if err != nil {
		return err
	}

	if !empty_zgensym_189e87a53e58dbf2_20[0] {
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

	if !empty_zgensym_189e87a53e58dbf2_20[1] {
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

	if !empty_zgensym_189e87a53e58dbf2_20[2] {
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

	if !empty_zgensym_189e87a53e58dbf2_20[3] {
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

	if !empty_zgensym_189e87a53e58dbf2_20[4] {
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
	const maxFields22zgensym_189e87a53e58dbf2_23 = 5

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields22zgensym_189e87a53e58dbf2_23 uint32
	if !nbs.AlwaysNil {
		totalEncodedFields22zgensym_189e87a53e58dbf2_23, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft22zgensym_189e87a53e58dbf2_23 := totalEncodedFields22zgensym_189e87a53e58dbf2_23
	missingFieldsLeft22zgensym_189e87a53e58dbf2_23 := maxFields22zgensym_189e87a53e58dbf2_23 - totalEncodedFields22zgensym_189e87a53e58dbf2_23

	var nextMiss22zgensym_189e87a53e58dbf2_23 int32 = -1
	var found22zgensym_189e87a53e58dbf2_23 [maxFields22zgensym_189e87a53e58dbf2_23]bool
	var curField22zgensym_189e87a53e58dbf2_23 string

doneWithStruct22zgensym_189e87a53e58dbf2_23:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft22zgensym_189e87a53e58dbf2_23 > 0 || missingFieldsLeft22zgensym_189e87a53e58dbf2_23 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft22zgensym_189e87a53e58dbf2_23, missingFieldsLeft22zgensym_189e87a53e58dbf2_23, msgp.ShowFound(found22zgensym_189e87a53e58dbf2_23[:]), unmarshalMsgFieldOrder22zgensym_189e87a53e58dbf2_23)
		if encodedFieldsLeft22zgensym_189e87a53e58dbf2_23 > 0 {
			encodedFieldsLeft22zgensym_189e87a53e58dbf2_23--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField22zgensym_189e87a53e58dbf2_23 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss22zgensym_189e87a53e58dbf2_23 < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss22zgensym_189e87a53e58dbf2_23 = 0
			}
			for nextMiss22zgensym_189e87a53e58dbf2_23 < maxFields22zgensym_189e87a53e58dbf2_23 && (found22zgensym_189e87a53e58dbf2_23[nextMiss22zgensym_189e87a53e58dbf2_23] || unmarshalMsgFieldSkip22zgensym_189e87a53e58dbf2_23[nextMiss22zgensym_189e87a53e58dbf2_23]) {
				nextMiss22zgensym_189e87a53e58dbf2_23++
			}
			if nextMiss22zgensym_189e87a53e58dbf2_23 == maxFields22zgensym_189e87a53e58dbf2_23 {
				// filled all the empty fields!
				break doneWithStruct22zgensym_189e87a53e58dbf2_23
			}
			missingFieldsLeft22zgensym_189e87a53e58dbf2_23--
			curField22zgensym_189e87a53e58dbf2_23 = unmarshalMsgFieldOrder22zgensym_189e87a53e58dbf2_23[nextMiss22zgensym_189e87a53e58dbf2_23]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField22zgensym_189e87a53e58dbf2_23)
		switch curField22zgensym_189e87a53e58dbf2_23 {
		// -- templateUnmarshalMsg ends here --

		case "FirstTm__tim":
			found22zgensym_189e87a53e58dbf2_23[0] = true
			z.FirstTm, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "LastTm__tim":
			found22zgensym_189e87a53e58dbf2_23[1] = true
			z.LastTm, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "SeenCount__i64":
			found22zgensym_189e87a53e58dbf2_23[2] = true
			z.SeenCount, bts, err = nbs.ReadInt64Bytes(bts)

			if err != nil {
				return
			}
		case "AcceptedCount__i64":
			found22zgensym_189e87a53e58dbf2_23[3] = true
			z.AcceptedCount, bts, err = nbs.ReadInt64Bytes(bts)

			if err != nil {
				return
			}
		case "PubFinger__str":
			found22zgensym_189e87a53e58dbf2_23[4] = true
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
	if nextMiss22zgensym_189e87a53e58dbf2_23 != -1 {
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
var unmarshalMsgFieldOrder22zgensym_189e87a53e58dbf2_23 = []string{"FirstTm__tim", "LastTm__tim", "SeenCount__i64", "AcceptedCount__i64", "PubFinger__str"}

var unmarshalMsgFieldSkip22zgensym_189e87a53e58dbf2_23 = []bool{false, false, false, false, false}

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
	const maxFields27zgensym_189e87a53e58dbf2_28 = 18

	// -- templateDecodeMsg starts here--
	var totalEncodedFields27zgensym_189e87a53e58dbf2_28 uint32
	totalEncodedFields27zgensym_189e87a53e58dbf2_28, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft27zgensym_189e87a53e58dbf2_28 := totalEncodedFields27zgensym_189e87a53e58dbf2_28
	missingFieldsLeft27zgensym_189e87a53e58dbf2_28 := maxFields27zgensym_189e87a53e58dbf2_28 - totalEncodedFields27zgensym_189e87a53e58dbf2_28

	var nextMiss27zgensym_189e87a53e58dbf2_28 int32 = -1
	var found27zgensym_189e87a53e58dbf2_28 [maxFields27zgensym_189e87a53e58dbf2_28]bool
	var curField27zgensym_189e87a53e58dbf2_28 string

doneWithStruct27zgensym_189e87a53e58dbf2_28:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft27zgensym_189e87a53e58dbf2_28 > 0 || missingFieldsLeft27zgensym_189e87a53e58dbf2_28 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft27zgensym_189e87a53e58dbf2_28, missingFieldsLeft27zgensym_189e87a53e58dbf2_28, msgp.ShowFound(found27zgensym_189e87a53e58dbf2_28[:]), decodeMsgFieldOrder27zgensym_189e87a53e58dbf2_28)
		if encodedFieldsLeft27zgensym_189e87a53e58dbf2_28 > 0 {
			encodedFieldsLeft27zgensym_189e87a53e58dbf2_28--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField27zgensym_189e87a53e58dbf2_28 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss27zgensym_189e87a53e58dbf2_28 < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss27zgensym_189e87a53e58dbf2_28 = 0
			}
			for nextMiss27zgensym_189e87a53e58dbf2_28 < maxFields27zgensym_189e87a53e58dbf2_28 && (found27zgensym_189e87a53e58dbf2_28[nextMiss27zgensym_189e87a53e58dbf2_28] || decodeMsgFieldSkip27zgensym_189e87a53e58dbf2_28[nextMiss27zgensym_189e87a53e58dbf2_28]) {
				nextMiss27zgensym_189e87a53e58dbf2_28++
			}
			if nextMiss27zgensym_189e87a53e58dbf2_28 == maxFields27zgensym_189e87a53e58dbf2_28 {
				// filled all the empty fields!
				break doneWithStruct27zgensym_189e87a53e58dbf2_28
			}
			missingFieldsLeft27zgensym_189e87a53e58dbf2_28--
			curField27zgensym_189e87a53e58dbf2_28 = decodeMsgFieldOrder27zgensym_189e87a53e58dbf2_28[nextMiss27zgensym_189e87a53e58dbf2_28]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField27zgensym_189e87a53e58dbf2_28)
		switch curField27zgensym_189e87a53e58dbf2_28 {
		// -- templateDecodeMsg ends here --

		case "MyEmail__str":
			found27zgensym_189e87a53e58dbf2_28[0] = true
			z.MyEmail, err = dc.ReadString()
			if err != nil {
				return
			}
		case "MyFullname__str":
			found27zgensym_189e87a53e58dbf2_28[1] = true
			z.MyFullname, err = dc.ReadString()
			if err != nil {
				return
			}
		case "MyLogin__str":
			found27zgensym_189e87a53e58dbf2_28[2] = true
			z.MyLogin, err = dc.ReadString()
			if err != nil {
				return
			}
		case "PublicKeyPath__str":
			found27zgensym_189e87a53e58dbf2_28[3] = true
			z.PublicKeyPath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "PrivateKeyPath__str":
			found27zgensym_189e87a53e58dbf2_28[4] = true
			z.PrivateKeyPath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "TOTPpath__str":
			found27zgensym_189e87a53e58dbf2_28[5] = true
			z.TOTPpath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "QrPath__str":
			found27zgensym_189e87a53e58dbf2_28[6] = true
			z.QrPath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Issuer__str":
			found27zgensym_189e87a53e58dbf2_28[7] = true
			z.Issuer, err = dc.ReadString()
			if err != nil {
				return
			}
		case "SeenPubKey__map":
			found27zgensym_189e87a53e58dbf2_28[9] = true
			var zgensym_189e87a53e58dbf2_29 uint32
			zgensym_189e87a53e58dbf2_29, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.SeenPubKey == nil && zgensym_189e87a53e58dbf2_29 > 0 {
				z.SeenPubKey = make(map[string]LoginRecord, zgensym_189e87a53e58dbf2_29)
			} else if len(z.SeenPubKey) > 0 {
				for key, _ := range z.SeenPubKey {
					delete(z.SeenPubKey, key)
				}
			}
			for zgensym_189e87a53e58dbf2_29 > 0 {
				zgensym_189e87a53e58dbf2_29--
				var zgensym_189e87a53e58dbf2_24 string
				var zgensym_189e87a53e58dbf2_25 LoginRecord
				zgensym_189e87a53e58dbf2_24, err = dc.ReadString()
				if err != nil {
					return
				}
				err = zgensym_189e87a53e58dbf2_25.DecodeMsg(dc)
				if err != nil {
					return
				}
				z.SeenPubKey[zgensym_189e87a53e58dbf2_24] = zgensym_189e87a53e58dbf2_25
			}
		case "ScryptedPassword__bin":
			found27zgensym_189e87a53e58dbf2_28[10] = true
			z.ScryptedPassword, err = dc.ReadBytes(z.ScryptedPassword)
			if err != nil {
				return
			}
		case "ClearPw__str":
			found27zgensym_189e87a53e58dbf2_28[11] = true
			z.ClearPw, err = dc.ReadString()
			if err != nil {
				return
			}
		case "TOTPorig__str":
			found27zgensym_189e87a53e58dbf2_28[12] = true
			z.TOTPorig, err = dc.ReadString()
			if err != nil {
				return
			}
		case "FirstLoginTime__tim":
			found27zgensym_189e87a53e58dbf2_28[13] = true
			z.FirstLoginTime, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "LastLoginTime__tim":
			found27zgensym_189e87a53e58dbf2_28[14] = true
			z.LastLoginTime, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "LastLoginAddr__str":
			found27zgensym_189e87a53e58dbf2_28[15] = true
			z.LastLoginAddr, err = dc.ReadString()
			if err != nil {
				return
			}
		case "IPwhitelist__slc":
			found27zgensym_189e87a53e58dbf2_28[16] = true
			var zgensym_189e87a53e58dbf2_30 uint32
			zgensym_189e87a53e58dbf2_30, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.IPwhitelist) >= int(zgensym_189e87a53e58dbf2_30) {
				z.IPwhitelist = (z.IPwhitelist)[:zgensym_189e87a53e58dbf2_30]
			} else {
				z.IPwhitelist = make([]string, zgensym_189e87a53e58dbf2_30)
			}
			for zgensym_189e87a53e58dbf2_26 := range z.IPwhitelist {
				z.IPwhitelist[zgensym_189e87a53e58dbf2_26], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "DisabledAcct__boo":
			found27zgensym_189e87a53e58dbf2_28[17] = true
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
	if nextMiss27zgensym_189e87a53e58dbf2_28 != -1 {
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
var decodeMsgFieldOrder27zgensym_189e87a53e58dbf2_28 = []string{"MyEmail__str", "MyFullname__str", "MyLogin__str", "PublicKeyPath__str", "PrivateKeyPath__str", "TOTPpath__str", "QrPath__str", "Issuer__str", "", "SeenPubKey__map", "ScryptedPassword__bin", "ClearPw__str", "TOTPorig__str", "FirstLoginTime__tim", "LastLoginTime__tim", "LastLoginAddr__str", "IPwhitelist__slc", "DisabledAcct__boo"}

var decodeMsgFieldSkip27zgensym_189e87a53e58dbf2_28 = []bool{false, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false, false, false}

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
	isempty[9] = (len(z.SeenPubKey) == 0) // string, omitempty
	if isempty[9] {
		fieldsInUse--
	}
	isempty[10] = (len(z.ScryptedPassword) == 0) // string, omitempty
	if isempty[10] {
		fieldsInUse--
	}
	isempty[11] = (len(z.ClearPw) == 0) // string, omitempty
	if isempty[11] {
		fieldsInUse--
	}
	isempty[12] = (len(z.TOTPorig) == 0) // string, omitempty
	if isempty[12] {
		fieldsInUse--
	}
	isempty[13] = (z.FirstLoginTime.IsZero()) // time.Time, omitempty
	if isempty[13] {
		fieldsInUse--
	}
	isempty[14] = (z.LastLoginTime.IsZero()) // time.Time, omitempty
	if isempty[14] {
		fieldsInUse--
	}
	isempty[15] = (len(z.LastLoginAddr) == 0) // string, omitempty
	if isempty[15] {
		fieldsInUse--
	}
	isempty[16] = (len(z.IPwhitelist) == 0) // string, omitempty
	if isempty[16] {
		fieldsInUse--
	}
	isempty[17] = (!z.DisabledAcct) // bool, omitempty
	if isempty[17] {
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
	var empty_zgensym_189e87a53e58dbf2_31 [18]bool
	fieldsInUse_zgensym_189e87a53e58dbf2_32 := z.fieldsNotEmpty(empty_zgensym_189e87a53e58dbf2_31[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zgensym_189e87a53e58dbf2_32)
	if err != nil {
		return err
	}

	if !empty_zgensym_189e87a53e58dbf2_31[0] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[1] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[2] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[3] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[4] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[5] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[6] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[7] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[9] {
		// write "SeenPubKey__map"
		err = en.Append(0xaf, 0x53, 0x65, 0x65, 0x6e, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x5f, 0x5f, 0x6d, 0x61, 0x70)
		if err != nil {
			return err
		}
		err = en.WriteMapHeader(uint32(len(z.SeenPubKey)))
		if err != nil {
			return
		}
		for zgensym_189e87a53e58dbf2_24, zgensym_189e87a53e58dbf2_25 := range z.SeenPubKey {
			err = en.WriteString(zgensym_189e87a53e58dbf2_24)
			if err != nil {
				return
			}
			err = zgensym_189e87a53e58dbf2_25.EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_31[10] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[11] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[12] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[13] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[14] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[15] {
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

	if !empty_zgensym_189e87a53e58dbf2_31[16] {
		// write "IPwhitelist__slc"
		err = en.Append(0xb0, 0x49, 0x50, 0x77, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x5f, 0x73, 0x6c, 0x63)
		if err != nil {
			return err
		}
		err = en.WriteArrayHeader(uint32(len(z.IPwhitelist)))
		if err != nil {
			return
		}
		for zgensym_189e87a53e58dbf2_26 := range z.IPwhitelist {
			err = en.WriteString(z.IPwhitelist[zgensym_189e87a53e58dbf2_26])
			if err != nil {
				return
			}
		}
	}

	if !empty_zgensym_189e87a53e58dbf2_31[17] {
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
	var empty [18]bool
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

	if !empty[9] {
		// string "SeenPubKey__map"
		o = append(o, 0xaf, 0x53, 0x65, 0x65, 0x6e, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x5f, 0x5f, 0x6d, 0x61, 0x70)
		o = msgp.AppendMapHeader(o, uint32(len(z.SeenPubKey)))
		for zgensym_189e87a53e58dbf2_24, zgensym_189e87a53e58dbf2_25 := range z.SeenPubKey {
			o = msgp.AppendString(o, zgensym_189e87a53e58dbf2_24)
			o, err = zgensym_189e87a53e58dbf2_25.MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}

	if !empty[10] {
		// string "ScryptedPassword__bin"
		o = append(o, 0xb5, 0x53, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x5f, 0x62, 0x69, 0x6e)
		o = msgp.AppendBytes(o, z.ScryptedPassword)
	}

	if !empty[11] {
		// string "ClearPw__str"
		o = append(o, 0xac, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x50, 0x77, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.ClearPw)
	}

	if !empty[12] {
		// string "TOTPorig__str"
		o = append(o, 0xad, 0x54, 0x4f, 0x54, 0x50, 0x6f, 0x72, 0x69, 0x67, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.TOTPorig)
	}

	if !empty[13] {
		// string "FirstLoginTime__tim"
		o = append(o, 0xb3, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x5f, 0x5f, 0x74, 0x69, 0x6d)
		o = msgp.AppendTime(o, z.FirstLoginTime)
	}

	if !empty[14] {
		// string "LastLoginTime__tim"
		o = append(o, 0xb2, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x5f, 0x5f, 0x74, 0x69, 0x6d)
		o = msgp.AppendTime(o, z.LastLoginTime)
	}

	if !empty[15] {
		// string "LastLoginAddr__str"
		o = append(o, 0xb2, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x5f, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.LastLoginAddr)
	}

	if !empty[16] {
		// string "IPwhitelist__slc"
		o = append(o, 0xb0, 0x49, 0x50, 0x77, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x5f, 0x73, 0x6c, 0x63)
		o = msgp.AppendArrayHeader(o, uint32(len(z.IPwhitelist)))
		for zgensym_189e87a53e58dbf2_26 := range z.IPwhitelist {
			o = msgp.AppendString(o, z.IPwhitelist[zgensym_189e87a53e58dbf2_26])
		}
	}

	if !empty[17] {
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
	const maxFields33zgensym_189e87a53e58dbf2_34 = 18

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields33zgensym_189e87a53e58dbf2_34 uint32
	if !nbs.AlwaysNil {
		totalEncodedFields33zgensym_189e87a53e58dbf2_34, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft33zgensym_189e87a53e58dbf2_34 := totalEncodedFields33zgensym_189e87a53e58dbf2_34
	missingFieldsLeft33zgensym_189e87a53e58dbf2_34 := maxFields33zgensym_189e87a53e58dbf2_34 - totalEncodedFields33zgensym_189e87a53e58dbf2_34

	var nextMiss33zgensym_189e87a53e58dbf2_34 int32 = -1
	var found33zgensym_189e87a53e58dbf2_34 [maxFields33zgensym_189e87a53e58dbf2_34]bool
	var curField33zgensym_189e87a53e58dbf2_34 string

doneWithStruct33zgensym_189e87a53e58dbf2_34:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft33zgensym_189e87a53e58dbf2_34 > 0 || missingFieldsLeft33zgensym_189e87a53e58dbf2_34 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft33zgensym_189e87a53e58dbf2_34, missingFieldsLeft33zgensym_189e87a53e58dbf2_34, msgp.ShowFound(found33zgensym_189e87a53e58dbf2_34[:]), unmarshalMsgFieldOrder33zgensym_189e87a53e58dbf2_34)
		if encodedFieldsLeft33zgensym_189e87a53e58dbf2_34 > 0 {
			encodedFieldsLeft33zgensym_189e87a53e58dbf2_34--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField33zgensym_189e87a53e58dbf2_34 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss33zgensym_189e87a53e58dbf2_34 < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss33zgensym_189e87a53e58dbf2_34 = 0
			}
			for nextMiss33zgensym_189e87a53e58dbf2_34 < maxFields33zgensym_189e87a53e58dbf2_34 && (found33zgensym_189e87a53e58dbf2_34[nextMiss33zgensym_189e87a53e58dbf2_34] || unmarshalMsgFieldSkip33zgensym_189e87a53e58dbf2_34[nextMiss33zgensym_189e87a53e58dbf2_34]) {
				nextMiss33zgensym_189e87a53e58dbf2_34++
			}
			if nextMiss33zgensym_189e87a53e58dbf2_34 == maxFields33zgensym_189e87a53e58dbf2_34 {
				// filled all the empty fields!
				break doneWithStruct33zgensym_189e87a53e58dbf2_34
			}
			missingFieldsLeft33zgensym_189e87a53e58dbf2_34--
			curField33zgensym_189e87a53e58dbf2_34 = unmarshalMsgFieldOrder33zgensym_189e87a53e58dbf2_34[nextMiss33zgensym_189e87a53e58dbf2_34]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField33zgensym_189e87a53e58dbf2_34)
		switch curField33zgensym_189e87a53e58dbf2_34 {
		// -- templateUnmarshalMsg ends here --

		case "MyEmail__str":
			found33zgensym_189e87a53e58dbf2_34[0] = true
			z.MyEmail, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "MyFullname__str":
			found33zgensym_189e87a53e58dbf2_34[1] = true
			z.MyFullname, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "MyLogin__str":
			found33zgensym_189e87a53e58dbf2_34[2] = true
			z.MyLogin, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "PublicKeyPath__str":
			found33zgensym_189e87a53e58dbf2_34[3] = true
			z.PublicKeyPath, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "PrivateKeyPath__str":
			found33zgensym_189e87a53e58dbf2_34[4] = true
			z.PrivateKeyPath, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "TOTPpath__str":
			found33zgensym_189e87a53e58dbf2_34[5] = true
			z.TOTPpath, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "QrPath__str":
			found33zgensym_189e87a53e58dbf2_34[6] = true
			z.QrPath, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Issuer__str":
			found33zgensym_189e87a53e58dbf2_34[7] = true
			z.Issuer, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "SeenPubKey__map":
			found33zgensym_189e87a53e58dbf2_34[9] = true
			if nbs.AlwaysNil {
				if len(z.SeenPubKey) > 0 {
					for key, _ := range z.SeenPubKey {
						delete(z.SeenPubKey, key)
					}
				}

			} else {

				var zgensym_189e87a53e58dbf2_35 uint32
				zgensym_189e87a53e58dbf2_35, bts, err = nbs.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				if z.SeenPubKey == nil && zgensym_189e87a53e58dbf2_35 > 0 {
					z.SeenPubKey = make(map[string]LoginRecord, zgensym_189e87a53e58dbf2_35)
				} else if len(z.SeenPubKey) > 0 {
					for key, _ := range z.SeenPubKey {
						delete(z.SeenPubKey, key)
					}
				}
				for zgensym_189e87a53e58dbf2_35 > 0 {
					var zgensym_189e87a53e58dbf2_24 string
					var zgensym_189e87a53e58dbf2_25 LoginRecord
					zgensym_189e87a53e58dbf2_35--
					zgensym_189e87a53e58dbf2_24, bts, err = nbs.ReadStringBytes(bts)
					if err != nil {
						return
					}
					bts, err = zgensym_189e87a53e58dbf2_25.UnmarshalMsg(bts)
					if err != nil {
						return
					}
					if err != nil {
						return
					}
					z.SeenPubKey[zgensym_189e87a53e58dbf2_24] = zgensym_189e87a53e58dbf2_25
				}
			}
		case "ScryptedPassword__bin":
			found33zgensym_189e87a53e58dbf2_34[10] = true
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
			found33zgensym_189e87a53e58dbf2_34[11] = true
			z.ClearPw, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "TOTPorig__str":
			found33zgensym_189e87a53e58dbf2_34[12] = true
			z.TOTPorig, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "FirstLoginTime__tim":
			found33zgensym_189e87a53e58dbf2_34[13] = true
			z.FirstLoginTime, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "LastLoginTime__tim":
			found33zgensym_189e87a53e58dbf2_34[14] = true
			z.LastLoginTime, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "LastLoginAddr__str":
			found33zgensym_189e87a53e58dbf2_34[15] = true
			z.LastLoginAddr, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "IPwhitelist__slc":
			found33zgensym_189e87a53e58dbf2_34[16] = true
			if nbs.AlwaysNil {
				(z.IPwhitelist) = (z.IPwhitelist)[:0]
			} else {

				var zgensym_189e87a53e58dbf2_36 uint32
				zgensym_189e87a53e58dbf2_36, bts, err = nbs.ReadArrayHeaderBytes(bts)
				if err != nil {
					return
				}
				if cap(z.IPwhitelist) >= int(zgensym_189e87a53e58dbf2_36) {
					z.IPwhitelist = (z.IPwhitelist)[:zgensym_189e87a53e58dbf2_36]
				} else {
					z.IPwhitelist = make([]string, zgensym_189e87a53e58dbf2_36)
				}
				for zgensym_189e87a53e58dbf2_26 := range z.IPwhitelist {
					z.IPwhitelist[zgensym_189e87a53e58dbf2_26], bts, err = nbs.ReadStringBytes(bts)

					if err != nil {
						return
					}
				}
			}
		case "DisabledAcct__boo":
			found33zgensym_189e87a53e58dbf2_34[17] = true
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
	if nextMiss33zgensym_189e87a53e58dbf2_34 != -1 {
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
var unmarshalMsgFieldOrder33zgensym_189e87a53e58dbf2_34 = []string{"MyEmail__str", "MyFullname__str", "MyLogin__str", "PublicKeyPath__str", "PrivateKeyPath__str", "TOTPpath__str", "QrPath__str", "Issuer__str", "", "SeenPubKey__map", "ScryptedPassword__bin", "ClearPw__str", "TOTPorig__str", "FirstLoginTime__tim", "LastLoginTime__tim", "LastLoginAddr__str", "IPwhitelist__slc", "DisabledAcct__boo"}

var unmarshalMsgFieldSkip33zgensym_189e87a53e58dbf2_34 = []bool{false, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false, false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *User) Msgsize() (s int) {
	s = 3 + 13 + msgp.StringPrefixSize + len(z.MyEmail) + 16 + msgp.StringPrefixSize + len(z.MyFullname) + 13 + msgp.StringPrefixSize + len(z.MyLogin) + 19 + msgp.StringPrefixSize + len(z.PublicKeyPath) + 20 + msgp.StringPrefixSize + len(z.PrivateKeyPath) + 14 + msgp.StringPrefixSize + len(z.TOTPpath) + 12 + msgp.StringPrefixSize + len(z.QrPath) + 12 + msgp.StringPrefixSize + len(z.Issuer) + 16 + msgp.MapHeaderSize
	if z.SeenPubKey != nil {
		for zgensym_189e87a53e58dbf2_24, zgensym_189e87a53e58dbf2_25 := range z.SeenPubKey {
			_ = zgensym_189e87a53e58dbf2_25
			_ = zgensym_189e87a53e58dbf2_24
			s += msgp.StringPrefixSize + len(zgensym_189e87a53e58dbf2_24) + zgensym_189e87a53e58dbf2_25.Msgsize()
		}
	}
	s += 22 + msgp.BytesPrefixSize + len(z.ScryptedPassword) + 13 + msgp.StringPrefixSize + len(z.ClearPw) + 14 + msgp.StringPrefixSize + len(z.TOTPorig) + 20 + msgp.TimeSize + 19 + msgp.TimeSize + 19 + msgp.StringPrefixSize + len(z.LastLoginAddr) + 17 + msgp.ArrayHeaderSize
	for zgensym_189e87a53e58dbf2_26 := range z.IPwhitelist {
		s += msgp.StringPrefixSize + len(z.IPwhitelist[zgensym_189e87a53e58dbf2_26])
	}
	s += 18 + msgp.BoolSize
	return
}
