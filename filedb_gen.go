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
	const maxFields2zgensym_995050f2db4fcb57_3 = 2

	// -- templateDecodeMsg starts here--
	var totalEncodedFields2zgensym_995050f2db4fcb57_3 uint32
	totalEncodedFields2zgensym_995050f2db4fcb57_3, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft2zgensym_995050f2db4fcb57_3 := totalEncodedFields2zgensym_995050f2db4fcb57_3
	missingFieldsLeft2zgensym_995050f2db4fcb57_3 := maxFields2zgensym_995050f2db4fcb57_3 - totalEncodedFields2zgensym_995050f2db4fcb57_3

	var nextMiss2zgensym_995050f2db4fcb57_3 int32 = -1
	var found2zgensym_995050f2db4fcb57_3 [maxFields2zgensym_995050f2db4fcb57_3]bool
	var curField2zgensym_995050f2db4fcb57_3 string

doneWithStruct2zgensym_995050f2db4fcb57_3:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft2zgensym_995050f2db4fcb57_3 > 0 || missingFieldsLeft2zgensym_995050f2db4fcb57_3 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft2zgensym_995050f2db4fcb57_3, missingFieldsLeft2zgensym_995050f2db4fcb57_3, msgp.ShowFound(found2zgensym_995050f2db4fcb57_3[:]), decodeMsgFieldOrder2zgensym_995050f2db4fcb57_3)
		if encodedFieldsLeft2zgensym_995050f2db4fcb57_3 > 0 {
			encodedFieldsLeft2zgensym_995050f2db4fcb57_3--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField2zgensym_995050f2db4fcb57_3 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss2zgensym_995050f2db4fcb57_3 < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss2zgensym_995050f2db4fcb57_3 = 0
			}
			for nextMiss2zgensym_995050f2db4fcb57_3 < maxFields2zgensym_995050f2db4fcb57_3 && (found2zgensym_995050f2db4fcb57_3[nextMiss2zgensym_995050f2db4fcb57_3] || decodeMsgFieldSkip2zgensym_995050f2db4fcb57_3[nextMiss2zgensym_995050f2db4fcb57_3]) {
				nextMiss2zgensym_995050f2db4fcb57_3++
			}
			if nextMiss2zgensym_995050f2db4fcb57_3 == maxFields2zgensym_995050f2db4fcb57_3 {
				// filled all the empty fields!
				break doneWithStruct2zgensym_995050f2db4fcb57_3
			}
			missingFieldsLeft2zgensym_995050f2db4fcb57_3--
			curField2zgensym_995050f2db4fcb57_3 = decodeMsgFieldOrder2zgensym_995050f2db4fcb57_3[nextMiss2zgensym_995050f2db4fcb57_3]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField2zgensym_995050f2db4fcb57_3)
		switch curField2zgensym_995050f2db4fcb57_3 {
		// -- templateDecodeMsg ends here --

		case "Filepath_zid00_str":
			found2zgensym_995050f2db4fcb57_3[0] = true
			z.Filepath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Map_zid01_map":
			found2zgensym_995050f2db4fcb57_3[1] = true
			var zgensym_995050f2db4fcb57_4 uint32
			zgensym_995050f2db4fcb57_4, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Map == nil && zgensym_995050f2db4fcb57_4 > 0 {
				z.Map = make(map[string]string, zgensym_995050f2db4fcb57_4)
			} else if len(z.Map) > 0 {
				for key, _ := range z.Map {
					delete(z.Map, key)
				}
			}
			for zgensym_995050f2db4fcb57_4 > 0 {
				zgensym_995050f2db4fcb57_4--
				var zgensym_995050f2db4fcb57_0 string
				var zgensym_995050f2db4fcb57_1 string
				zgensym_995050f2db4fcb57_0, err = dc.ReadString()
				if err != nil {
					return
				}
				zgensym_995050f2db4fcb57_1, err = dc.ReadString()
				if err != nil {
					return
				}
				z.Map[zgensym_995050f2db4fcb57_0] = zgensym_995050f2db4fcb57_1
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	if nextMiss2zgensym_995050f2db4fcb57_3 != -1 {
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
var decodeMsgFieldOrder2zgensym_995050f2db4fcb57_3 = []string{"Filepath_zid00_str", "Map_zid01_map"}

var decodeMsgFieldSkip2zgensym_995050f2db4fcb57_3 = []bool{false, false}

// fieldsNotEmpty supports omitempty tags
func (z *Filedb) fieldsNotEmpty(isempty []bool) uint32 {
	if len(isempty) == 0 {
		return 2
	}
	var fieldsInUse uint32 = 2
	isempty[0] = (len(z.Filepath) == 0) // string, omitempty
	if isempty[0] {
		fieldsInUse--
	}
	isempty[1] = (len(z.Map) == 0) // string, omitempty
	if isempty[1] {
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
	var empty_zgensym_995050f2db4fcb57_5 [2]bool
	fieldsInUse_zgensym_995050f2db4fcb57_6 := z.fieldsNotEmpty(empty_zgensym_995050f2db4fcb57_5[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zgensym_995050f2db4fcb57_6)
	if err != nil {
		return err
	}

	if !empty_zgensym_995050f2db4fcb57_5[0] {
		// write "Filepath_zid00_str"
		err = en.Append(0xb2, 0x46, 0x69, 0x6c, 0x65, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x30, 0x5f, 0x73, 0x74, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteString(z.Filepath)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_995050f2db4fcb57_5[1] {
		// write "Map_zid01_map"
		err = en.Append(0xad, 0x4d, 0x61, 0x70, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x31, 0x5f, 0x6d, 0x61, 0x70)
		if err != nil {
			return err
		}
		err = en.WriteMapHeader(uint32(len(z.Map)))
		if err != nil {
			return
		}
		for zgensym_995050f2db4fcb57_0, zgensym_995050f2db4fcb57_1 := range z.Map {
			err = en.WriteString(zgensym_995050f2db4fcb57_0)
			if err != nil {
				return
			}
			err = en.WriteString(zgensym_995050f2db4fcb57_1)
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
	var empty [2]bool
	fieldsInUse := z.fieldsNotEmpty(empty[:])
	o = msgp.AppendMapHeader(o, fieldsInUse)

	if !empty[0] {
		// string "Filepath_zid00_str"
		o = append(o, 0xb2, 0x46, 0x69, 0x6c, 0x65, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x30, 0x5f, 0x73, 0x74, 0x72)
		o = msgp.AppendString(o, z.Filepath)
	}

	if !empty[1] {
		// string "Map_zid01_map"
		o = append(o, 0xad, 0x4d, 0x61, 0x70, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x31, 0x5f, 0x6d, 0x61, 0x70)
		o = msgp.AppendMapHeader(o, uint32(len(z.Map)))
		for zgensym_995050f2db4fcb57_0, zgensym_995050f2db4fcb57_1 := range z.Map {
			o = msgp.AppendString(o, zgensym_995050f2db4fcb57_0)
			o = msgp.AppendString(o, zgensym_995050f2db4fcb57_1)
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
	const maxFields7zgensym_995050f2db4fcb57_8 = 2

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields7zgensym_995050f2db4fcb57_8 uint32
	if !nbs.AlwaysNil {
		totalEncodedFields7zgensym_995050f2db4fcb57_8, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft7zgensym_995050f2db4fcb57_8 := totalEncodedFields7zgensym_995050f2db4fcb57_8
	missingFieldsLeft7zgensym_995050f2db4fcb57_8 := maxFields7zgensym_995050f2db4fcb57_8 - totalEncodedFields7zgensym_995050f2db4fcb57_8

	var nextMiss7zgensym_995050f2db4fcb57_8 int32 = -1
	var found7zgensym_995050f2db4fcb57_8 [maxFields7zgensym_995050f2db4fcb57_8]bool
	var curField7zgensym_995050f2db4fcb57_8 string

doneWithStruct7zgensym_995050f2db4fcb57_8:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft7zgensym_995050f2db4fcb57_8 > 0 || missingFieldsLeft7zgensym_995050f2db4fcb57_8 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft7zgensym_995050f2db4fcb57_8, missingFieldsLeft7zgensym_995050f2db4fcb57_8, msgp.ShowFound(found7zgensym_995050f2db4fcb57_8[:]), unmarshalMsgFieldOrder7zgensym_995050f2db4fcb57_8)
		if encodedFieldsLeft7zgensym_995050f2db4fcb57_8 > 0 {
			encodedFieldsLeft7zgensym_995050f2db4fcb57_8--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField7zgensym_995050f2db4fcb57_8 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss7zgensym_995050f2db4fcb57_8 < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss7zgensym_995050f2db4fcb57_8 = 0
			}
			for nextMiss7zgensym_995050f2db4fcb57_8 < maxFields7zgensym_995050f2db4fcb57_8 && (found7zgensym_995050f2db4fcb57_8[nextMiss7zgensym_995050f2db4fcb57_8] || unmarshalMsgFieldSkip7zgensym_995050f2db4fcb57_8[nextMiss7zgensym_995050f2db4fcb57_8]) {
				nextMiss7zgensym_995050f2db4fcb57_8++
			}
			if nextMiss7zgensym_995050f2db4fcb57_8 == maxFields7zgensym_995050f2db4fcb57_8 {
				// filled all the empty fields!
				break doneWithStruct7zgensym_995050f2db4fcb57_8
			}
			missingFieldsLeft7zgensym_995050f2db4fcb57_8--
			curField7zgensym_995050f2db4fcb57_8 = unmarshalMsgFieldOrder7zgensym_995050f2db4fcb57_8[nextMiss7zgensym_995050f2db4fcb57_8]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField7zgensym_995050f2db4fcb57_8)
		switch curField7zgensym_995050f2db4fcb57_8 {
		// -- templateUnmarshalMsg ends here --

		case "Filepath_zid00_str":
			found7zgensym_995050f2db4fcb57_8[0] = true
			z.Filepath, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Map_zid01_map":
			found7zgensym_995050f2db4fcb57_8[1] = true
			if nbs.AlwaysNil {
				if len(z.Map) > 0 {
					for key, _ := range z.Map {
						delete(z.Map, key)
					}
				}

			} else {

				var zgensym_995050f2db4fcb57_9 uint32
				zgensym_995050f2db4fcb57_9, bts, err = nbs.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				if z.Map == nil && zgensym_995050f2db4fcb57_9 > 0 {
					z.Map = make(map[string]string, zgensym_995050f2db4fcb57_9)
				} else if len(z.Map) > 0 {
					for key, _ := range z.Map {
						delete(z.Map, key)
					}
				}
				for zgensym_995050f2db4fcb57_9 > 0 {
					var zgensym_995050f2db4fcb57_0 string
					var zgensym_995050f2db4fcb57_1 string
					zgensym_995050f2db4fcb57_9--
					zgensym_995050f2db4fcb57_0, bts, err = nbs.ReadStringBytes(bts)
					if err != nil {
						return
					}
					zgensym_995050f2db4fcb57_1, bts, err = nbs.ReadStringBytes(bts)

					if err != nil {
						return
					}
					z.Map[zgensym_995050f2db4fcb57_0] = zgensym_995050f2db4fcb57_1
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	if nextMiss7zgensym_995050f2db4fcb57_8 != -1 {
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
var unmarshalMsgFieldOrder7zgensym_995050f2db4fcb57_8 = []string{"Filepath_zid00_str", "Map_zid01_map"}

var unmarshalMsgFieldSkip7zgensym_995050f2db4fcb57_8 = []bool{false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Filedb) Msgsize() (s int) {
	s = 1 + 19 + msgp.StringPrefixSize + len(z.Filepath) + 14 + msgp.MapHeaderSize
	if z.Map != nil {
		for zgensym_995050f2db4fcb57_0, zgensym_995050f2db4fcb57_1 := range z.Map {
			_ = zgensym_995050f2db4fcb57_1
			_ = zgensym_995050f2db4fcb57_0
			s += msgp.StringPrefixSize + len(zgensym_995050f2db4fcb57_0) + msgp.StringPrefixSize + len(zgensym_995050f2db4fcb57_1)
		}
	}
	return
}
