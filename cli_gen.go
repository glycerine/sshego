package sshego

// NOTE: THIS FILE WAS PRODUCED BY THE
// GREENPACK CODE GENERATION TOOL (github.com/glycerine/greenpack)
// DO NOT EDIT

import (
	"github.com/glycerine/greenpack/msgp"
)

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *KeepAlivePing) DecodeMsg(dc *msgp.Reader) (err error) {
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
	const maxFields0zgensym_c523013f9c573deb_1 = 3

	// -- templateDecodeMsg starts here--
	var totalEncodedFields0zgensym_c523013f9c573deb_1 uint32
	totalEncodedFields0zgensym_c523013f9c573deb_1, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft0zgensym_c523013f9c573deb_1 := totalEncodedFields0zgensym_c523013f9c573deb_1
	missingFieldsLeft0zgensym_c523013f9c573deb_1 := maxFields0zgensym_c523013f9c573deb_1 - totalEncodedFields0zgensym_c523013f9c573deb_1

	var nextMiss0zgensym_c523013f9c573deb_1 int32 = -1
	var found0zgensym_c523013f9c573deb_1 [maxFields0zgensym_c523013f9c573deb_1]bool
	var curField0zgensym_c523013f9c573deb_1 string

doneWithStruct0zgensym_c523013f9c573deb_1:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft0zgensym_c523013f9c573deb_1 > 0 || missingFieldsLeft0zgensym_c523013f9c573deb_1 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft0zgensym_c523013f9c573deb_1, missingFieldsLeft0zgensym_c523013f9c573deb_1, msgp.ShowFound(found0zgensym_c523013f9c573deb_1[:]), decodeMsgFieldOrder0zgensym_c523013f9c573deb_1)
		if encodedFieldsLeft0zgensym_c523013f9c573deb_1 > 0 {
			encodedFieldsLeft0zgensym_c523013f9c573deb_1--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField0zgensym_c523013f9c573deb_1 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss0zgensym_c523013f9c573deb_1 < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss0zgensym_c523013f9c573deb_1 = 0
			}
			for nextMiss0zgensym_c523013f9c573deb_1 < maxFields0zgensym_c523013f9c573deb_1 && (found0zgensym_c523013f9c573deb_1[nextMiss0zgensym_c523013f9c573deb_1] || decodeMsgFieldSkip0zgensym_c523013f9c573deb_1[nextMiss0zgensym_c523013f9c573deb_1]) {
				nextMiss0zgensym_c523013f9c573deb_1++
			}
			if nextMiss0zgensym_c523013f9c573deb_1 == maxFields0zgensym_c523013f9c573deb_1 {
				// filled all the empty fields!
				break doneWithStruct0zgensym_c523013f9c573deb_1
			}
			missingFieldsLeft0zgensym_c523013f9c573deb_1--
			curField0zgensym_c523013f9c573deb_1 = decodeMsgFieldOrder0zgensym_c523013f9c573deb_1[nextMiss0zgensym_c523013f9c573deb_1]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField0zgensym_c523013f9c573deb_1)
		switch curField0zgensym_c523013f9c573deb_1 {
		// -- templateDecodeMsg ends here --

		case "Sent_zid00_tim":
			found0zgensym_c523013f9c573deb_1[0] = true
			z.Sent, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "Replied_zid01_tim":
			found0zgensym_c523013f9c573deb_1[1] = true
			z.Replied, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "Serial_zid02_i64":
			found0zgensym_c523013f9c573deb_1[2] = true
			z.Serial, err = dc.ReadInt64()
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
	if nextMiss0zgensym_c523013f9c573deb_1 != -1 {
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

// fields of KeepAlivePing
var decodeMsgFieldOrder0zgensym_c523013f9c573deb_1 = []string{"Sent_zid00_tim", "Replied_zid01_tim", "Serial_zid02_i64"}

var decodeMsgFieldSkip0zgensym_c523013f9c573deb_1 = []bool{false, false, false}

// fieldsNotEmpty supports omitempty tags
func (z KeepAlivePing) fieldsNotEmpty(isempty []bool) uint32 {
	if len(isempty) == 0 {
		return 3
	}
	var fieldsInUse uint32 = 3
	isempty[0] = (z.Sent.IsZero()) // time.Time, omitempty
	if isempty[0] {
		fieldsInUse--
	}
	isempty[1] = (z.Replied.IsZero()) // time.Time, omitempty
	if isempty[1] {
		fieldsInUse--
	}
	isempty[2] = (z.Serial == 0) // number, omitempty
	if isempty[2] {
		fieldsInUse--
	}

	return fieldsInUse
}

// EncodeMsg implements msgp.Encodable
func (z KeepAlivePing) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// honor the omitempty tags
	var empty_zgensym_c523013f9c573deb_2 [3]bool
	fieldsInUse_zgensym_c523013f9c573deb_3 := z.fieldsNotEmpty(empty_zgensym_c523013f9c573deb_2[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zgensym_c523013f9c573deb_3)
	if err != nil {
		return err
	}

	if !empty_zgensym_c523013f9c573deb_2[0] {
		// write "Sent_zid00_tim"
		err = en.Append(0xae, 0x53, 0x65, 0x6e, 0x74, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x30, 0x5f, 0x74, 0x69, 0x6d)
		if err != nil {
			return err
		}
		err = en.WriteTime(z.Sent)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_c523013f9c573deb_2[1] {
		// write "Replied_zid01_tim"
		err = en.Append(0xb1, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x65, 0x64, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x31, 0x5f, 0x74, 0x69, 0x6d)
		if err != nil {
			return err
		}
		err = en.WriteTime(z.Replied)
		if err != nil {
			return
		}
	}

	if !empty_zgensym_c523013f9c573deb_2[2] {
		// write "Serial_zid02_i64"
		err = en.Append(0xb0, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x32, 0x5f, 0x69, 0x36, 0x34)
		if err != nil {
			return err
		}
		err = en.WriteInt64(z.Serial)
		if err != nil {
			return
		}
	}

	return
}

// MarshalMsg implements msgp.Marshaler
func (z KeepAlivePing) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())

	// honor the omitempty tags
	var empty [3]bool
	fieldsInUse := z.fieldsNotEmpty(empty[:])
	o = msgp.AppendMapHeader(o, fieldsInUse)

	if !empty[0] {
		// string "Sent_zid00_tim"
		o = append(o, 0xae, 0x53, 0x65, 0x6e, 0x74, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x30, 0x5f, 0x74, 0x69, 0x6d)
		o = msgp.AppendTime(o, z.Sent)
	}

	if !empty[1] {
		// string "Replied_zid01_tim"
		o = append(o, 0xb1, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x65, 0x64, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x31, 0x5f, 0x74, 0x69, 0x6d)
		o = msgp.AppendTime(o, z.Replied)
	}

	if !empty[2] {
		// string "Serial_zid02_i64"
		o = append(o, 0xb0, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x7a, 0x69, 0x64, 0x30, 0x32, 0x5f, 0x69, 0x36, 0x34)
		o = msgp.AppendInt64(o, z.Serial)
	}

	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *KeepAlivePing) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *KeepAlivePing) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	var field []byte
	_ = field
	const maxFields4zgensym_c523013f9c573deb_5 = 3

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields4zgensym_c523013f9c573deb_5 uint32
	if !nbs.AlwaysNil {
		totalEncodedFields4zgensym_c523013f9c573deb_5, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft4zgensym_c523013f9c573deb_5 := totalEncodedFields4zgensym_c523013f9c573deb_5
	missingFieldsLeft4zgensym_c523013f9c573deb_5 := maxFields4zgensym_c523013f9c573deb_5 - totalEncodedFields4zgensym_c523013f9c573deb_5

	var nextMiss4zgensym_c523013f9c573deb_5 int32 = -1
	var found4zgensym_c523013f9c573deb_5 [maxFields4zgensym_c523013f9c573deb_5]bool
	var curField4zgensym_c523013f9c573deb_5 string

doneWithStruct4zgensym_c523013f9c573deb_5:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft4zgensym_c523013f9c573deb_5 > 0 || missingFieldsLeft4zgensym_c523013f9c573deb_5 > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft4zgensym_c523013f9c573deb_5, missingFieldsLeft4zgensym_c523013f9c573deb_5, msgp.ShowFound(found4zgensym_c523013f9c573deb_5[:]), unmarshalMsgFieldOrder4zgensym_c523013f9c573deb_5)
		if encodedFieldsLeft4zgensym_c523013f9c573deb_5 > 0 {
			encodedFieldsLeft4zgensym_c523013f9c573deb_5--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField4zgensym_c523013f9c573deb_5 = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss4zgensym_c523013f9c573deb_5 < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss4zgensym_c523013f9c573deb_5 = 0
			}
			for nextMiss4zgensym_c523013f9c573deb_5 < maxFields4zgensym_c523013f9c573deb_5 && (found4zgensym_c523013f9c573deb_5[nextMiss4zgensym_c523013f9c573deb_5] || unmarshalMsgFieldSkip4zgensym_c523013f9c573deb_5[nextMiss4zgensym_c523013f9c573deb_5]) {
				nextMiss4zgensym_c523013f9c573deb_5++
			}
			if nextMiss4zgensym_c523013f9c573deb_5 == maxFields4zgensym_c523013f9c573deb_5 {
				// filled all the empty fields!
				break doneWithStruct4zgensym_c523013f9c573deb_5
			}
			missingFieldsLeft4zgensym_c523013f9c573deb_5--
			curField4zgensym_c523013f9c573deb_5 = unmarshalMsgFieldOrder4zgensym_c523013f9c573deb_5[nextMiss4zgensym_c523013f9c573deb_5]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField4zgensym_c523013f9c573deb_5)
		switch curField4zgensym_c523013f9c573deb_5 {
		// -- templateUnmarshalMsg ends here --

		case "Sent_zid00_tim":
			found4zgensym_c523013f9c573deb_5[0] = true
			z.Sent, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "Replied_zid01_tim":
			found4zgensym_c523013f9c573deb_5[1] = true
			z.Replied, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "Serial_zid02_i64":
			found4zgensym_c523013f9c573deb_5[2] = true
			z.Serial, bts, err = nbs.ReadInt64Bytes(bts)

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
	if nextMiss4zgensym_c523013f9c573deb_5 != -1 {
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

// fields of KeepAlivePing
var unmarshalMsgFieldOrder4zgensym_c523013f9c573deb_5 = []string{"Sent_zid00_tim", "Replied_zid01_tim", "Serial_zid02_i64"}

var unmarshalMsgFieldSkip4zgensym_c523013f9c573deb_5 = []bool{false, false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z KeepAlivePing) Msgsize() (s int) {
	s = 1 + 15 + msgp.TimeSize + 18 + msgp.TimeSize + 17 + msgp.Int64Size
	return
}
