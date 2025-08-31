package jtt

import (
	"fmt"
)

// T808_0x8300 文本信息下发
// Type 2019 版本新增
type T808_0x8300 struct {
	// 标志
	//	2013 版本:
	//	位		标志
	//	0		1:紧急
	//	1		保留
	//	2		1:终端显示器显示
	//	3 		1:终端TTS插读
	//	4		1:广告屏显示
	//	5 		0:中心导航信息,1:CAN故障码信息
	//	6~7		保留
	//	2019 版本:
	//	位		标志
	//	0~1		01:服务;10:紧急;11:通知
	//	2		1:终端显示器显示
	//	3 		1:终端TTS插读
	//	4		---
	//	5 		0:中心导航信息,1:CAN故障码信息
	//	6~7		保留
	Flag byte
	// 1 通知，2服务. 2019 版本新增
	Type byte
	// 文本信息
	Text string

	flag            *T808_0x8300_Flag
	protocolVersion VersionType
}

// T808_0x8300_Flag 文本信息标志位含义
//
//	2013 版本:
//	位		标志
//	0		1:紧急
//	1		保留
//	2		1:终端显示器显示
//	3 		1:终端TTS插读
//	4		1:广告屏显示
//	5 		0:中心导航信息,1:CAN故障码信息
//	6~7		保留
//	2019 版本:
//	位		标志
//	0~1		01:服务;10:紧急;11:通知
//	2		1:终端显示器显示
//	3 		1:终端TTS插读
//	4		---
//	5 		0:中心导航信息,1:CAN故障码信息
//	6~7		保留
type T808_0x8300_Flag struct {
	Urgent          bool // 紧急
	Serve           bool // 服务, 2019 版本新增
	Notify          bool // 通知, 2019 版本新增
	TerminalMonitor bool // 终端显示器显示
	TTS             bool // 终端 TTS 插读
	Ad              bool // 广告屏显示
	ErrCode         bool // true: CAN故障码信息, default: 中心导航信息
}

func (flag *T808_0x8300_Flag) Encode(protocolVersion VersionType) byte {
	var f byte
	switch protocolVersion {
	case 1: // 2019
		if flag.Urgent {
			f |= 0b00000010
		}
		if flag.Serve {
			f |= 0b00000001
		}
		if flag.Notify {
			f |= 0b00000011
		}
		if flag.TerminalMonitor {
			f |= 0b00000100
		}
		if flag.TTS {
			f |= 0b00001000
		}
		if flag.ErrCode {
			f |= 0b00100000
		}

	default: // 2013
		if flag.Urgent {
			f |= 1
		}
		if flag.TerminalMonitor {
			f |= 0b00000100
		}
		if flag.TTS {
			f |= 0b00001000
		}
		if flag.Ad {
			f |= 0b00010000
		}
		if flag.ErrCode {
			f |= 0b00100000
		}
	}
	return f
}

func (entity *T808_0x8300) MsgID() MsgID { return MsgT808_0x8300 }

func (entity *T808_0x8300) SetFlag(flag *T808_0x8300_Flag) {
	entity.flag = flag
}

func (entity *T808_0x8300) SetProtocolVersion(protocolVersion VersionType) {
	entity.protocolVersion = protocolVersion
}

func (entity *T808_0x8300) Encode() ([]byte, error) {
	writer := NewWriter()
	if entity.flag != nil {
		entity.Flag = entity.flag.Encode(entity.protocolVersion)
	}
	writer.WriteByte(entity.Flag)

	if entity.protocolVersion == Version2019 {
		writer.WriteByte(entity.Type)
	}

	if len(entity.Text) > 0 {
		if err := writer.WriteString(entity.Text); err != nil {
			return nil, err
		}
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8300) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}

	r := NewReader(data)
	var err error
	if entity.Flag, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read flag: %w", err)
	}
	if entity.protocolVersion == Version2019 {
		if entity.Type, err = r.ReadByte(); err != nil {
			return 0, fmt.Errorf("read type: %w", err)
		}
	}
	if r.Len() > 0 {
		data, err := r.ReadString()
		if err != nil {
			return 0, err
		}
		entity.Text = data
	}
	return len(data) - r.Len(), nil
}
