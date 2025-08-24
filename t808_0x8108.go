package jtt

import (
	"fmt"
	"math"
)

// UpgradeType 升级类型（BYTE）
//
//	UpgradeTypeTerminal: 终端
//	UpgradeTypeICCard: 道路运输 IC 卡读卡器
//	UpgradeTypeBDModule: 北斗卫星定位模块
type UpgradeType byte

const (
	UpgradeTypeTerminal UpgradeType = 0
	UpgradeTypeICCard   UpgradeType = 12
	UpgradeTypeBDModule UpgradeType = 52
)

// T808_0x8108 下发终端升级包
type T808_0x8108 struct {
	// 升级类型
	//	UpgradeTypeTerminal: 终端
	//	UpgradeTypeICCard: 道路运输 IC 卡读卡器
	//	UpgradeTypeBDModule: 北斗卫星定位模块
	Type UpgradeType
	// 制造商ID，固定5字节
	ManufacturerID string
	// 版本号
	Version string
	// 升级数据包
	Data []byte
}

func (m *T808_0x8108) MsgID() MsgID { return MsgT808_0x8108 }

func (m *T808_0x8108) Encode() ([]byte, error) {
	w := NewWriter()
	// 升级类型
	w.WriteByte(byte(m.Type))
	// 制造商ID BYTE[5]
	if n, err := GB18030Length(m.ManufacturerID); err != nil {
		return nil, fmt.Errorf("get GB18030 length for ManufacturerID: %w", err)
	} else if n != 5 {
		return nil, fmt.Errorf("invalid ManufacturerID: %w (need 5 bytes, got %d)", ErrInvalidBody, n)
	}
	if err := w.WriteString(m.ManufacturerID, 5); err != nil {
		return nil, fmt.Errorf("write ManufacturerID: %w", err)
	}
	// 版本号长度 + 版本号
	if n, err := GB18030Length(m.Version); err != nil {
		return nil, fmt.Errorf("get GB18030 length for Version: %w", err)
	} else {
		if n > int(math.MaxUint8) {
			return nil, fmt.Errorf("invalid Version: %w (need <=%d bytes, got %d)", ErrInvalidBody, math.MaxUint8, n)
		}
		w.WriteByte(byte(n))
		if err := w.WriteString(m.Version, n); err != nil {
			return nil, fmt.Errorf("write Version: %w", err)
		}
	}
	// 升级数据包长度 + 数据
	w.WriteDWord(uint32(len(m.Data)))
	if len(m.Data) > 0 {
		w.Write(m.Data)
	}
	return w.Bytes(), nil
}

func (m *T808_0x8108) Decode(data []byte) (int, error) {
	r := NewReader(data)
	var err error
	// 升级类型
	b, err := r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read Type: %w", err)
	}
	m.Type = UpgradeType(b)
	// 制造商ID BYTE[5]
	var buf []byte
	if buf, err = r.Read(5); err != nil {
		return 0, fmt.Errorf("read ManufacturerID: %w", err)
	}
	m.ManufacturerID = BytesToString(buf)
	// 版本号长度 + 版本号
	var n byte
	if n, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read Version length: %w", err)
	}
	if n > 0 {
		if m.Version, err = r.ReadString(int(n)); err != nil {
			return 0, fmt.Errorf("read Version: %w", err)
		}
	}
	// 升级数据包长度 + 数据
	ln, err := r.ReadDWord()
	if err != nil {
		return 0, fmt.Errorf("read Data length: %w", err)
	}
	if ln > uint32(r.Len()) {
		return 0, fmt.Errorf("invalid Data length: %w (remain %d, want %d)", ErrInvalidBody, r.Len(), ln)
	}
	if ln > 0 {
		if m.Data, err = r.Read(int(ln)); err != nil {
			return 0, fmt.Errorf("read Data: %w", err)
		}
	} else {
		m.Data = nil
	}
	return len(data) - r.Len(), nil
}
