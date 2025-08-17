package jtt

import (
	"fmt"
)

// 消息体属性字段的 bit 位.
const (
	bodyLengthBit    uint16 = 0b0000001111111111
	encryptionBit    uint16 = 0b0001110000000000
	fragmentationBit uint16 = 0b0010000000000000
	versionSignBit   uint16 = 0b0100000000000000
	extraBit         uint16 = 0b1000000000000000
)

// VersionType
//
// Note: Property.VersionSign only contains 0/1 to mean (<=2013) vs (>=2019).
// We keep a separate MsgHeader.Version to distinguish 2011 vs 2013 at encode time.
type VersionType int8

const (
	Version2011 VersionType = -1
	Version2013 VersionType = 0
	Version2019 VersionType = 1
)

// MsgHeader 定义消息头
type MsgHeader struct {
	MsgID           MsgID        `json:"msgID"`           // 消息ID
	Property        *Property    `json:"property"`        // 消息体属性
	ProtocolVersion byte         `json:"protocolVersion"` // 协议版本号，默认 0 表示 2011/2013 版本，其他为 2019 后续版本，每次修订递增，初始为1
	PhoneNumber     string       `json:"phoneNumber"`     // 终端手机号
	SerialNumber    uint16       `json:"serialNumber"`    // 消息流水号
	SegmentInfo     *SegmentInfo `json:"segmentInfo"`     // 消息包封装项
	// Version: 由解码默认映射或由用户在编码前设置。
	// 	解码：
	//	  VersionSign==0 => Version2013
	//	  VersionSign==1 => Version2019
	//	编码：
	// 	  当设置为 Version2011/Version2013/Version2019 时，按对应规则写入
	Version VersionType `json:"version"`
}

// Decode 将 []byte 解码成消息头结构体
func (h *MsgHeader) Decode(data []byte) error {
	if len(data) < Message2013HeaderSize {
		return fmt.Errorf("invalid data length: %d", len(data))
	}

	reader := NewReader(data)

	// 读取消息ID
	msgID, err := reader.ReadWord()
	if err != nil {
		return err
	}
	h.MsgID = MsgID(msgID)

	// 读取消息体属性
	attr, err := reader.ReadWord()
	if err != nil {
		return err
	}
	h.Property = &Property{}
	if err := h.Property.Decode(attr); err != nil { // 消息体属性 [2,4) 位
		return err
	}

	// 2013 版本，phoneNumber [4,10) 位 长度 6 位；
	// 2019 版本，phoneNumber [5,15) 位 长度 10 位，第 4 位版本标识。
	switch h.Property.VersionSign {
	case Version2019:
		h.Version = Version2019
		h.ProtocolVersion, err = reader.ReadByte() // 2019 版本，协议版本号 第 4 位
		if err != nil {
			return err
		}

		h.PhoneNumber, err = reader.ReadBcd(10)
	default:
		// 未知比特，按旧版安全处理
		h.Version = Version2013
		h.PhoneNumber, err = reader.ReadBcd(6)
	}
	if err != nil {
		return err
	}

	h.SerialNumber, err = reader.ReadWord()
	if err != nil {
		return err
	}

	if h.Property.IsSegment() { // 消息包封装项
		h.SegmentInfo = &SegmentInfo{}
		h.SegmentInfo.Total, err = reader.ReadWord()
		if err != nil {
			return err
		}
		h.SegmentInfo.Index, err = reader.ReadWord()
		if err != nil {
			return err
		}
	}
	return nil
}

// Encode 将消息头结构体编码成 []byte
func (h *MsgHeader) Encode() (pkt []byte, err error) {
	writer := NewWriter()

	writer.WriteWord(uint16(h.MsgID)) // 消息id

	// 将用户设定的 Version 同步至属性位：2019 写 1，否则写 0。
	if h.Property == nil {
		h.Property = &Property{}
	}
	if h.Version == Version2019 {
		h.Property.VersionSign = Version2019
	} else {
		h.Property.VersionSign = Version2013
	}
	writer.WriteWord(h.Property.Encode()) // 消息体属性

	switch h.Version {
	case Version2019:
		writer.WriteByte(h.ProtocolVersion) // 协议版本号
		writer.WriteBcd(h.PhoneNumber, 10)  // 手机号
	default:
		// 2011 与 2013 在头部仅手机号长度不同于 2019，二者同为 6 位 BCD。
		writer.WriteBcd(h.PhoneNumber, 6)
	}

	writer.WriteWord(h.SerialNumber) // 消息流水号
	if h.SegmentInfo != nil {
		writer.WriteWord(h.SegmentInfo.Total)
		writer.WriteWord(h.SegmentInfo.Index) // 消息包封装项
	}
	return writer.Bytes(), nil
}

func (h *MsgHeader) GetVersion() VersionType {
	return h.Property.Version()
}

func (h *MsgHeader) IsSegment() bool {
	return h.Property.IsSegment()
}

// Property 定义消息体属性.
type Property struct {
	BodyLength   uint16      `json:"bodyLength"`   // 消息体长度
	Encryption   byte        `json:"encryption"`   // 加密类型
	Segmentation byte        `json:"segmentation"` // 分包标识，1：长消息，有分包；0：无分包
	VersionSign  VersionType `json:"versionSign"`  // 版本标识，1：2019 版本；0：2013/2011 版本
	Extra        byte        `json:"extra"`        // 预留一个bit位的保留字段
}

func (p *Property) Decode(bitNum uint16) error {
	p.BodyLength = bitNum & bodyLengthBit // 消息体长度 低 10 位

	// 加密方式 10-12位
	p.Encryption = byte((bitNum & encryptionBit) >> 10)
	p.Segmentation = byte(bitNum & fragmentationBit >> 13)     // 分包 13 位
	p.VersionSign = VersionType(bitNum & versionSignBit >> 14) // 版本标识 14 位
	p.Extra = byte(bitNum & extraBit >> 15)                    // 保留 15位
	return nil
}

func (p *Property) IsSegment() bool {
	return p.Segmentation == 1
}

// Version 返回版本类型，1：2019 版本；0：默认 2013 版本
func (p *Property) Version() VersionType {
	if p.VersionSign == Version2019 {
		return Version2019
	}
	// 其余（含 Version2011/Version2013）均按 2013 处理
	return Version2013
}

func (p *Property) Encode() uint16 {
	var bitNum uint16
	bitNum += p.BodyLength                 // 消息体长度
	bitNum += uint16(p.Encryption) << 10   // 加密方式
	bitNum += uint16(p.Segmentation) << 13 // 分包
	// 版本标识位：仅当 2019 时置 1，2011/2013 置 0
	if p.VersionSign == Version2019 {
		bitNum += 1 << 14
	}
	bitNum += uint16(p.Extra) << 15 // 保留位
	return bitNum
}

// SegmentInfo 定义分包的封装项
type SegmentInfo struct {
	Total uint16 `json:"total"` // 分包后的包总数
	Index uint16 `json:"index"` // 包序号，从1开始
}
