package jtt

import "fmt"

// T808_0x0100 终端注册
//
// 按协议版本（2011/2013/2019）字段长度与填充规则不同：
//   - 2011:
//     ManufacturerID 固定 5 字节（右补 0x00）
//     TerminalModel 最长 8 字节（不足右补空格 0x20）
//     TerminalID 固定 7 字节（右补 0x00，字符集必须为 [A-Z0-9]）
//   - 2013:
//     ManufacturerID 固定 5 字节（右补 0x00）
//     TerminalModel 最长 20 字节（不足右补 0x00）
//     TerminalID 最长 7 字节（不足右补 0x00，字符集必须为 [A-Z0-9]）
//   - 2019:
//     ManufacturerID 固定 11 字节（右补 0x00）
//     TerminalModel 最长 30 字节（不足左补 0x00）
//     TerminalID 固定 30 字节（右补 0x00，字符集必须为 [A-Z0-9]）
//
// 说明：
//   - 所有字符串以 GB18030 编码计算长度与写入；Decode 时按对应规则去除补位字节。
//   - PlateColor 为 1 字节；PlateNumber 为尾随可变长 GB18030 字符串。
type T808_0x0100 struct {
	// 省域ID，标示终端安装车辆所在的省域，0 保留，由平台取默认值。
	// 省域 ID 采用 GB/T 2260 中规定的行政区划代码六位中前两位。
	ProvinceID uint16 `json:"provinceID"`
	// 市县域ID，标示终端安装车辆所在的市域和县域，0 保留，由平台取默认值。
	// 市县域 ID 采用 GB/T 2260 中规定的行政区划代码六位中后四位。
	CityID uint16 `json:"cityID"`
	// 制造商ID
	//	- 2011 版本，终端制造商编码，长度=5，不足/超过报错
	//	- 2013 版本，终端制造商编码，长度=5，不足/超过报错
	//	- 2019 版本，由车载终端厂商所在地行政区划代码和制造商ID组成，长度=11，不足/超过报错
	ManufacturerID string `json:"manufacturerID"`
	// 终端型号
	//	- 2011 版本，长度<=8，不足右补空格 0x20，超过报错
	//	- 2013 版本，长度<=20，不足右补 0x00，超过报错
	//	- 2019 版本，长度<=30，不足左补 0x00，超过报错
	TerminalModel string `json:"terminalModel"`
	// 终端ID
	//	- 2011 版本，由大写字母和数字组成，长度=7，不足/超过报错
	//	- 2013 版本，由大写字母和数字组成，长度<=7，不足右补 0x00，超过报错
	//	- 2019 版本，由大写字母和数字组成，长度=30，不足/超过报错
	TerminalID string `json:"terminalID"`
	// 车牌颜色
	// 	- 2011/2013 版本，按照 JT/T 415-2006 的 5.4.12
	// 	- 2019 版本，按照 JT/T 697.7-2014 中的规定，未上牌车辆填 0
	PlateColor byte `json:"plateColor"`
	// 车牌 STRING，公安交通管理部门颁发的机动车号牌
	//	- 2013 版本，车牌颜色为 0 时，表示车辆 VIN
	//	- 2019 版本，如果车辆未上牌则填写车架号
	PlateNumber string `json:"plateNumber"`

	// 协议版本号：-1=2011, 0=2013(默认), 1=2019
	// 若未显式设置，Encode/Decode 均按 2013 处理。
	protocolVersion VersionType
}

func (entity *T808_0x0100) MsgID() MsgID { return MsgT808_0x0100 }

// ProtocolVersion 返回当前消息体所采用的协议版本。
func (entity *T808_0x0100) ProtocolVersion() VersionType { return entity.protocolVersion }

// SetProtocolVersion 设置该消息体的协议版本。
func (entity *T808_0x0100) SetProtocolVersion(protocolVersion VersionType) {
	entity.protocolVersion = protocolVersion
}

func (entity *T808_0x0100) Encode() ([]byte, error) {
	writer := NewWriter()

	// WORD fields
	writer.WriteWord(entity.ProvinceID)
	writer.WriteWord(entity.CityID)

	// Version-aware fixed-length fields
	ver := entity.protocolVersion
	if ver == 0 { // default old version
		ver = Version2013
	}

	// helper to get GB18030 byte length
	gbLen := func(s string) (int, error) { return GB18030Length(s) }

	switch ver {
	case Version2011:
		// ManufacturerID BYTE[5]
		if n, err := gbLen(entity.ManufacturerID); err != nil {
			return nil, fmt.Errorf("encode 0x0100 manufacturerID(2011) gb18030 err: %w", err)
		} else if n != 5 {
			return nil, fmt.Errorf("encode 0x0100 manufacturerID length invalid for 2011: got %d bytes, require 5", n)
		}
		if err := writer.WriteString(entity.ManufacturerID, 5); err != nil {
			return nil, err
		}

		// TerminalModel BYTE[8], pad with spaces 0x20 on the right
		if n, err := gbLen(entity.TerminalModel); err != nil {
			return nil, fmt.Errorf("encode 0x0100 terminalModel(2011) gb18030 err: %w", err)
		} else if n > 8 {
			return nil, fmt.Errorf("encode 0x0100 terminalModel too long for 2011: got %d bytes, limit 8", n)
		}
		// manually space-pad to 8
		if err := writer.WriteString(entity.TerminalModel); err != nil {
			return nil, err
		}
		if n, _ := gbLen(entity.TerminalModel); n < 8 {
			for i := 0; i < 8-n; i++ {
				writer.Write([]byte{0x20})
			}
		}

		// TerminalID BYTE[7] (right-pad 0x00)
		if n, err := gbLen(entity.TerminalID); err != nil {
			return nil, fmt.Errorf("encode 0x0100 terminalID(2011) gb18030 err: %w", err)
		} else if n != 7 {
			return nil, fmt.Errorf("encode 0x0100 terminalID length invalid for 2011: got %d bytes, require 7", n)
		}
		if err := writer.WriteString(entity.TerminalID, 7); err != nil {
			return nil, err
		}

	case Version2013:
		// ManufacturerID BYTE[5]
		if n, err := gbLen(entity.ManufacturerID); err != nil {
			return nil, fmt.Errorf("encode 0x0100 manufacturerID(2013) gb18030 err: %w", err)
		} else if n != 5 {
			return nil, fmt.Errorf("encode 0x0100 manufacturerID length invalid for 2013: got %d bytes, require 5", n)
		}
		if err := writer.WriteString(entity.ManufacturerID, 5); err != nil {
			return nil, err
		}

		// TerminalModel BYTE[20], right-pad 0x00
		if n, err := gbLen(entity.TerminalModel); err != nil {
			return nil, fmt.Errorf("encode 0x0100 terminalModel(2013) gb18030 err: %w", err)
		} else if n > 20 {
			return nil, fmt.Errorf("encode 0x0100 terminalModel too long for 2013: got %d bytes, limit 20", n)
		}
		if err := writer.WriteString(entity.TerminalModel, 20); err != nil {
			return nil, err
		}

		// TerminalID BYTE[7]
		if n, err := gbLen(entity.TerminalID); err != nil {
			return nil, fmt.Errorf("encode 0x0100 terminalID(2013) gb18030 err: %w", err)
		} else if n > 7 {
			return nil, fmt.Errorf("encode 0x0100 terminalID too long for 2013: got %d bytes, limit 7", n)
		}
		if err := writer.WriteString(entity.TerminalID, 7); err != nil {
			return nil, err
		}

	case Version2019:
		// ManufacturerID BYTE[11]
		if n, err := gbLen(entity.ManufacturerID); err != nil {
			return nil, fmt.Errorf("encode 0x0100 manufacturerID(2019) gb18030 err: %w", err)
		} else if n != 11 {
			return nil, fmt.Errorf("encode 0x0100 manufacturerID length invalid for 2019: got %d bytes, require 11", n)
		}
		if err := writer.WriteString(entity.ManufacturerID, 11); err != nil {
			return nil, err
		}

		// TerminalModel BYTE[30], left-pad 0x00
		if n, err := gbLen(entity.TerminalModel); err != nil {
			return nil, fmt.Errorf("encode 0x0100 terminalModel(2019) gb18030 err: %w", err)
		} else if n > 30 {
			return nil, fmt.Errorf("encode 0x0100 terminalModel too long for 2019: got %d bytes, limit 30", n)
		}
		// 说明：writer.WriteString 仅支持右侧 0x00 填充，无法满足 2019 左侧 0x00 的需求，
		// 因此这里通过临时 Writer 获取 GB18030 编码后的原始字节，再手动在左侧补 0x00。
		tmp := NewWriter()
		if err := tmp.WriteString(entity.TerminalModel); err != nil {
			return nil, err
		}
		modelBytes := tmp.Bytes()
		pad := 30 - len(modelBytes)
		// 左侧补 0x00 直至 30 字节
		buf := make([]byte, pad)
		buf = append(buf, modelBytes...)
		writer.Write(buf)

		// TerminalID BYTE[30] right-pad 0x00
		if n, err := gbLen(entity.TerminalID); err != nil {
			return nil, fmt.Errorf("encode 0x0100 terminalID(2019) gb18030 err: %w", err)
		} else if n != 30 {
			return nil, fmt.Errorf("encode 0x0100 terminalID length invalid for 2019: got %d bytes, require 30", n)
		}
		if err := writer.WriteString(entity.TerminalID, 30); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("encode 0x0100: unsupported protocolVersion: %d", ver)
	}

	// Plate color
	writer.WriteByte(entity.PlateColor)

	// Vehicle plate STRING (GB18030)
	if err := writer.WriteString(entity.PlateNumber); err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x0100) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	if entity.ProvinceID, err = reader.ReadWord(); err != nil {
		return 0, fmt.Errorf("read ProvinceID: %w", err)
	}
	if entity.CityID, err = reader.ReadWord(); err != nil {
		return 0, fmt.Errorf("read CityID: %w", err)
	}

	ver := entity.protocolVersion
	if ver == 0 {
		ver = Version2013
	}

	var b []byte
	switch ver {
	case Version2011:
		if b, err = reader.Read(5); err != nil {
			return 0, fmt.Errorf("read ManufacturerID(2011): %w", err)
		}
		entity.ManufacturerID = string(trimRightZeros(b))
		if b, err = reader.Read(8); err != nil {
			return 0, fmt.Errorf("read TerminalModel(2011): %w", err)
		}
		entity.TerminalModel = string(trimRightSpaces(b))
		if b, err = reader.Read(7); err != nil {
			return 0, fmt.Errorf("read TerminalID(2011): %w", err)
		}
		entity.TerminalID = string(trimRightZeros(b))

	case Version2013:
		if b, err = reader.Read(5); err != nil {
			return 0, fmt.Errorf("read ManufacturerID(2013): %w", err)
		}
		entity.ManufacturerID = BytesToString(b)
		if b, err = reader.Read(20); err != nil {
			return 0, fmt.Errorf("read TerminalModel(2013): %w", err)
		}
		entity.TerminalModel = BytesToString(b)
		if b, err = reader.Read(7); err != nil {
			return 0, fmt.Errorf("read TerminalID(2013): %w", err)
		}
		entity.TerminalID = BytesToString(b)

	case Version2019:
		if b, err = reader.Read(11); err != nil {
			return 0, fmt.Errorf("read ManufacturerID(2019): %w", err)
		}
		entity.ManufacturerID = BytesToString(b)
		if b, err = reader.Read(30); err != nil {
			return 0, fmt.Errorf("read TerminalModel(2019): %w", err)
		}
		b = trimLeftZeros(b)
		entity.TerminalModel = string(b)
		if b, err = reader.Read(30); err != nil {
			return 0, fmt.Errorf("read TerminalID(2019): %w", err)
		}
		entity.TerminalID = BytesToString(b)
	default:
		return 0, fmt.Errorf("decode 0x0100: unsupported protocolVersion: %d", ver)
	}

	if entity.PlateColor, err = reader.ReadByte(); err != nil {
		return 0, fmt.Errorf("read PlateColor: %w", err)
	}

	// Remaining is STRING for PlateNumber (may be empty)
	if entity.PlateNumber, err = reader.ReadString(); err != nil {
		return 0, fmt.Errorf("read PlateNumber: %w", err)
	}
	return len(data) - reader.Len(), nil
}
