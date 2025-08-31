package jtt

import "fmt"

// AreaType 区域类型
//
//	AreaTypeCircle	圆形区域
//	AreaTypeRectangle	矩形区域
//	AreaTypePolygon	多边形区域
//	AreaTypeRoute	线路
type AreaType byte

const (
	AreaTypeCircle    AreaType = 1
	AreaTypeRectangle AreaType = 2
	AreaTypePolygon   AreaType = 3
	AreaTypeRoute     AreaType = 4
)

// T808_0x8608 查询区域或线路数据
type T808_0x8608 struct {
	// 查询类型
	// AreaTypeCircle: 查询圆形区域数据
	// AreaTypeRectangle: 查询矩形区域数据
	// AreaTypePolygon: 查询多边形区域数据
	// AreaTypeRoute: 查询线路数据
	Type AreaType
	// 区域或线路ID列表，空列表表示查询所有
	IDs []uint32
}

func (m *T808_0x8608) MsgID() MsgID { return MsgT808_0x8608 }

func (m *T808_0x8608) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(byte(m.Type))
	w.WriteDWord(uint32(len(m.IDs)))

	for _, id := range m.IDs {
		w.WriteDWord(id)
	}
	return w.Bytes(), nil
}

func (m *T808_0x8608) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	var typeByte byte
	if typeByte, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read type: %w", err)
	}
	m.Type = AreaType(typeByte)

	var count uint32
	if count, err = r.ReadDWord(); err != nil {
		return 0, fmt.Errorf("read ID count: %w", err)
	}

	m.IDs = make([]uint32, 0, int(count))
	for i := 0; i < int(count) && r.Len() >= 4; i++ {
		var id uint32
		if id, err = r.ReadDWord(); err != nil {
			return 0, fmt.Errorf("read ID %d: %w", i, err)
		}
		m.IDs = append(m.IDs, id)
	}

	return len(data) - r.Len(), nil
}
