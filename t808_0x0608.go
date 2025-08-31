package jtt

import "fmt"

// T808_0x0608 查询区域或线路数据应答（2019版本新增）
type T808_0x0608 struct {
	// 查询类型
	// AreaTypeCircle: 查询圆形区域数据，返回 CircleAreas
	// AreaTypeRectangle: 查询矩形区域数据，返回 RectangleAreas
	// AreaTypePolygon: 查询多边形区域数据，返回 PolygonAreas
	// AreaTypeRoute: 查询线路数据，返回 RouteAreas
	Type AreaType
	// 数据数量
	Count uint32

	// 圆形区域项（Type 为 AreaTypeCircle 时）
	CircleAreas []T808_0x8600
	// 矩形区域项（Type 为 AreaTypeRectangle 时）
	RectangleAreas []T808_0x8602
	// 多边形区域项（Type 为 AreaTypePolygon 时）
	PolygonAreas []T808_0x8604
	// 线路项（Type 为 AreaTypeRoute 时）
	RouteAreas []T808_0x8606
}

func (m *T808_0x0608) MsgID() MsgID { return MsgT808_0x0608 }

func (m *T808_0x0608) Decode(data []byte) (int, error) {
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

	if m.Count, err = r.ReadDWord(); err != nil {
		return 0, fmt.Errorf("read count: %w", err)
	}

	buf, err := r.Read(r.Len())
	if err != nil {
		return 0, fmt.Errorf("read buffer: %w", err)
	}

	switch m.Type {
	case AreaTypeCircle:
		m.CircleAreas = make([]T808_0x8600, 0, int(m.Count))
		for i := 0; i < int(m.Count) && len(buf) >= 4; i++ {
			var area T808_0x8600
			area.SetProtocolVersion(Version2019)
			index, err := area.Decode(buf)
			if err != nil {
				return 0, fmt.Errorf("read circle area %d: %w", i, err)
			}
			m.CircleAreas = append(m.CircleAreas, area)
			buf = buf[index:]
		}
	case AreaTypeRectangle:
		m.RectangleAreas = make([]T808_0x8602, 0, int(m.Count))
		for i := 0; i < int(m.Count) && len(buf) >= 4; i++ {
			var area T808_0x8602
			area.SetProtocolVersion(Version2019)
			index, err := area.Decode(buf)
			if err != nil {
				return 0, fmt.Errorf("read rectangle area %d: %w", i, err)
			}
			m.RectangleAreas = append(m.RectangleAreas, area)
			buf = buf[index:]
		}
	case AreaTypePolygon:
		m.PolygonAreas = make([]T808_0x8604, 0, int(m.Count))
		for i := 0; i < int(m.Count) && len(buf) >= 4; i++ {
			var area T808_0x8604
			area.SetProtocolVersion(Version2019)
			index, err := area.Decode(buf)
			if err != nil {
				return 0, fmt.Errorf("read polygon area %d: %w", i, err)
			}
			m.PolygonAreas = append(m.PolygonAreas, area)
			buf = buf[index:]
		}
	case AreaTypeRoute:
		m.RouteAreas = make([]T808_0x8606, 0, int(m.Count))
		for i := 0; i < int(m.Count) && len(buf) >= 4; i++ {
			var area T808_0x8606
			area.SetProtocolVersion(Version2019)
			index, err := area.Decode(buf)
			if err != nil {
				return 0, fmt.Errorf("read route area %d: %w", i, err)
			}
			m.RouteAreas = append(m.RouteAreas, area)
			buf = buf[index:]
		}
	default:
		return 0, fmt.Errorf("invalid area type: %d", m.Type)
	}

	return len(buf), nil
}

func (m *T808_0x0608) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(byte(m.Type))

	switch m.Type {
	case AreaTypeCircle:
		w.WriteDWord(uint32(len(m.CircleAreas)))
		for _, area := range m.CircleAreas {
			area.SetProtocolVersion(Version2019)
			if _, err := area.Encode(); err != nil {
				return nil, fmt.Errorf("encode area: %w", err)
			}
		}
	case AreaTypeRectangle:
		w.WriteDWord(uint32(len(m.RectangleAreas)))
		for _, area := range m.RectangleAreas {
			area.SetProtocolVersion(Version2019)
			if _, err := area.Encode(); err != nil {
				return nil, fmt.Errorf("encode area: %w", err)
			}
		}
	case AreaTypePolygon:
		w.WriteDWord(uint32(len(m.PolygonAreas)))
		for _, area := range m.PolygonAreas {
			area.SetProtocolVersion(Version2019)
			if _, err := area.Encode(); err != nil {
				return nil, fmt.Errorf("encode area: %w", err)
			}
		}
	case AreaTypeRoute:
		w.WriteDWord(uint32(len(m.RouteAreas)))
		for _, area := range m.RouteAreas {
			area.SetProtocolVersion(Version2019)
			if _, err := area.Encode(); err != nil {
				return nil, fmt.Errorf("encode area: %w", err)
			}
		}
	default:
		return nil, fmt.Errorf("invalid area type: %d", m.Type)
	}
	return w.Bytes(), nil
}
