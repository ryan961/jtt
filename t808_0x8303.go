package jtt

import "fmt"

// T808_0x8303 信息点播菜单设置
type T808_0x8303 struct {
	// 设置类型
	// 0：删除全部信息项
	// 1：更新菜单
	// 2：追加菜单
	// 3：修改菜单
	SettingType byte
	// 信息项列表
	Items []T808_0x8303_Item
}

type T808_0x8303_Item struct {
	// 信息类型
	InfoType byte
	// 信息名称，经 GBK 编码
	Name string
}

func (m *T808_0x8303) MsgID() MsgID { return MsgT808_0x8303 }

func (m *T808_0x8303) Encode() ([]byte, error) {
	w := NewWriter()
	// 设置类型
	w.WriteByte(m.SettingType)
	if m.SettingType == 0 { // 删除全部信息项
		return w.Bytes(), nil
	}

	// 信息项总数
	w.WriteByte(byte(len(m.Items)))

	// 列表
	for _, it := range m.Items {
		w.WriteByte(it.InfoType)
		ln, err := GB18030Length(it.Name)
		if err != nil {
			return nil, fmt.Errorf("get name length: %w", err)
		}
		w.WriteWord(uint16(ln))
		if ln > 0 {
			if err := w.WriteString(it.Name); err != nil {
				return nil, fmt.Errorf("write name: %w", err)
			}
		}
	}
	return w.Bytes(), nil
}

func (m *T808_0x8303) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	m.SettingType = data[0]
	if m.SettingType == 0 { // 无后续
		return 1, nil
	}
	if len(data) < 2 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}

	r := NewReader(data[1:])
	cnt, err := r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read item count: %w", err)
	}
	m.Items = make([]T808_0x8303_Item, 0, int(cnt))
	for i := 0; i < int(cnt) && r.Len() > 0; i++ {
		var tp byte
		if tp, err = r.ReadByte(); err != nil {
			return 0, fmt.Errorf("read info type: %w", err)
		}
		var l uint16
		if l, err = r.ReadWord(); err != nil {
			return 0, fmt.Errorf("read name length: %w", err)
		}
		item := T808_0x8303_Item{InfoType: tp}
		if l > 0 {
			if item.Name, err = r.ReadString(int(l)); err != nil {
				return 0, fmt.Errorf("read name: %w", err)
			}
		}
		m.Items = append(m.Items, item)
	}
	return len(data) - r.Len(), nil
}
