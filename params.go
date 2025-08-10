package jtt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// ParamID 参数 ID
type ParamID uint32

func (paramID ParamID) String() string {
	return fmt.Sprintf("0x%04X", uint32(paramID))
}

const (
	// ParamHeartbeatInterval DWORD 终端心跳发送间隔，单位为秒(s)
	ParamHeartbeatInterval ParamID = 0x0001
	// ParamTCPRetryInterval DWORD TCP消息应答超时时间，单位为秒(s)
	ParamTCPRetryInterval ParamID = 0x0002
	// ParamTCPRetryTimes DWORD TCP消息重传次数
	ParamTCPRetryTimes ParamID = 0x0003
	// ParamUDPRetryInterval DWORD UDP消息应答超时时间，单位为秒(s)
	ParamUDPRetryInterval ParamID = 0x0004
	// ParamUDPRetryTimes DWORD UDP消息重传次数
	ParamUDPRetryTimes ParamID = 0x0005
	// ParamSMSRetryInterval DWORD SMS消息应答超时时间，单位为秒(s)
	ParamSMSRetryInterval ParamID = 0x0006
	// ParamSMSRetryTimes DWORD SMS消息重传次数
	ParamSMSRetryTimes ParamID = 0x0007

	// ParamServerAPN STRING 主服务器APN，无线通信拨号访问点。若网络制式为CDMA，则该处为PPP拨号号码
	ParamServerAPN ParamID = 0x0010
	// ParamServerUser STRING 主服务器无线通信拨号用户名
	ParamServerUser ParamID = 0x0011
	// ParamServerPassword STRING 主服务器无线通信拨号密码
	ParamServerPassword ParamID = 0x0012
	// ParamServerAddress STRING 主服务器地址,IP地址或域名(主机名，IP或域名，可带端口号，多个服务器使用分号分割)
	ParamServerAddress ParamID = 0x0013
	// ParamBackupServerAPN STRING 备份服务器APN
	ParamBackupServerAPN ParamID = 0x0014
	// ParamBackupServerUser STRING 备份服务器无线通信拨号用户名
	ParamBackupServerUser ParamID = 0x0015
	// ParamBackupServerPassword STRING 备份服务器无线通信拨号密码
	ParamBackupServerPassword ParamID = 0x0016
	// ParamBackupServerAddress STRING 备份服务器地址,IP地址或域名(主机名，IP或域名，可带端口号，多个服务器使用分号分割)
	ParamBackupServerAddress ParamID = 0x0017

	// ParamICClientDomainName STRING 道路运输证IC卡认证主服务器IP地址或域名
	ParamICClientDomainName ParamID = 0x001A
	// ParamICClientTCPPort DWORD 道路运输证IC卡认证主服务器TCP端口
	ParamICClientTCPPort ParamID = 0x001B
	// ParamICClientUDPPort DWORD 道路运输证IC卡认证主服务器UDP端口
	ParamICClientUDPPort ParamID = 0x001C
	// ParamICClientBackupDomainName STRING 道路运输证IC卡认证备份服务器IP地址或域名
	ParamICClientBackupDomainName ParamID = 0x001D

	// ParamLocationReportStrategy DWORD 位置汇报策略，0:定时汇报; 1:定距汇报; 2:定时和定距汇报
	ParamLocationReportStrategy ParamID = 0x0020
	// ParamLocationReportScheme DWORD 位置汇报方案，0:根据ACC状态; 1:根据登录状态和ACC状态，先判断登录状态，若登录再根据ACC状态
	ParamLocationReportScheme ParamID = 0x0021
	// ParamDriverUnloginReportInterval DWORD 驾驶员未登录汇报时间间隔，单位为秒(s)，值为0时禁用
	ParamDriverUnloginReportInterval ParamID = 0x0022
	// ParamSleepReportInterval STRING 休眠时汇报时间间隔，单位为秒(s)，值为0时睡眠时不汇报
	ParamSleepReportInterval ParamID = 0x0023
	// ParamEmergencyReportInterval STRING 紧急报警时汇报时间间隔，单位为秒(s)，值为0时紧急报警时不汇报
	ParamEmergencyReportInterval ParamID = 0x0024
	// ParamDefaultReportInterval STRING 缺省时间汇报间隔，单位为秒(s)，值为0时不汇报
	ParamDefaultReportInterval ParamID = 0x0025
	// ParamDefaultDistanceReportInterval STRING 缺省距离汇报间隔，单位为米(m)，值为0时不汇报
	ParamDefaultDistanceReportInterval ParamID = 0x0026
	// ParamDriverLoginReportInterval DWORD 驾驶员登录汇报时间间隔，单位为秒(s)，值为0时禁用
	ParamDriverLoginReportInterval ParamID = 0x0027
	// ParamAccOnReportInterval DWORD ACC开汇报时间间隔，单位为秒(s)，值为0时禁用
	ParamAccOnReportInterval ParamID = 0x0028
	// ParamAccOffReportInterval DWORD ACC关汇报时间间隔，单位为秒(s)，值为0时禁用
	ParamAccOffReportInterval ParamID = 0x0029

	// ParamAccOffDistanceInterval DWORD ACC关闭后汇报时间间隔，单位为秒(s)，值为0时禁用
	ParamAccOffDistanceInterval ParamID = 0x002C
	// ParamAccOnDistanceInterval DWORD ACC开启后汇报时间间隔，单位为秒(s)，值为0时禁用
	ParamAccOnDistanceInterval ParamID = 0x002D
	// ParamAccOffDistanceReportInterval DWORD 停车时汇报距离间隔，单位为米(m)，值为0时禁用
	ParamAccOffDistanceReportInterval ParamID = 0x002E
	// ParamAccOnDistanceReportInterval DWORD 行驶时汇报距离间隔，单位为米(m)，值为0时禁用
	ParamAccOnDistanceReportInterval ParamID = 0x002F
	// ParamTurnAngleReport DWORD 拐点补传角度，值为0～180°
	ParamTurnAngleReport ParamID = 0x0030
	// ParamElectronicFence DWORD 电子围栏半径(非法位移阈值)，单位为米
	ParamElectronicFence ParamID = 0x0031
	// ParamTimeSection BYTE[4] 违规行驶时段范围，精确到分钟（BYTE1..BYTE2：开始时间的小时和分钟；BYTE3..BYTE4:结束时间的小时和分钟）
	ParamTimeSection ParamID = 0x0032

	// ParamPhoneNumber STRING 监控平台电话号码
	ParamPhoneNumber ParamID = 0x0040
	// ParamResetPhoneNumber STRING 复位电话号码，可采用此电话号码拨打终端电话让终端复位
	ParamResetPhoneNumber ParamID = 0x0041
	// ParamRestoreFactoryPhoneNumber STRING 恢复出厂设置电话号码，可采用此电话号码拨打终端电话让终端恢复出厂设置
	ParamRestoreFactoryPhoneNumber ParamID = 0x0042
	// ParamSMSPhoneNumber STRING 监控平台SMS电话号码
	ParamSMSPhoneNumber ParamID = 0x0043
	// ParamSMSEventPhoneNumber STRING 接收终端SMS文本报警号码
	ParamSMSEventPhoneNumber ParamID = 0x0044
	// ParamAnswerPhoneStrategy DWORD 终端电话接听策略，0:自动接听；1:ACC ON时自动接听，OFF时不自动接听
	ParamAnswerPhoneStrategy ParamID = 0x0045
	// ParamMaxCallTime DWORD 每次最长通话时间，单位为秒(s)，0为不限制
	ParamMaxCallTime ParamID = 0x0046
	// ParamMaxCallTimeInMonth DWORD 当月最长通话时间，单位为秒(s)，0为不限制
	ParamMaxCallTimeInMonth ParamID = 0x0047
	// ParamMonitorPhoneNumber STRING 监听电话号码
	ParamMonitorPhoneNumber ParamID = 0x0048
	// ParamSupervisorPhoneNumber STRING 监管平台特权短信号码
	ParamSupervisorPhoneNumber ParamID = 0x0049

	// ParamAlarmMask DWORD 报警屏蔽字，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警被屏蔽
	ParamAlarmMask ParamID = 0x0050
	// ParamSMSAlarmMask DWORD 报警发送文本SMS开关，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警时发送文本SMS
	ParamSMSAlarmMask ParamID = 0x0051
	// ParamPhoneAlarmMask DWORD 报警拍摄开关，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警时摄像头拍摄
	ParamPhoneAlarmMask ParamID = 0x0052
	// ParamPhoneAlarmSaveMask DWORD 报警拍摄存储标志，与位置信息汇报消息中的报警标志相对应，相应位为1则对相应报警时拍摄的照片进行存储，否则实时上传
	ParamPhoneAlarmSaveMask ParamID = 0x0053
	// ParamAlarmShootMask DWORD 关键标志，与位置信息汇报消息中的报警标志相对应，相应位为1则对相应报警为关键报警
	ParamAlarmShootMask ParamID = 0x0054
	// ParamAlarmPhotographMask DWORD 报警点播放声音开关，与位置信息汇报消息中的报警标志相对应，相应位为1则对相应报警时终端可播放声音
	ParamAlarmPhotographMask ParamID = 0x0055
	// ParamAlarmSMSMask DWORD 报警短信开关，与位置信息汇报消息中的报警标志相对应，相应位为1则对相应报警时发送短信
	ParamAlarmSMSMask ParamID = 0x0056
	// ParamRunningTimeInterval DWORD 运行时间间隔[HH]，值为为秒(s)
	ParamRunningTimeInterval ParamID = 0x0057
	// ParamBaseStationReportTimeinterval DWORD 当天累计驾驶时间门限，单位为秒(s)
	ParamBaseStationReportTimeinterval ParamID = 0x0058
	// ParamStopCarTimeThreshold DWORD 最小休息时间，单位为秒(s)
	ParamStopCarTimeThreshold ParamID = 0x0059
	// ParamMaxDrivingTimeOnce DWORD 最长持续驾驶时间，单位为秒(s)
	ParamMaxDrivingTimeOnce ParamID = 0x005A
	// ParamOverspeedThreshold WORD 超速报警预警差值，单位为1/10Km/h
	ParamOverspeedThreshold ParamID = 0x005B
	// ParamDriverDutyTime WORD 驾驶员疲劳驾驶预警差值，单位为秒(s)，
	ParamDriverDutyTime ParamID = 0x005C
	// ParamSpeedThreshold WORD 碰撞报警参数设置：碰撞时间，单位为毫秒
	ParamSpeedThreshold ParamID = 0x005D

	// ParamVideoRecordingSettings DWORD 定时拍照控制，参数详细定义见表14
	ParamVideoRecordingSettings ParamID = 0x0064
	// ParamVideoRecordingStore DWORD 定距拍照控制，参数详细定义见表15
	ParamVideoRecordingStore ParamID = 0x0065

	// ParamImageQualitySetting DWORD 图像/视频质量，值高表示清晰度高
	ParamImageQualitySetting ParamID = 0x0070
	// ParamBrightness DWORD 亮度，值大表示亮度高
	ParamBrightness ParamID = 0x0071
	// ParamContrast DWORD 对比度，值大表示对比度大
	ParamContrast ParamID = 0x0072
	// ParamSaturation DWORD 饱和度，值大表示饱和度大
	ParamSaturation ParamID = 0x0073
	// ParamChroma DWORD 色度，值大表示色度大
	ParamChroma ParamID = 0x0074

	// ParamDeviceOdometer DWORD 车辆里程表读数，值为/10km
	ParamDeviceOdometer ParamID = 0x0080
	// ParamDeviceProvinceID WORD 车辆所在的省域ID
	ParamDeviceProvinceID ParamID = 0x0081
	// ParamDeviceCityID WORD 车辆所在的市域ID
	ParamDeviceCityID ParamID = 0x0082
	// ParamDevicePlateNumber STRING 公安交通管理部门颁发的机动车号牌
	ParamDevicePlateNumber ParamID = 0x0083
	// ParamDevicePlateColor BYTE 车牌颜色，按照JT/T 697-2014中的规定
	ParamDevicePlateColor ParamID = 0x0084
	// ParamGNSS BYTE GNSS 定位模式，定义如下：bit0，0:禁用GPS定位，1:启用GPS定位；bit1，0:禁用北斗定位，1:启用北斗定位；bit2，0:禁用GLONASS定位，1:启用GLONASS定位；bit3，0:禁用Galileo定位，1:启用Galileo定位
	ParamGNSS ParamID = 0x0090
	// ParamGNSSBaudRate BYTE GNSS 波特率，定义如下：0x00，4800；0x01，9600；0x02，19200；0x03，38400；0x04，57600；0x05，115200
	ParamGNSSBaudRate ParamID = 0x0091
	// ParamGNSSOutputFrequency BYTE GNSS 模块详细定位数据输出频率，定义如下：0x00，500ms；0x01，1000ms (默认值)；0x02，2000ms；0x03，3000ms
	ParamGNSSOutputFrequency ParamID = 0x0092
	// ParamGNSSCollectFrequency DWORD GNSS 模块详细定位数据采集频率，单位为秒(s)，默认为1
	ParamGNSSCollectFrequency ParamID = 0x0093
	// ParamGNSSUploadMode BYTE GNSS 模块详细定位数据上传方式：0x00，本地存储，不上传（默认值）；0x01，按时间间隔上传；0x02，按距离间隔上传；0xO3，按累计时距上传；0x04，按累计距离上传，达到传输条件后自动停止上传；0x05，按累计时距上传，达到传输条件后自动停止上传
	ParamGNSSUploadMode ParamID = 0x0094
	// ParamGNSSUploadSetting DWORD GNSS 模块详细定位数据上传设置：上传方式为0x01时，单位为秒(s)；上传方式为0x02时，单位为米(m)；上传方式为0x03时，单位为秒(s)或米(m)；上传方式为0x04时，单位为米(m)；上传方式为0x05时，单位为秒(s)
	ParamGNSSUploadSetting ParamID = 0x0095
	// ParamCANBusChannel1CollectInterval DWORD CAN 总线通道1 采集时间间隔，单位为毫秒(ms)，0表示不采集
	ParamCANBusChannel1CollectInterval ParamID = 0x0100
	// ParamCANBusChannel1UploadInterval WORD CAN 总线通道1 上传时间间隔，单位为秒(s)，0表示不上传
	ParamCANBusChannel1UploadInterval ParamID = 0x0101
	// ParamCANBusChannel2CollectInterval DWORD CAN 总线通道2 采集时间间隔，单位为毫秒(ms)，0表示不采集
	ParamCANBusChannel2CollectInterval ParamID = 0x0102
	// ParamCANBusChannel2UploadInterval WORD CAN 总线通道2 上传时间间隔，单位为秒(s)，0表示不上传
	ParamCANBusChannel2UploadInterval ParamID = 0x0103
	// ParamCANID BYTE[8] CAN 总线ID 单独采集设置：bit63-bit32 表示此ID 采集时间间隔（ms），0 表示不采集；bit31 表示CAN 通道号，0：CAN1，1：CAN2；bit30 表示帧类型，0：标准帧，1：扩展帧；bit29 表示数据采集方式，0：原始数据，1：采集区间的计算值；bit28-bit0 表示CAN 总线ID
	ParamCANID ParamID = 0x0110
)

// Param 终端参数
type Param struct {
	id         ParamID
	serialized []byte
}

// ID 参数ID
func (param *Param) ID() ParamID {
	return param.id
}

// SetByte 设为Byte
func (param *Param) SetByte(id ParamID, b byte) *Param {
	param.id = id
	param.serialized = []byte{b}
	return param
}

// SetBytes 设为Bytes
func (param *Param) SetBytes(id ParamID, b []byte) *Param {
	param.id = id
	buffer := make([]byte, len(b))
	copy(buffer, b)
	param.serialized = buffer
	return param
}

// SetUint16 设为Uint16
func (param *Param) SetUint16(id ParamID, n uint16) *Param {
	param.id = id
	var buffer [2]byte
	binary.BigEndian.PutUint16(buffer[:], n)
	param.serialized = buffer[:]
	return param
}

// SetUint32 设为Uint32
func (param *Param) SetUint32(id ParamID, n uint32) *Param {
	param.id = id
	var buffer [4]byte
	binary.BigEndian.PutUint32(buffer[:], n)
	param.serialized = buffer[:]
	return param
}

// SetString 设为字符串
func (param *Param) SetString(id ParamID, s string) *Param {
	if len(s) == 0 {
		return param.SetBytes(id, nil)
	}
	data, _ := io.ReadAll(transform.NewReader(
		bytes.NewReader([]byte(s)), simplifiedchinese.GB18030.NewEncoder()))
	return param.SetBytes(id, data)
}

// GetByte 读取Byte
func (param *Param) GetByte() (byte, error) {
	if len(param.serialized) < 1 {
		return 0, ErrInvalidBody
	}
	return param.serialized[0], nil
}

// GetBytes 读取Bytes
func (param *Param) GetBytes() ([]byte, error) {
	data := make([]byte, len(param.serialized))
	copy(data, param.serialized)
	return data, nil
}

// GetUint16 读取Uint16
func (param *Param) GetUint16() (uint16, error) {
	if len(param.serialized) < 2 {
		return 0, ErrInvalidBody
	}
	return binary.BigEndian.Uint16(param.serialized[:2]), nil
}

// GetUint32 读取Uint32
func (param *Param) GetUint32() (uint32, error) {
	if len(param.serialized) < 4 {
		return 0, ErrInvalidBody
	}
	return binary.BigEndian.Uint32(param.serialized[:4]), nil
}

// GetString 读取字符串
func (param *Param) GetString() (string, error) {
	data, err := io.ReadAll(transform.NewReader(
		bytes.NewReader(param.serialized), simplifiedchinese.GB18030.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
