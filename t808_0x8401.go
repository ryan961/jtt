package jtt

import "fmt"

// T808_0x8401 设置电话本
type T808_0x8401 struct {
	// 设置类型
	// 0: 删除终端上所有存储的联系人
	// 1: 表示更新电话本(删除终端中已有全部联系人并追加消息中的联系人)
	// 2: 表示追加电话本
	// 3: 表示修改电话本(以联系人为索引)
	SettingType byte
	// 联系人列表
	Contacts []T808_0x8401_Contact
}

type T808_0x8401_Contact struct {
	// 标志
	// 1: 呼入
	// 2: 呼出
	// 3: 呼入/呼出
	Flag byte
	// 电话号码
	Phone string
	// 联系人，经 GBK 编码
	ContactName string
}

func (m *T808_0x8401) MsgID() MsgID { return MsgT808_0x8401 }

func (m *T808_0x8401) Encode() ([]byte, error) {
	w := NewWriter()
	// 设置类型
	w.WriteByte(m.SettingType)
	if m.SettingType == 0 { // 删除终端上所有存储的联系人，无后续
		return w.Bytes(), nil
	}

	// 联系人总数
	w.WriteByte(byte(len(m.Contacts)))

	for _, c := range m.Contacts {
		w.WriteByte(c.Flag)
		// 电话号码 GBK: 这里按字节长度写入
		phoneLen, err := GB18030Length(c.Phone)
		if err != nil {
			return nil, fmt.Errorf("get phone length: %w", err)
		}
		if phoneLen > 0xFF {
			return nil, fmt.Errorf("phone too long: %d", phoneLen)
		}
		w.WriteByte(byte(phoneLen))
		if phoneLen > 0 {
			if err := w.WriteString(c.Phone); err != nil {
				return nil, fmt.Errorf("write phone: %w", err)
			}
		}

		nameLen, err := GB18030Length(c.ContactName)
		if err != nil {
			return nil, fmt.Errorf("get contact name length: %w", err)
		}
		if nameLen > 0xFF {
			return nil, fmt.Errorf("contact name too long: %d", nameLen)
		}
		w.WriteByte(byte(nameLen))
		if nameLen > 0 {
			if err := w.WriteString(c.ContactName); err != nil {
				return nil, fmt.Errorf("write contact name: %w", err)
			}
		}
	}

	return w.Bytes(), nil
}

func (m *T808_0x8401) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	m.SettingType = data[0]
	if m.SettingType == 0 { // 删除所有联系人
		return 1, nil
	}

	if len(data) < 2 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}

	r := NewReader(data[1:])
	cnt, err := r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read contact count: %w", err)
	}

	m.Contacts = make([]T808_0x8401_Contact, 0, int(cnt))
	for i := 0; i < int(cnt) && r.Len() > 0; i++ {
		var c T808_0x8401_Contact
		if c.Flag, err = r.ReadByte(); err != nil {
			return 0, fmt.Errorf("read contact flag: %w", err)
		}
		var pl byte
		if pl, err = r.ReadByte(); err != nil {
			return 0, fmt.Errorf("read phone length: %w", err)
		}
		if pl > 0 {
			if c.Phone, err = r.ReadString(int(pl)); err != nil {
				return 0, fmt.Errorf("read phone: %w", err)
			}
		}
		var nl byte
		if nl, err = r.ReadByte(); err != nil {
			return 0, fmt.Errorf("read contact name length: %w", err)
		}
		if nl > 0 {
			if c.ContactName, err = r.ReadString(int(nl)); err != nil {
				return 0, fmt.Errorf("read contact name: %w", err)
			}
		}
		m.Contacts = append(m.Contacts, c)
	}

	return len(data) - r.Len(), nil
}
