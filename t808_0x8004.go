package jtt

import (
	"fmt"
	"time"
)

// T808_0x8004 查询服务器时间应答
type T808_0x8004 struct {
	// 服务器时间，UTC时间(即协调世界时)按照年月日时分秒排列
	// 例如:2017-03-15 17:09:23 表示为 0x170315170923
	Time time.Time
}

func (entity *T808_0x8004) MsgID() MsgID {
	return MsgT808_0x8004
}

func (entity *T808_0x8004) Encode() ([]byte, error) {
	return ToBCDTime(entity.Time), nil
}

func (entity *T808_0x8004) Decode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, fmt.Errorf("invalid data length: %w (need >=6 bytes, got %d)", ErrInvalidBody, len(data))
	}

	var err error
	entity.Time, err = FromBCDTime(data)
	if err != nil {
		return 0, fmt.Errorf("read Time: %w", err)
	}
	return 6, nil
}
