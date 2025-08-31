package jtt

import "fmt"

// T808_0x8500 车辆控制（2013版和2019版消息体格式不同）
//
//   - 2013版：Flag 控制标志，bit0: 车门控制，0: 车门解锁，1: 车门加锁
//   - 2019版：Items 控制类型数据
type T808_0x8500 struct {
	// 控制标志（2013版）
	// bit0: 车门控制，0: 车门解锁，1: 车门加锁
	// bit1～bit7: 保留
	Flag byte

	// 控制类型数据（2019版）
	Items []T808_0x8500_Item

	protoVersion VersionType
}

// 控制类型数据格式（2019版）
type T808_0x8500_Item struct {
	// 控制类型ID
	// 0x0001: 车门控制
	// 0x0002 ~ 0x8000：为标准修订预留
	// 0x8001 ~ 0xFFFF：为厂家自定义控制类型
	ID uint16
	// 车门控制参数（当 ID==0x0001 时使用，0:车门锁闭;1:车门开启）
	DoorParam *byte
	// 未知/厂商扩展参数的原始字节（Encode时原样写入；Decode时长度未知，保持为空）
	ParamRaw []byte
}

func (m *T808_0x8500) MsgID() MsgID { return MsgT808_0x8500 }

func (m *T808_0x8500) SetProtoVersion(v VersionType) { m.protoVersion = v }

func (m *T808_0x8500) ProtoVersion() VersionType { return m.protoVersion }

func (m *T808_0x8500) Encode() ([]byte, error) {
	w := NewWriter()
	if m.protoVersion == Version2013 {
		w.WriteByte(m.Flag)
		return w.Bytes(), nil
	}
	// 控制类型数量
	w.WriteWord(uint16(len(m.Items)))

	for i, it := range m.Items {
		// 控制类型ID
		w.WriteWord(it.ID)
		switch it.ID {
		case 0x0001: // 车门：参数1字节
			if it.DoorParam == nil {
				return nil, fmt.Errorf("item %d id=0x0001 missing DoorParam", i)
			}
			w.WriteByte(*it.DoorParam)
		default:
			// 其它类型按原始参数写入（若为空则不写入任何字节）
			if len(it.ParamRaw) > 0 {
				w.Write(it.ParamRaw)
			}
		}
	}
	return w.Bytes(), nil
}

func (m *T808_0x8500) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)

	if m.protoVersion == Version2013 {
		var err error
		m.Flag, err = r.ReadByte()
		return len(data) - r.Len(), fmt.Errorf("read flag: %w", err)
	}

	cnt, err := r.ReadWord()
	if err != nil {
		return len(data) - r.Len(), fmt.Errorf("read item count: %w", err)
	}
	m.Items = make([]T808_0x8500_Item, 0, int(cnt))

	// 由于不同控制类型的参数长度由类型定义，这里仅对已知ID进行解析
	for i := 0; i < int(cnt) && r.Len() > 0; i++ {
		var id uint16
		if id, err = r.ReadWord(); err != nil {
			return 0, fmt.Errorf("read item id: %w", err)
		}
		item := T808_0x8500_Item{ID: id}
		switch id {
		case 0x0001: // 车门：后续1字节
			var p byte
			if p, err = r.ReadByte(); err != nil {
				return 0, fmt.Errorf("read door param: %w", err)
			}
			item.DoorParam = &p
		default:
			// 未知类型无法确定参数长度，保持不读取额外字节。
			// 如需支持，请在此按厂商协议扩展解析逻辑。
		}
		m.Items = append(m.Items, item)
	}

	return len(data) - r.Len(), nil
}
