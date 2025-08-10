package jtt

import "fmt"

// MsgID 消息 ID 枚举.
type MsgID uint16

func (msgID MsgID) String() string {
	return fmt.Sprintf("0x%04X", uint16(msgID))
}

const (
	// MsgT808_0x0001 终端通用应答
	MsgT808_0x0001 MsgID = 0x0001
	// MsgT808_0x0002 终端心跳
	MsgT808_0x0002 MsgID = 0x0002
	// MsgT808_0x0004 查询服务器时间请求
	MsgT808_0x0004 MsgID = 0x0004
	// MsgT808_0x8004 查询服务器时间应答
	MsgT808_0x8004 MsgID = 0x8004
	// MsgT808_0x0104 查询终端参数应答
	MsgT808_0x0104 MsgID = 0x0104
	// MsgT808_0x0200 汇报位置
	MsgT808_0x0200 MsgID = 0x0200
	// MsgT808_0x0704 定位数据批量上传
	MsgT808_0x0704 MsgID = 0x0704

	// MsgT808_0x8001 平台通用应答
	MsgT808_0x8001 MsgID = 0x8001
	// MsgT808_0x8104 查询终端参数
	MsgT808_0x8104 MsgID = 0x8104
	// MsgT808_0x8300 文本信息下发
	MsgT808_0x8300 MsgID = 0x8300
	// MsgT808_0x8600 设置圆形区域
	MsgT808_0x8600 MsgID = 0x8600
	// MsgT808_0x8601 删除圆形区域
	MsgT808_0x8601 MsgID = 0x8601
	// MsgT808_0x8602 设置矩形区域
	MsgT808_0x8602 MsgID = 0x8602
	// MsgT808_0x8603 删除矩形区域
	MsgT808_0x8603 MsgID = 0x8603
	// MsgT808_0x8604 设置多边形区域
	MsgT808_0x8604 MsgID = 0x8604
	// MsgT808_0x8605 删除多边形区域
	MsgT808_0x8605 MsgID = 0x8605
	// MsgT808_0x8606 设置路线
	MsgT808_0x8606 MsgID = 0x8606
	// MsgT808_0x8607 删除路线
	MsgT808_0x8607 MsgID = 0x8607

	// 0x0F00 ~ 0x0FFF 终端上行信息保留（自定义）

	//	驾驶员身份识别指令（吉标）
	// MsgT808_0x8E11 驾驶员身份信息下发（吉标）
	MsgT808_0x8E11 MsgID = 0x8E11
	// MsgT808_0x0E11 驾驶员身份库数据下载应答（吉标）
	MsgT808_0x0E11 MsgID = 0x0E11
	// MsgT808_0x8E12 驾驶员身份库信息查询（吉标）
	MsgT808_0x8E12 MsgID = 0x8E12
	// MsgT808_0x0E12 驾驶员身份库查询应答（吉标）
	MsgT808_0x0E12 MsgID = 0x0E12
	// MsgT808_0x0E10 驾驶员身份识别上报（吉标）
	MsgT808_0x0E10 MsgID = 0x0E10
	// MsgT808_0x8E10 驾驶员身份识别上报应答（吉标）
	MsgT808_0x8E10 MsgID = 0x8E10
)

const (
	// MsgT1078_0x9101 实时音视频传输请求
	MsgT1078_0x9101 MsgID = 0x9101
	// MsgT1078_0x9205 查询资源列表
	MsgT1078_0x9205 MsgID = 0x9205
	// MsgT1078_0x1205 终端上传音视频资源列表
	MsgT1078_0x1205 MsgID = 0x1205
)

// Msg 消息体
type Msg interface {
	MsgID() MsgID
	Encode() ([]byte, error)
	Decode([]byte) (int, error)
}

// Message 消息包
type Message struct {
	Header *MsgHeader
	Body   Msg
}
