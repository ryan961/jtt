package jtt

import (
	"fmt"
)

// T808_0x0805 摄像头立即拍摄命令应答
type T808_0x0805 struct {
	// 应答流水号，对应平台消息的流水号
	ReplyMsgSerialNo uint16
	// 结果
	// 0-成功/确认
	// 1-失败
	// 2-消息有误
	// 3-不支持
	// 4-报警处理确认
	Result uint8
	// 多媒体ID列表，成功时包含拍摄的多媒体数据ID
	MultimediaIDs []uint16
}

func (entity *T808_0x0805) MsgID() MsgID { return MsgT808_0x0805 }

func (entity *T808_0x0805) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入应答流水号
	writer.WriteWord(entity.ReplyMsgSerialNo)

	// 写入结果
	writer.WriteByte(entity.Result)

	// 如果成功且有多媒体ID，写入多媒体ID列表
	if entity.Result == 0 && len(entity.MultimediaIDs) > 0 {
		// 写入多媒体ID数量
		writer.WriteWord(uint16(len(entity.MultimediaIDs)))

		// 写入每个多媒体ID
		for _, id := range entity.MultimediaIDs {
			writer.WriteDWord(uint32(id))
		}
	}

	return writer.Bytes(), nil
}

func (entity *T808_0x0805) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}

	r := NewReader(data)
	var err error

	// 读取应答流水号
	if entity.ReplyMsgSerialNo, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read reply msg serial no: %w", err)
	}

	// 读取结果
	if entity.Result, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read result: %w", err)
	}

	// 如果成功且还有数据，读取多媒体ID列表
	if entity.Result == 0 && r.Len() >= 2 {
		var count uint16
		if count, err = r.ReadWord(); err != nil {
			return 0, fmt.Errorf("read multimedia count: %w", err)
		}

		if r.Len() < int(count)*4 {
			return 0, fmt.Errorf("insufficient data for multimedia IDs")
		}

		entity.MultimediaIDs = make([]uint16, count)
		for i := uint16(0); i < count; i++ {
			var id uint32
			if id, err = r.ReadDWord(); err != nil {
				return 0, fmt.Errorf("read multimedia ID %d: %w", i, err)
			}
			entity.MultimediaIDs[i] = uint16(id)
		}
	}

	return len(data) - r.Len(), nil
}
