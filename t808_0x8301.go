package jtt

import "fmt"

// T808_0x8301 事件设置
type T808_0x8301 struct {
	// 设置类型 BYTE
	//   0: 删除终端现有所有事件（命令后不带后续字节）
	//   1: 更新事件
	//   2: 追加事件
	//   3: 修改事件
	//   4: 删除特定几项事件（之后事件项中无需带事件内容）
	SettingType byte
	// 事件项列表
	Items []T808_0x8301_EventItem
}

// T808_0x8301_EventItem 事件项
type T808_0x8301_EventItem struct {
	ID      byte
	Content string // 当 SettingType==4 时可为空
}

func (m *T808_0x8301) MsgID() MsgID { return MsgT808_0x8301 }

func (m *T808_0x8301) Encode() ([]byte, error) {
	w := NewWriter()
	// 设置类型
	w.WriteByte(m.SettingType)
	if m.SettingType == 0 { // 删除所有事件，不带后续字节
		return w.Bytes(), nil
	}

	// 事件项总数
	w.WriteByte(byte(len(m.Items)))

	// 事件项列表
	for _, it := range m.Items {
		w.WriteByte(it.ID)
		if m.SettingType == 4 { // 删除特定几项事件，不带内容
			continue
		}
		length, err := GB18030Length(it.Content)
		if err != nil {
			return nil, fmt.Errorf("get event content length: %w", err)
		}
		if length > 0xFF {
			return nil, fmt.Errorf("event content too long: %d", length)
		}
		w.WriteByte(byte(length))
		if err := w.WriteString(it.Content); err != nil {
			return nil, fmt.Errorf("write event content: %w", err)
		}
	}
	return w.Bytes(), nil
}

func (m *T808_0x8301) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	m.SettingType = data[0]
	if m.SettingType == 0 { // 删除所有事件，无后续
		return 1, nil
	}

	if len(data) < 2 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}

	// 使用 Reader 解析后续字段
	r := NewReader(data[1:])
	packCount, err := r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read pack count: %w", err)
	}

	m.Items = make([]T808_0x8301_EventItem, 0, int(packCount))
	for i := 0; i < int(packCount) && r.Len() > 0; i++ {
		var id byte
		if id, err = r.ReadByte(); err != nil {
			return 0, fmt.Errorf("read event ID: %w", err)
		}
		item := T808_0x8301_EventItem{ID: id}
		if m.SettingType != 4 { // 携带内容
			var ln byte
			if ln, err = r.ReadByte(); err != nil {
				return 0, fmt.Errorf("read event content length: %w", err)
			}
			if ln > 0 {
				var s string
				if s, err = r.ReadString(int(ln)); err != nil {
					return 0, fmt.Errorf("read event content: %w", err)
				}
				item.Content = s
			}
		}
		m.Items = append(m.Items, item)
	}

	// 总读取长度 = 1(SettingType) + 消耗的 reader 长度
	return len(data) - r.Len(), nil
}
