package jtt

// T1078_0x1205 终端上传音视频资源列表
type T1078_0x1205 struct {
	ReplyMsgSerialNo uint16 `json:"replyMsgSerialNo"` // 流水号，对应查询音视频资源列表消息的流水号
	MediaCount       uint32 `json:"mediaCount"`       // 音视频资源总数
}

func (entity *T1078_0x1205) MsgID() MsgID {
	return MsgT1078_0x1205
}

func (entity *T1078_0x1205) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T1078_0x1205) Decode(data []byte) (int, error) {
	return 0, nil
}

type DeviceMedia struct {
	DeviceMediaQuery
	Size uint32 // 文件大小，单位Byte
}

func (m *DeviceMedia) Encode() ([]byte, error) {
	writer := NewWriter()
	_bytes, err := m.DeviceMediaQuery.Encode() // encode deviceMediaQuery
	if err != nil {
		return nil, err
	}
	writer.Write(_bytes)
	writer.WriteDWord(m.Size)
	return writer.Bytes(), nil

}

func (m *DeviceMedia) Decode(data []byte) (int, error) {
	idx, err := m.DeviceMediaQuery.Decode(data)
	if err != nil {
		return 0, err
	}

	reader := NewReader(data[idx:])
	m.Size, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	return 0, nil
}
