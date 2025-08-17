package jtt

import "fmt"

// T808_0x0102 终端鉴权
type T808_0x0102 struct {
	// 鉴权码
	//	- 2011/2013 版本，终端重连后上报鉴权码
	//	- 2019 版本，鉴权码长度 + 鉴权码内容
	AuthCode string `json:"authCode"`
	// 终端 IMEI，BYTE[15]，2019 版本新增
	IMEI string `json:"imei"`
	// 软件版本号，BYTE[20]，2019 版本新增，厂家自定义版本号，位数不足时，后补"0x00"
	SoftwareVersion string `json:"softwareVersion"`

	// 协议版本号：-1=2011, 0=2013(默认), 1=2019
	// 若未显式设置，Encode/Decode 均按 2013 处理。
	protocolVersion VersionType
}

func (entity *T808_0x0102) MsgID() MsgID { return MsgT808_0x0102 }

// GetProtocolVersion 获取协议版本
func (entity *T808_0x0102) GetProtocolVersion() VersionType {
	return entity.protocolVersion
}

// SetProtocolVersion 设置协议版本
func (entity *T808_0x0102) SetProtocolVersion(protocolVersion VersionType) {
	entity.protocolVersion = protocolVersion
}

func (entity *T808_0x0102) Encode() ([]byte, error) {
	writer := NewWriter()

	ver := entity.protocolVersion
	if ver == 0 {
		ver = Version2013
	}

	switch ver {
	case Version2019:
		// 写入鉴权码长度 + 内容（GB18030）
		n, err := GB18030Length(entity.AuthCode)
		if err != nil {
			return nil, fmt.Errorf("encode 0x0102 authCode gb18030 err: %w", err)
		}
		if n > 255 {
			return nil, fmt.Errorf("encode 0x0102 authCode too long: %d", n)
		}
		writer.WriteByte(byte(n))
		if n > 0 {
			if err := writer.WriteString(entity.AuthCode); err != nil {
				return nil, err
			}
		}
		// IMEI BYTE[15]
		if err := writer.WriteString(entity.IMEI, 15); err != nil {
			return nil, err
		}
		// 软件版本 BYTE[20]
		if err := writer.WriteString(entity.SoftwareVersion, 20); err != nil {
			return nil, err
		}
	default: // 2011/2013: 仅鉴权码 STRING
		if len(entity.AuthCode) > 0 {
			if err := writer.WriteString(entity.AuthCode); err != nil {
				return nil, err
			}
		}
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x0102) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	ver := entity.protocolVersion
	if ver == 0 {
		ver = Version2013
	}

	switch ver {
	case Version2019:
		// 长度 + 内容
		ln, err := reader.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("read authCode length: %w", err)
		}
		if ln > 0 {
			if entity.AuthCode, err = reader.ReadString(int(ln)); err != nil {
				return 0, fmt.Errorf("read authCode: %w", err)
			}
		}
		if entity.IMEI, err = reader.ReadString(15); err != nil {
			return 0, fmt.Errorf("read IMEI: %w", err)
		}
		if entity.SoftwareVersion, err = reader.ReadString(20); err != nil {
			return 0, fmt.Errorf("read softwareVersion: %w", err)
		}
	default:
		// 2011/2013: 剩余为 STRING 鉴权码
		s, err := reader.ReadString()
		if err != nil {
			return 0, err
		}
		entity.AuthCode = s
	}
	return len(data) - reader.Len(), nil
}
