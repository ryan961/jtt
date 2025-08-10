package jtt

// Type 附加消息类型
type Type byte

const (
	// TypeExtra_0x01 里程
	TypeExtra_0x01 Type = 0x01
	// TypeExtra_0x02 油量
	TypeExtra_0x02 Type = 0x02
	// TypeExtra_0x03 速度
	TypeExtra_0x03 Type = 0x03
	// TypeExtra_0x04 报警确认
	TypeExtra_0x04 Type = 0x04
	// TypeExtra_0x11 超速报警
	TypeExtra_0x11 Type = 0x11
	// TypeExtra_0x12 进出区域报警
	TypeExtra_0x12 Type = 0x12
	// TypeExtra_0x13 路段行驶时间报警
	TypeExtra_0x13 Type = 0x13
	// TypeExtra_0x25 扩展车辆信号状态位
	TypeExtra_0x25 Type = 0x25
	// TypeExtra_0x2a IO状态位
	TypeExtra_0x2a Type = 0x2a
	// TypeExtra_0x2b 模拟量
	TypeExtra_0x2b Type = 0x2b
	// TypeExtra_0x30 无线通信网络信号强度
	TypeExtra_0x30 Type = 0x30
	// TypeExtra_0x31 GNSS定位卫星数
	TypeExtra_0x31 Type = 0x31
)

// T808_0x0200_Extra 附加信息
type T808_0x0200_Extra struct {
	Id     Type
	Length byte
	Data   []byte
}
