package sshego

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import "github.com/tinylib/msgp/msgp"

// DecodeMsg implements msgp.Decodable
func (z *HostDb) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxvk uint32
	zxvk, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxvk > 0 {
		zxvk--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Users":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Users = nil
			} else {
				if z.Users == nil {
					z.Users = new(AtomicUserMap)
				}
				err = z.Users.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "HostPrivateKeyPath":
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
	return
}

// EncodeMsg implements msgp.Encodable
func (z *HostDb) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Users"
	err = en.Append(0x82, 0xa5, 0x55, 0x73, 0x65, 0x72, 0x73)
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
	// write "HostPrivateKeyPath"
	err = en.Append(0xb2, 0x48, 0x6f, 0x73, 0x74, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68)
	if err != nil {
		return err
	}
	err = en.WriteString(z.HostPrivateKeyPath)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *HostDb) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Users"
	o = append(o, 0x82, 0xa5, 0x55, 0x73, 0x65, 0x72, 0x73)
	if z.Users == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Users.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "HostPrivateKeyPath"
	o = append(o, 0xb2, 0x48, 0x6f, 0x73, 0x74, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68)
	o = msgp.AppendString(o, z.HostPrivateKeyPath)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *HostDb) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zbzg uint32
	zbzg, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zbzg > 0 {
		zbzg--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Users":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Users = nil
			} else {
				if z.Users == nil {
					z.Users = new(AtomicUserMap)
				}
				bts, err = z.Users.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "HostPrivateKeyPath":
			z.HostPrivateKeyPath, bts, err = msgp.ReadStringBytes(bts)
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
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *HostDb) Msgsize() (s int) {
	s = 1 + 6
	if z.Users == nil {
		s += msgp.NilSize
	} else {
		s += z.Users.Msgsize()
	}
	s += 19 + msgp.StringPrefixSize + len(z.HostPrivateKeyPath)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *LoginRecord) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "FirstTm":
			z.FirstTm, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "LastTm":
			z.LastTm, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "SeenCount":
			z.SeenCount, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "AcceptedCount":
			z.AcceptedCount, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "PubFinger":
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
	return
}

// EncodeMsg implements msgp.Encodable
func (z *LoginRecord) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "FirstTm"
	err = en.Append(0x85, 0xa7, 0x46, 0x69, 0x72, 0x73, 0x74, 0x54, 0x6d)
	if err != nil {
		return err
	}
	err = en.WriteTime(z.FirstTm)
	if err != nil {
		return
	}
	// write "LastTm"
	err = en.Append(0xa6, 0x4c, 0x61, 0x73, 0x74, 0x54, 0x6d)
	if err != nil {
		return err
	}
	err = en.WriteTime(z.LastTm)
	if err != nil {
		return
	}
	// write "SeenCount"
	err = en.Append(0xa9, 0x53, 0x65, 0x65, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt64(z.SeenCount)
	if err != nil {
		return
	}
	// write "AcceptedCount"
	err = en.Append(0xad, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt64(z.AcceptedCount)
	if err != nil {
		return
	}
	// write "PubFinger"
	err = en.Append(0xa9, 0x50, 0x75, 0x62, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteString(z.PubFinger)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *LoginRecord) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "FirstTm"
	o = append(o, 0x85, 0xa7, 0x46, 0x69, 0x72, 0x73, 0x74, 0x54, 0x6d)
	o = msgp.AppendTime(o, z.FirstTm)
	// string "LastTm"
	o = append(o, 0xa6, 0x4c, 0x61, 0x73, 0x74, 0x54, 0x6d)
	o = msgp.AppendTime(o, z.LastTm)
	// string "SeenCount"
	o = append(o, 0xa9, 0x53, 0x65, 0x65, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74)
	o = msgp.AppendInt64(o, z.SeenCount)
	// string "AcceptedCount"
	o = append(o, 0xad, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74)
	o = msgp.AppendInt64(o, z.AcceptedCount)
	// string "PubFinger"
	o = append(o, 0xa9, 0x50, 0x75, 0x62, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72)
	o = msgp.AppendString(o, z.PubFinger)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *LoginRecord) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcmr uint32
	zcmr, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcmr > 0 {
		zcmr--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "FirstTm":
			z.FirstTm, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				return
			}
		case "LastTm":
			z.LastTm, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				return
			}
		case "SeenCount":
			z.SeenCount, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "AcceptedCount":
			z.AcceptedCount, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "PubFinger":
			z.PubFinger, bts, err = msgp.ReadStringBytes(bts)
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
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *LoginRecord) Msgsize() (s int) {
	s = 1 + 8 + msgp.TimeSize + 7 + msgp.TimeSize + 10 + msgp.Int64Size + 14 + msgp.Int64Size + 10 + msgp.StringPrefixSize + len(z.PubFinger)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *User) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zcua uint32
	zcua, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zcua > 0 {
		zcua--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "MyEmail":
			z.MyEmail, err = dc.ReadString()
			if err != nil {
				return
			}
		case "MyFullname":
			z.MyFullname, err = dc.ReadString()
			if err != nil {
				return
			}
		case "MyLogin":
			z.MyLogin, err = dc.ReadString()
			if err != nil {
				return
			}
		case "PublicKeyPath":
			z.PublicKeyPath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "PrivateKeyPath":
			z.PrivateKeyPath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "TOTPpath":
			z.TOTPpath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "QrPath":
			z.QrPath, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Issuer":
			z.Issuer, err = dc.ReadString()
			if err != nil {
				return
			}
		case "SeenPubKey":
			var zxhx uint32
			zxhx, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.SeenPubKey == nil && zxhx > 0 {
				z.SeenPubKey = make(map[string]LoginRecord, zxhx)
			} else if len(z.SeenPubKey) > 0 {
				for key, _ := range z.SeenPubKey {
					delete(z.SeenPubKey, key)
				}
			}
			for zxhx > 0 {
				zxhx--
				var zajw string
				var zwht LoginRecord
				zajw, err = dc.ReadString()
				if err != nil {
					return
				}
				err = zwht.DecodeMsg(dc)
				if err != nil {
					return
				}
				z.SeenPubKey[zajw] = zwht
			}
		case "ScryptedPassword":
			z.ScryptedPassword, err = dc.ReadBytes(z.ScryptedPassword)
			if err != nil {
				return
			}
		case "ClearPw":
			z.ClearPw, err = dc.ReadString()
			if err != nil {
				return
			}
		case "TOTPorig":
			z.TOTPorig, err = dc.ReadString()
			if err != nil {
				return
			}
		case "FirstLoginTime":
			z.FirstLoginTime, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "LastLoginTime":
			z.LastLoginTime, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "LastLoginAddr":
			z.LastLoginAddr, err = dc.ReadString()
			if err != nil {
				return
			}
		case "IPwhitelist":
			var zlqf uint32
			zlqf, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.IPwhitelist) >= int(zlqf) {
				z.IPwhitelist = (z.IPwhitelist)[:zlqf]
			} else {
				z.IPwhitelist = make([]string, zlqf)
			}
			for zhct := range z.IPwhitelist {
				z.IPwhitelist[zhct], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "DisabledAcct":
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
	return
}

// EncodeMsg implements msgp.Encodable
func (z *User) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 17
	// write "MyEmail"
	err = en.Append(0xde, 0x0, 0x11, 0xa7, 0x4d, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c)
	if err != nil {
		return err
	}
	err = en.WriteString(z.MyEmail)
	if err != nil {
		return
	}
	// write "MyFullname"
	err = en.Append(0xaa, 0x4d, 0x79, 0x46, 0x75, 0x6c, 0x6c, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.MyFullname)
	if err != nil {
		return
	}
	// write "MyLogin"
	err = en.Append(0xa7, 0x4d, 0x79, 0x4c, 0x6f, 0x67, 0x69, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteString(z.MyLogin)
	if err != nil {
		return
	}
	// write "PublicKeyPath"
	err = en.Append(0xad, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68)
	if err != nil {
		return err
	}
	err = en.WriteString(z.PublicKeyPath)
	if err != nil {
		return
	}
	// write "PrivateKeyPath"
	err = en.Append(0xae, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68)
	if err != nil {
		return err
	}
	err = en.WriteString(z.PrivateKeyPath)
	if err != nil {
		return
	}
	// write "TOTPpath"
	err = en.Append(0xa8, 0x54, 0x4f, 0x54, 0x50, 0x70, 0x61, 0x74, 0x68)
	if err != nil {
		return err
	}
	err = en.WriteString(z.TOTPpath)
	if err != nil {
		return
	}
	// write "QrPath"
	err = en.Append(0xa6, 0x51, 0x72, 0x50, 0x61, 0x74, 0x68)
	if err != nil {
		return err
	}
	err = en.WriteString(z.QrPath)
	if err != nil {
		return
	}
	// write "Issuer"
	err = en.Append(0xa6, 0x49, 0x73, 0x73, 0x75, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Issuer)
	if err != nil {
		return
	}
	// write "SeenPubKey"
	err = en.Append(0xaa, 0x53, 0x65, 0x65, 0x6e, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.SeenPubKey)))
	if err != nil {
		return
	}
	for zajw, zwht := range z.SeenPubKey {
		err = en.WriteString(zajw)
		if err != nil {
			return
		}
		err = zwht.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "ScryptedPassword"
	err = en.Append(0xb0, 0x53, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.ScryptedPassword)
	if err != nil {
		return
	}
	// write "ClearPw"
	err = en.Append(0xa7, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x50, 0x77)
	if err != nil {
		return err
	}
	err = en.WriteString(z.ClearPw)
	if err != nil {
		return
	}
	// write "TOTPorig"
	err = en.Append(0xa8, 0x54, 0x4f, 0x54, 0x50, 0x6f, 0x72, 0x69, 0x67)
	if err != nil {
		return err
	}
	err = en.WriteString(z.TOTPorig)
	if err != nil {
		return
	}
	// write "FirstLoginTime"
	err = en.Append(0xae, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteTime(z.FirstLoginTime)
	if err != nil {
		return
	}
	// write "LastLoginTime"
	err = en.Append(0xad, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteTime(z.LastLoginTime)
	if err != nil {
		return
	}
	// write "LastLoginAddr"
	err = en.Append(0xad, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x41, 0x64, 0x64, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteString(z.LastLoginAddr)
	if err != nil {
		return
	}
	// write "IPwhitelist"
	err = en.Append(0xab, 0x49, 0x50, 0x77, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69, 0x73, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.IPwhitelist)))
	if err != nil {
		return
	}
	for zhct := range z.IPwhitelist {
		err = en.WriteString(z.IPwhitelist[zhct])
		if err != nil {
			return
		}
	}
	// write "DisabledAcct"
	err = en.Append(0xac, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x41, 0x63, 0x63, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteBool(z.DisabledAcct)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *User) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 17
	// string "MyEmail"
	o = append(o, 0xde, 0x0, 0x11, 0xa7, 0x4d, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c)
	o = msgp.AppendString(o, z.MyEmail)
	// string "MyFullname"
	o = append(o, 0xaa, 0x4d, 0x79, 0x46, 0x75, 0x6c, 0x6c, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.MyFullname)
	// string "MyLogin"
	o = append(o, 0xa7, 0x4d, 0x79, 0x4c, 0x6f, 0x67, 0x69, 0x6e)
	o = msgp.AppendString(o, z.MyLogin)
	// string "PublicKeyPath"
	o = append(o, 0xad, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68)
	o = msgp.AppendString(o, z.PublicKeyPath)
	// string "PrivateKeyPath"
	o = append(o, 0xae, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x74, 0x68)
	o = msgp.AppendString(o, z.PrivateKeyPath)
	// string "TOTPpath"
	o = append(o, 0xa8, 0x54, 0x4f, 0x54, 0x50, 0x70, 0x61, 0x74, 0x68)
	o = msgp.AppendString(o, z.TOTPpath)
	// string "QrPath"
	o = append(o, 0xa6, 0x51, 0x72, 0x50, 0x61, 0x74, 0x68)
	o = msgp.AppendString(o, z.QrPath)
	// string "Issuer"
	o = append(o, 0xa6, 0x49, 0x73, 0x73, 0x75, 0x65, 0x72)
	o = msgp.AppendString(o, z.Issuer)
	// string "SeenPubKey"
	o = append(o, 0xaa, 0x53, 0x65, 0x65, 0x6e, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79)
	o = msgp.AppendMapHeader(o, uint32(len(z.SeenPubKey)))
	for zajw, zwht := range z.SeenPubKey {
		o = msgp.AppendString(o, zajw)
		o, err = zwht.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "ScryptedPassword"
	o = append(o, 0xb0, 0x53, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64)
	o = msgp.AppendBytes(o, z.ScryptedPassword)
	// string "ClearPw"
	o = append(o, 0xa7, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x50, 0x77)
	o = msgp.AppendString(o, z.ClearPw)
	// string "TOTPorig"
	o = append(o, 0xa8, 0x54, 0x4f, 0x54, 0x50, 0x6f, 0x72, 0x69, 0x67)
	o = msgp.AppendString(o, z.TOTPorig)
	// string "FirstLoginTime"
	o = append(o, 0xae, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65)
	o = msgp.AppendTime(o, z.FirstLoginTime)
	// string "LastLoginTime"
	o = append(o, 0xad, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65)
	o = msgp.AppendTime(o, z.LastLoginTime)
	// string "LastLoginAddr"
	o = append(o, 0xad, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x41, 0x64, 0x64, 0x72)
	o = msgp.AppendString(o, z.LastLoginAddr)
	// string "IPwhitelist"
	o = append(o, 0xab, 0x49, 0x50, 0x77, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69, 0x73, 0x74)
	o = msgp.AppendArrayHeader(o, uint32(len(z.IPwhitelist)))
	for zhct := range z.IPwhitelist {
		o = msgp.AppendString(o, z.IPwhitelist[zhct])
	}
	// string "DisabledAcct"
	o = append(o, 0xac, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x41, 0x63, 0x63, 0x74)
	o = msgp.AppendBool(o, z.DisabledAcct)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *User) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zdaf uint32
	zdaf, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zdaf > 0 {
		zdaf--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "MyEmail":
			z.MyEmail, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "MyFullname":
			z.MyFullname, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "MyLogin":
			z.MyLogin, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "PublicKeyPath":
			z.PublicKeyPath, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "PrivateKeyPath":
			z.PrivateKeyPath, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "TOTPpath":
			z.TOTPpath, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "QrPath":
			z.QrPath, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Issuer":
			z.Issuer, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "SeenPubKey":
			var zpks uint32
			zpks, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.SeenPubKey == nil && zpks > 0 {
				z.SeenPubKey = make(map[string]LoginRecord, zpks)
			} else if len(z.SeenPubKey) > 0 {
				for key, _ := range z.SeenPubKey {
					delete(z.SeenPubKey, key)
				}
			}
			for zpks > 0 {
				var zajw string
				var zwht LoginRecord
				zpks--
				zajw, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				bts, err = zwht.UnmarshalMsg(bts)
				if err != nil {
					return
				}
				z.SeenPubKey[zajw] = zwht
			}
		case "ScryptedPassword":
			z.ScryptedPassword, bts, err = msgp.ReadBytesBytes(bts, z.ScryptedPassword)
			if err != nil {
				return
			}
		case "ClearPw":
			z.ClearPw, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "TOTPorig":
			z.TOTPorig, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "FirstLoginTime":
			z.FirstLoginTime, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				return
			}
		case "LastLoginTime":
			z.LastLoginTime, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				return
			}
		case "LastLoginAddr":
			z.LastLoginAddr, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "IPwhitelist":
			var zjfb uint32
			zjfb, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.IPwhitelist) >= int(zjfb) {
				z.IPwhitelist = (z.IPwhitelist)[:zjfb]
			} else {
				z.IPwhitelist = make([]string, zjfb)
			}
			for zhct := range z.IPwhitelist {
				z.IPwhitelist[zhct], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "DisabledAcct":
			z.DisabledAcct, bts, err = msgp.ReadBoolBytes(bts)
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
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *User) Msgsize() (s int) {
	s = 3 + 8 + msgp.StringPrefixSize + len(z.MyEmail) + 11 + msgp.StringPrefixSize + len(z.MyFullname) + 8 + msgp.StringPrefixSize + len(z.MyLogin) + 14 + msgp.StringPrefixSize + len(z.PublicKeyPath) + 15 + msgp.StringPrefixSize + len(z.PrivateKeyPath) + 9 + msgp.StringPrefixSize + len(z.TOTPpath) + 7 + msgp.StringPrefixSize + len(z.QrPath) + 7 + msgp.StringPrefixSize + len(z.Issuer) + 11 + msgp.MapHeaderSize
	if z.SeenPubKey != nil {
		for zajw, zwht := range z.SeenPubKey {
			_ = zwht
			s += msgp.StringPrefixSize + len(zajw) + zwht.Msgsize()
		}
	}
	s += 17 + msgp.BytesPrefixSize + len(z.ScryptedPassword) + 8 + msgp.StringPrefixSize + len(z.ClearPw) + 9 + msgp.StringPrefixSize + len(z.TOTPorig) + 15 + msgp.TimeSize + 14 + msgp.TimeSize + 14 + msgp.StringPrefixSize + len(z.LastLoginAddr) + 12 + msgp.ArrayHeaderSize
	for zhct := range z.IPwhitelist {
		s += msgp.StringPrefixSize + len(z.IPwhitelist[zhct])
	}
	s += 13 + msgp.BoolSize
	return
}
