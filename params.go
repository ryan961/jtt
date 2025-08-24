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
//
// 终端参数设置各参数项定义及说明：
//   - 0x0001 ParamHeartbeatInterval: 终端心跳发送间隔（秒）
//   - 0x0002 ParamTCPRetryInterval: TCP消息应答超时时间（秒）
//   - 0x0003 ParamTCPRetryTimes: TCP消息重传次数
//   - 0x0004 ParamUDPRetryInterval: UDP消息应答超时时间（秒）
//   - 0x0005 ParamUDPRetryTimes: UDP消息重传次数
//   - 0x0006 ParamSMSRetryInterval: SMS消息应答超时时间（秒）
//   - 0x0007 ParamSMSRetryTimes: SMS消息重传次数
//   - 0x0008 ~ 0x000F 保留
//   - 0x0010 ParamServerAPN: 主服务器APN（CDMA为PPP拨号号码）
//   - 0x0011 ParamServerUser: 主服务器无线拨号用户名
//   - 0x0012 ParamServerPassword: 主服务器无线拨号密码
//   - 0x0013 ParamServerAddress: 主服务器地址（host:port，分号分隔多个）
//   - 0x0014 ParamBackupServerAPN: 备份服务器APN
//   - 0x0015 ParamBackupServerUser: 备份服务器无线拨号用户名
//   - 0x0016 ParamBackupServerPassword: 备份服务器无线拨号密码
//   - 0x0017 ParamBackupServerAddress: 备份服务器地址（host:port，分号分隔多个）
//   - 0x0018 ~ 0x0019 保留
//   - 0x001A ParamICClientDomainName: 道路运输证IC卡认证主服务器IP/域名
//   - 0x001B ParamICClientTCPPort: 道路运输证IC卡认证主服务器TCP端口
//   - 0x001C ParamICClientUDPPort: 道路运输证IC卡认证主服务器UDP端口
//   - 0x001D ParamICClientBackupDomainName: 道路运输证IC卡认证备份服务器IP/域名
//   - 0x001E ~ 0x001F 保留
//   - 0x0020 ParamLocationReportStrategy: 位置汇报策略（0定时；1定距；2定时+定距）
//   - 0x0021 ParamLocationReportScheme: 位置汇报方案（ACC/登录状态）
//   - 0x0022 ParamDriverUnloginReportInterval: 驾驶员未登录汇报时间间隔（秒，>0）
//   - 0x0023 ParamSlaveServerAPN: 从服务器APN（为空则同主服务器）
//   - 0x0024 ParamSlaveServerUser: 从服务器无线拨号用户名（为空则同主服务器）
//   - 0x0025 ParamSlaveServerPassword: 从服务器无线拨号密码（为空则同主服务器）
//   - 0x0026 ParamSlaveServerAddress: 从服务器地址（host:port，分号分隔多个）
//   - 0x0027 ParamSleepReportInterval: 休眠时汇报时间间隔（秒，>0）
//   - 0x0028 ParamEmergencyReportInterval: 紧急报警时汇报时间间隔（秒，>0）
//   - 0x0029 ParamDefaultReportInterval: 缺省时间汇报间隔（秒，>0）
//   - 0x002A ~ 0x002B 保留
//   - 0x002C ParamDefaultDistanceReportInterval: 缺省距离汇报间隔（米，>0）
//   - 0x002D ParamDriverUnloginDistanceReportInterval: 驾驶员未登录时距离汇报间隔（米，>0）
//   - 0x002E ParamSleepDistanceReportInterval: 休眠时距离汇报间隔（米，>0）
//   - 0x002F ParamEmergencyDistanceReportInterval: 紧急报警时距离汇报间隔（米，>0）
//   - 0x0030 ParamTurnAngleReport: 拐点补传角度（<180°）
//   - 0x0031 ParamElectronicFence: 电子围栏半径（非法位移阈值，米）
//   - 0x0032 ParamTimeSection: 违规行驶时段范围（BYTE1..2起始时分；BYTE3..4结束时分）
//   - 0x0033 ~ 0x003F 保留
//   - 0x0040 ParamPhoneNumber: 监控平台电话号码
//   - 0x0041 ParamResetPhoneNumber: 复位电话号码
//   - 0x0042 ParamRestoreFactoryPhoneNumber: 恢复出厂设置电话号码
//   - 0x0043 ParamSMSPhoneNumber: 监控平台SMS电话号码
//   - 0x0044 ParamSMSEventPhoneNumber: 接收终端SMS文本报警号码
//   - 0x0045 ParamAnswerPhoneStrategy: 终端电话接听策略（0自动；1 ACC ON自动）
//   - 0x0046 ParamMaxCallTime: 每次最长通话时间（秒；0不允许；0xFFFFFFFF不限制）
//   - 0x0047 ParamMaxCallTimeInMonth: 当月最长通话时间（秒；0不允许；0xFFFFFFFF不限制）
//   - 0x0048 ParamMonitorPhoneNumber: 监听电话号码
//   - 0x0049 ParamSupervisorPhoneNumber: 监管平台特权短信号码
//   - 0x0050 ParamAlarmMask: 报警屏蔽字（与位置报警标志对应；1屏蔽）
//   - 0x0051 ParamSMSAlarmMask: 报警发送文本SMS开关（与位置报警标志对应；1发送）
//   - 0x0052 ParamPhoneAlarmMask: 报警拍摄开关（与位置报警标志对应；1拍摄）
//   - 0x0053 ParamPhoneAlarmSaveMask: 报警拍摄存储标志（1存储；否则实时上传）
//   - 0x0054 ParamAlarmShootMask: 关键标志（与位置报警标志对应；1关键）
//   - 0x0055 ParamMaxSpeed: 最高速度（km/h）
//   - 0x0056 ParamOverspeedDuration: 超速持续时间（秒）
//   - 0x0057 ParamRunningTimeInterval: 连续驾驶时间门限（秒）
//   - 0x0058 ParamBaseStationReportTimeinterval: 当天累计驾驶时间门限（秒）
//   - 0x0059 ParamStopCarTimeThreshold: 最小休息时间（秒）
//   - 0x005A ParamMaxDrivingTimeOnce: 最长停车时间（秒）
//   - 0x005B ParamOverspeedThreshold: 超速报警预警差值（1/10Km/h）
//   - 0x005C ParamDriverDutyTime: 驾驶员疲劳驾驶预警差值（秒，>0）
//   - 0x005D ParamSpeedThreshold: 碰撞报警参数设置（低8位时间ms；高8位加速度0.1g）
//   - 0x005E ParamSideFlipThreshold: 侧翻报警参数设置（角度°，默认30°）
//   - 0x005F ~ 0x0063 保留
//   - 0x0064 ParamVideoRecordingSettings: 定时拍照控制（通道开关/存储、时间单位、间隔）
//   - 0x0065 ParamVideoRecordingStore: 定距拍照控制（通道开关/存储、距离单位、间隔）
//   - 0x0066 ~ 0x006F 保留
//   - 0x0070 ParamImageQualitySetting: 图像/视频质量（1~10，1最优）
//   - 0x0071 ParamBrightness: 亮度（0~255）
//   - 0x0072 ParamContrast: 对比度（0~127）
//   - 0x0073 ParamSaturation: 饱和度（0~127）
//   - 0x0074 ParamChroma: 色度（0~255）
//   - 0x0075 ~ 0x007F -
//   - 0x0080 ParamDeviceOdometer: 车辆里程表读数（1/10km）
//   - 0x0081 ParamDeviceProvinceID: 车辆所在的省域ID（WORD）
//   - 0x0082 ParamDeviceCityID: 车辆所在的市域ID（WORD）
//   - 0x0083 ParamDevicePlateNumber: 公安交通管理部门颁发的机动车号牌
//   - 0x0084 ParamDevicePlateColor: 车牌颜色（JT/T 697.7-2014；未上牌填0）
//   - 0x0090 ParamGNSS: GNSS 定位模式（bit0 GPS；bit1 北斗；bit2 GLONASS；bit3 Galileo）
//   - 0x0091 ParamGNSSBaudRate: GNSS 波特率（0x00 4800；0x01 9600；0x02 19200；0x03 38400；0x04 57600；0x05 115200）
//   - 0x0092 ParamGNSSOutputFrequency: GNSS 详细定位数据输出频率（0x00 500ms；0x01 1000ms(默认)；0x02 2000ms；0x03 3000ms）
//   - 0x0093 ParamGNSSCollectFrequency: GNSS 详细定位数据采集频率（秒，默认1）
//   - 0x0094 ParamGNSSUploadMode: GNSS 详细定位数据上传方式（0x00本地存储；0x01按时间；0x02按距离；0x0B累计时间；0x0C累计距离；0x0D累计条数）
//   - 0x0095 ParamGNSSUploadSetting: GNSS 详细定位数据上传设置（单位随上传方式变化）
//   - 0x0100 ParamCANBusChannel1CollectInterval: CAN通道1采集时间间隔（ms，0不采集）
//   - 0x0101 ParamCANBusChannel1UploadInterval: CAN通道1上传时间间隔（秒，0不上传）
//   - 0x0102 ParamCANBusChannel2CollectInterval: CAN通道2采集时间间隔（ms，0不采集）
//   - 0x0103 ParamCANBusChannel2UploadInterval: CAN通道2上传时间间隔（秒，0不上传）
//   - 0x0110 ParamCANID: CAN 总线ID 单独采集设置（采集间隔/通道/帧类型/采集方式/ID）
//   - 0x0111 ~ 0x01FF 用于其他CAN总线ID单独采集设置
//   - 0xF000 ~ 0xFFFF 厂商自定义
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
	// ParamServerAddress STRING 主服务器地址,IP或域名,以冒号分割主机和端口,多个服务器使用分号分割
	ParamServerAddress ParamID = 0x0013
	// ParamBackupServerAPN STRING 备份服务器APN
	ParamBackupServerAPN ParamID = 0x0014
	// ParamBackupServerUser STRING 备份服务器无线通信拨号用户名
	ParamBackupServerUser ParamID = 0x0015
	// ParamBackupServerPassword STRING 备份服务器无线通信拨号密码
	ParamBackupServerPassword ParamID = 0x0016
	// ParamBackupServerAddress STRING 备份服务器备份地址,IP地址或域名,以冒号分割主机和端口,多个服务器使用分号分割
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
	// ParamDriverUnloginReportInterval DWORD 驾驶员未登录汇报时间间隔，单位为秒(s)，值大于0。
	ParamDriverUnloginReportInterval ParamID = 0x0022
	// ParamSlaveServerAPN STRING 从服务器 APN。该值为空时，终端应使用主服务器相同配置
	ParamSlaveServerAPN ParamID = 0x0023
	// ParamSlaveServerUser STRING 从服务器无线通信拨号用户名。该值为空时，终端应使用主服务器相同配置
	ParamSlaveServerUser ParamID = 0x0024
	// ParamSlaveServerPassword STRING 从服务器无线通信拨号密码。该值为空时，终端应使用主服务器相同配置
	ParamSlaveServerPassword ParamID = 0x0025
	// ParamSlaveServerAddress STRING 从服务器地址，IP或域名，主机和端口可用冒号分割，多个服务器使用分号分割
	ParamSlaveServerAddress ParamID = 0x0026
	// ParamSleepReportInterval DWORD 休眠时汇报时间间隔，单位为秒(s)，值大于0。
	ParamSleepReportInterval ParamID = 0x0027
	// ParamEmergencyReportInterval DWORD 紧急报警时汇报时间间隔，单位为秒(s)，值大于0。
	ParamEmergencyReportInterval ParamID = 0x0028
	// ParamDefaultReportInterval DWORD 缺省时间汇报间隔，单位为秒(s)，值大于0。
	ParamDefaultReportInterval ParamID = 0x0029
	// ParamDefaultDistanceReportInterval DWORD 缺省距离汇报间隔，单位为米(m)，值大于0。
	ParamDefaultDistanceReportInterval ParamID = 0x002C
	// ParamDriverUnloginDistanceReportInterval DWORD 驾驶员未登录时汇报距离间隔，单位为米(m)，值大于0。
	ParamDriverUnloginDistanceReportInterval ParamID = 0x002D
	// ParamSleepDistanceReportInterval DWORD 休眠时汇报距离间隔，单位为米(m)，值大于0。
	ParamSleepDistanceReportInterval ParamID = 0x002E
	// ParamEmergencyDistanceReportInterval DWORD 紧急报警时汇报距离间隔，单位为米(m)，值大于0。
	ParamEmergencyDistanceReportInterval ParamID = 0x002F

	// ParamTurnAngleReport DWORD 拐点补传角度，值小于180°
	ParamTurnAngleReport ParamID = 0x0030
	// ParamElectronicFence WORD 电子围栏半径（非法位移阈值），单位为米(m)
	ParamElectronicFence ParamID = 0x0031
	// ParamTimeSection BYTE[4] 违规行驶时段范围，精确到分钟（BYTE1..BYTE2：开始时间的小时和分钟；BYTE3..BYTE4:结束时间的小时和分钟）
	//
	// 参数定义见 TimeSection
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
	// ParamAnswerPhoneStrategy DWORD 终端电话接听策略，
	// 0:自动接听；
	// 1:ACC ON时自动接听，ACC OFF时不自动接听
	ParamAnswerPhoneStrategy ParamID = 0x0045
	// ParamMaxCallTime DWORD 每次最长通话时间，单位为秒(s)，0为不允许通话，0xFFFFFFFF为不限制
	ParamMaxCallTime ParamID = 0x0046
	// ParamMaxCallTimeInMonth DWORD 当月最长通话时间，单位为秒(s)，0为不允许通话，0xFFFFFFFF为不限制
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
	// ParamMaxSpeed DWORD 最高速度，单位 km/h
	ParamMaxSpeed ParamID = 0x0055
	// ParamOverspeedDuration DWORD 超速持续时间，单位秒(s)
	ParamOverspeedDuration ParamID = 0x0056
	// ParamRunningTimeInterval DWORD 连续驾驶时间门限，单位为秒(s)
	ParamRunningTimeInterval ParamID = 0x0057
	// ParamBaseStationReportTimeinterval DWORD 当天累计驾驶时间门限，单位为秒(s)
	ParamBaseStationReportTimeinterval ParamID = 0x0058
	// ParamStopCarTimeThreshold DWORD 最小休息时间，单位为秒(s)
	ParamStopCarTimeThreshold ParamID = 0x0059
	// ParamMaxDrivingTimeOnce DWORD 最长停车时间，单位为秒(s)
	ParamMaxDrivingTimeOnce ParamID = 0x005A
	// ParamOverspeedThreshold WORD 超速报警预警差值，单位为1/10Km/h
	ParamOverspeedThreshold ParamID = 0x005B
	// ParamDriverDutyTime WORD 驾驶员疲劳驾驶预警差值，单位为秒(s)，值大于0
	ParamDriverDutyTime ParamID = 0x005C
	// ParamSpeedThreshold WORD 碰撞报警参数设置：
	//
	//	b7-b0: 为碰撞时间,单位为毫秒(ms)
	//	b15-b8: 为碰撞加速度,单位为0.1g，设置范围为0~79，默认为10
	//
	// 参数定义见 SpeedThreshold
	ParamSpeedThreshold ParamID = 0x005D
	// ParamSideFlipThreshold WORD 侧翻报警参数设置：侧翻角度，单位为度(°)，默认为30°
	ParamSideFlipThreshold ParamID = 0x005E

	// ParamTimerShootingControl DWORD 定时拍照控制。位定义：
	//
	//	bit0~bit4: 摄像通道1~5 定时拍照开关标志（0:不允许；1:允许）
	//	bit5~bit7: 保留
	//	bit8~bit12: 摄像通道1~5 定时拍照存储标志（0:存储；1:上传）
	//	bit13~bit15: 保留
	//	bit16: 定时时间单位（0:秒s，小于5s按5s处理；1:分）
	//	bit17~bit31: 定时时间间隔（收到参数设置或重启后执行）
	//
	// 参数定义见 ParamTimerShootingFlags
	ParamTimerShootingControl ParamID = 0x0064
	// ParamDistanceShootingControl DWORD 定距拍照控制。位定义：
	//
	//	bit0~bit4: 摄像通道1~5 定距拍照开关标志（0:不允许；1:允许）
	//	bit5~bit7: 保留
	//	bit8~bit12: 摄像通道1~5 定距拍照存储标志（0:存储；1:上传）
	//	bit13~bit15: 保留
	//	bit16: 定距距离单位（0:米m，小于100m按100m处理；1:千米km）
	//	bit17~bit31: 定距离间隔（收到参数设置或重启后执行）
	//
	// 参数定义见 ParamDistanceShootingFlags
	ParamDistanceShootingControl ParamID = 0x0065

	// ParamImageQualitySetting DWORD 图像/视频质量，设置范围为1~10，1表示最优质量
	ParamImageQualitySetting ParamID = 0x0070
	// ParamBrightness DWORD 亮度，设置范围为 0～255
	ParamBrightness ParamID = 0x0071
	// ParamContrast DWORD 对比度，设置范围为 0～127
	ParamContrast ParamID = 0x0072
	// ParamSaturation DWORD 饱和度，设置范围为 0～127
	ParamSaturation ParamID = 0x0073
	// ParamChroma DWORD 色度，设置范围为 0～255
	ParamChroma ParamID = 0x0074

	// ParamDeviceOdometer DWORD 车辆里程表读数，单位：1/10km
	ParamDeviceOdometer ParamID = 0x0080
	// ParamDeviceProvinceID WORD 车辆所在的省域ID
	ParamDeviceProvinceID ParamID = 0x0081
	// ParamDeviceCityID WORD 车辆所在的市域ID
	ParamDeviceCityID ParamID = 0x0082
	// ParamDevicePlateNumber STRING 公安交通管理部门颁发的机动车号牌
	ParamDevicePlateNumber ParamID = 0x0083
	// ParamDevicePlateColor BYTE 车牌颜色，按照JT/T 697.7-2014中的规定，未上牌车辆填0
	ParamDevicePlateColor ParamID = 0x0084
	// ParamGNSS DWORD GNSS 定位模式，定义如下：
	//
	//	bit0，0:禁用 GPS 定位，1:启用 GPS 定位
	//	bit1，0:禁用 北斗 定位，1:启用 北斗 定位
	//	bit2，0:禁用 GLONASS 定位，1:启用 GLONASS 定位
	//	bit3，0:禁用 Galileo 定位，1:启用 Galileo 定位
	//
	// 参数定义见 GNSS
	ParamGNSS ParamID = 0x0090
	// ParamGNSSBaudRate DWORD GNSS 波特率，定义如下：
	//
	//	0x00：4800
	//	0x01：9600
	//	0x02：19200
	//	0x03：38400
	//	0x04：57600
	//	0x05：115200
	ParamGNSSBaudRate ParamID = 0x0091
	// ParamGNSSOutputFrequency BYTE GNSS 模块详细定位数据输出频率，定义如下：
	//
	//	0x00：500ms
	//	0x01：1000ms (默认值)
	//	0x02：2000ms
	//	0x03：3000ms
	ParamGNSSOutputFrequency ParamID = 0x0092
	// ParamGNSSCollectFrequency DWORD GNSS 模块详细定位数据采集频率，单位为秒(s)，默认为1
	ParamGNSSCollectFrequency ParamID = 0x0093
	// ParamGNSSUploadMode BYTE GNSS 模块详细定位数据上传方式：
	//
	//	0x00：本地存储，不上传（默认值）
	//	0x01：按时间间隔上传
	//	0x02：按距离间隔上传
	//	0xOB：按累计时间上传，达到传输时间后自动停止上传
	//	0x0C：按累计距离上传，达到距离后自动停止上传
	//	0x0D：按累计条数上传，达到上传条数后自动停止上传
	ParamGNSSUploadMode ParamID = 0x0094
	// ParamGNSSUploadSetting DWORD GNSS 模块详细定位数据上传设置：
	//
	//	上传方式为0x01时，单位为秒(s)；
	//	上传方式为0x02时，单位为米(m)；
	//	上传方式为0x0B时，单位为秒(s)；
	//	上传方式为0x0C时，单位为米(m)；
	//	上传方式为0x0D时，单位为条
	ParamGNSSUploadSetting ParamID = 0x0095
	// ParamCANBusChannel1CollectInterval DWORD CAN 总线通道1 采集时间间隔，单位为毫秒(ms)，0表示不采集
	ParamCANBusChannel1CollectInterval ParamID = 0x0100
	// ParamCANBusChannel1UploadInterval WORD CAN 总线通道1 上传时间间隔，单位为秒(s)，0表示不上传
	ParamCANBusChannel1UploadInterval ParamID = 0x0101
	// ParamCANBusChannel2CollectInterval DWORD CAN 总线通道2 采集时间间隔，单位为毫秒(ms)，0表示不采集
	ParamCANBusChannel2CollectInterval ParamID = 0x0102
	// ParamCANBusChannel2UploadInterval WORD CAN 总线通道2 上传时间间隔，单位为秒(s)，0表示不上传
	ParamCANBusChannel2UploadInterval ParamID = 0x0103
	// ParamCANID BYTE[8] CAN 总线 ID 单独采集设置：
	//
	//	bit63-bit32 表示此 ID 采集时间间隔（ms），0 表示不采集
	//	bit31 表示 CAN 通道号，0：CAN1，1：CAN2
	//	bit30 表示帧类型，0：标准帧，1：扩展帧
	//	bit29 表示数据采集方式，0：原始数据，1：采集区间的计算值
	//	bit28-bit0 表示 CAN 总线 ID
	//
	// 参数定义见 CANID
	ParamCANID ParamID = 0x0110
)

// Param 终端参数
type Param struct {
	Id   ParamID
	Data []byte
}

// ID 参数ID
func (param *Param) ID() ParamID {
	return param.Id
}

// expectID 校验当前参数ID是否为期望的ID
func expectID(p *Param, id ParamID, name string) error {
	if p.Id != id {
		return fmt.Errorf("ParamID mismatch: expect %s(%s), got %s", name, id.String(), p.Id.String())
	}
	return nil
}

// SetByte 设为Byte
func (param *Param) SetByte(id ParamID, b byte) *Param {
	param.Id = id
	param.Data = []byte{b}
	return param
}

// SetBytes 设为Bytes
func (param *Param) SetBytes(id ParamID, b []byte) *Param {
	param.Id = id
	buffer := make([]byte, len(b))
	copy(buffer, b)
	param.Data = buffer
	return param
}

// SetUint16 设为Uint16
func (param *Param) SetUint16(id ParamID, n uint16) *Param {
	param.Id = id
	var buffer [2]byte
	binary.BigEndian.PutUint16(buffer[:], n)
	param.Data = buffer[:]
	return param
}

// SetUint32 设为Uint32
func (param *Param) SetUint32(id ParamID, n uint32) *Param {
	param.Id = id
	var buffer [4]byte
	binary.BigEndian.PutUint32(buffer[:], n)
	param.Data = buffer[:]
	return param
}

// SetUint64 设为Uint64
func (param *Param) SetUint64(id ParamID, n uint64) *Param {
	param.Id = id
	var buffer [8]byte
	binary.BigEndian.PutUint64(buffer[:], n)
	param.Data = buffer[:]
	return param
}

// SetString 设为字符串
func (param *Param) SetString(id ParamID, s string) *Param {
	if len(s) == 0 {
		return param.SetBytes(id, nil)
	}
	data, _ := io.ReadAll(transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GB18030.NewEncoder()))
	return param.SetBytes(id, data)
}

// GetByte 读取Byte
func (param *Param) GetByte() (byte, error) {
	if len(param.Data) < 1 {
		return 0, fmt.Errorf("fail to get byte: %w", ErrInvalidBody)
	}
	return param.Data[0], nil
}

// GetBytes 读取Bytes
func (param *Param) GetBytes() ([]byte, error) {
	data := make([]byte, len(param.Data))
	copy(data, param.Data)
	return data, nil
}

// GetUint16 读取Uint16
func (param *Param) GetUint16() (uint16, error) {
	if len(param.Data) < 2 {
		return 0, fmt.Errorf("fail to get uint16: %w", ErrInvalidBody)
	}
	return binary.BigEndian.Uint16(param.Data[:2]), nil
}

// GetUint32 读取Uint32
func (param *Param) GetUint32() (uint32, error) {
	if len(param.Data) < 4 {
		return 0, fmt.Errorf("fail to get uint32: %w", ErrInvalidBody)
	}
	return binary.BigEndian.Uint32(param.Data[:4]), nil
}

// GetUint64 读取Uint64
func (param *Param) GetUint64() (uint64, error) {
	if len(param.Data) < 8 {
		return 0, fmt.Errorf("fail to get uint64: %w", ErrInvalidBody)
	}
	return binary.BigEndian.Uint64(param.Data[:8]), nil
}

// GetString 读取字符串
func (param *Param) GetString() (string, error) {
	if len(param.Data) == 0 {
		return "", nil
	}
	s, err := io.ReadAll(transform.NewReader(bytes.NewReader(param.Data), simplifiedchinese.GB18030.NewDecoder()))
	if err != nil {
		return "", fmt.Errorf("fail to get string: %w", err)
	}
	return string(s), nil
}

// SetHeartbeatInterval 设置参数 0x0001（终端心跳发送间隔，单位秒）。
func (p *Param) SetHeartbeatInterval(v uint32) *Param { return p.SetUint32(ParamHeartbeatInterval, v) }

// GetHeartbeatInterval 读取参数 0x0001（终端心跳发送间隔，单位秒）。
func (p *Param) GetHeartbeatInterval() (uint32, error) {
	if err := expectID(p, ParamHeartbeatInterval, "心跳发送间隔"); err != nil {
		return 0, fmt.Errorf("fail to get heartbeat interval: %w", err)
	}
	return p.GetUint32()
}

// SetTCPRetryInterval 设置参数 0x0002（TCP消息应答超时时间，单位秒）。
func (p *Param) SetTCPRetryInterval(v uint32) *Param { return p.SetUint32(ParamTCPRetryInterval, v) }

// GetTCPRetryInterval 读取参数 0x0002（TCP消息应答超时时间，单位秒）。
func (p *Param) GetTCPRetryInterval() (uint32, error) {
	if err := expectID(p, ParamTCPRetryInterval, "TCP消息应答超时时间"); err != nil {
		return 0, fmt.Errorf("fail to get TCP retry interval: %w", err)
	}
	return p.GetUint32()
}

// SetUDPRetryInterval 设置参数 0x0004（UDP消息应答超时时间，单位秒）。
func (p *Param) SetUDPRetryInterval(v uint32) *Param { return p.SetUint32(ParamUDPRetryInterval, v) }

// GetUDPRetryInterval 读取参数 0x0004（UDP消息应答超时时间，单位秒）。
func (p *Param) GetUDPRetryInterval() (uint32, error) {
	if err := expectID(p, ParamUDPRetryInterval, "UDP消息应答超时时间"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetUDPRetryTimes 设置参数 0x0005（UDP消息重传次数）。
func (p *Param) SetUDPRetryTimes(v uint32) *Param { return p.SetUint32(ParamUDPRetryTimes, v) }

// GetUDPRetryTimes 读取参数 0x0005（UDP消息重传次数）。
func (p *Param) GetUDPRetryTimes() (uint32, error) {
	if err := expectID(p, ParamUDPRetryTimes, "UDP消息重传次数"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetSMSRetryInterval 设置参数 0x0006（SMS消息应答超时时间，单位秒）。
func (p *Param) SetSMSRetryInterval(v uint32) *Param { return p.SetUint32(ParamSMSRetryInterval, v) }

// GetSMSRetryInterval 读取参数 0x0006（SMS消息应答超时时间，单位秒）。
func (p *Param) GetSMSRetryInterval() (uint32, error) {
	if err := expectID(p, ParamSMSRetryInterval, "SMS消息应答超时时间"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetSMSRetryTimes 设置参数 0x0007（SMS消息重传次数）。
func (p *Param) SetSMSRetryTimes(v uint32) *Param { return p.SetUint32(ParamSMSRetryTimes, v) }

// GetSMSRetryTimes 读取参数 0x0007（SMS消息重传次数）。
func (p *Param) GetSMSRetryTimes() (uint32, error) {
	if err := expectID(p, ParamSMSRetryTimes, "SMS消息重传次数"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetServerAPN 设置参数 0x0010（主服务器APN）。
func (p *Param) SetServerAPN(s string) *Param { return p.SetString(ParamServerAPN, s) }

// GetServerAPN 读取参数 0x0010（主服务器APN）。
func (p *Param) GetServerAPN() (string, error) {
	if err := expectID(p, ParamServerAPN, "ServerAPN"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetServerUser 设置参数 0x0011（主服务器无线通信拨号用户名）。
func (p *Param) SetServerUser(s string) *Param { return p.SetString(ParamServerUser, s) }

// GetServerUser 读取参数 0x0011（主服务器无线通信拨号用户名）。
func (p *Param) GetServerUser() (string, error) {
	if err := expectID(p, ParamServerUser, "ServerUser"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetServerPassword 设置参数 0x0012（主服务器无线通信拨号密码）。
func (p *Param) SetServerPassword(s string) *Param { return p.SetString(ParamServerPassword, s) }

// GetServerPassword 读取参数 0x0012（主服务器无线通信拨号密码）。
func (p *Param) GetServerPassword() (string, error) {
	if err := expectID(p, ParamServerPassword, "ServerPassword"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetServerAddress 设置参数 0x0013（主服务器地址,IP或域名,以冒号分割主机和端口,多个服务器使用分号分割）。
func (p *Param) SetServerAddress(s string) *Param { return p.SetString(ParamServerAddress, s) }

// GetServerAddress 读取参数 0x0013（主服务器地址,IP或域名,以冒号分割主机和端口,多个服务器使用分号分割）。
func (p *Param) GetServerAddress() (string, error) {
	if err := expectID(p, ParamServerAddress, "ServerAddress"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetBackupServerAPN 设置参数 0x0014（备份服务器APN）。
func (p *Param) SetBackupServerAPN(s string) *Param { return p.SetString(ParamBackupServerAPN, s) }

// GetBackupServerAPN 读取参数 0x0014（备份服务器APN）。
func (p *Param) GetBackupServerAPN() (string, error) {
	if err := expectID(p, ParamBackupServerAPN, "BackupServerAPN"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetBackupServerUser 设置参数 0x0015（备份服务器无线通信拨号用户名）。
func (p *Param) SetBackupServerUser(s string) *Param { return p.SetString(ParamBackupServerUser, s) }

// GetBackupServerUser 读取参数 0x0015（备份服务器无线通信拨号用户名）。
func (p *Param) GetBackupServerUser() (string, error) {
	if err := expectID(p, ParamBackupServerUser, "BackupServerUser"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetBackupServerPassword 设置参数 0x0016（备份服务器无线通信拨号密码）。
func (p *Param) SetBackupServerPassword(s string) *Param {
	return p.SetString(ParamBackupServerPassword, s)
}

// GetBackupServerPassword 读取参数 0x0016（备份服务器无线通信拨号密码）。
func (p *Param) GetBackupServerPassword() (string, error) {
	if err := expectID(p, ParamBackupServerPassword, "BackupServerPassword"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetBackupServerAddress 设置参数 0x0017（备份服务器地址）。
func (p *Param) SetBackupServerAddress(s string) *Param {
	return p.SetString(ParamBackupServerAddress, s)
}

// GetBackupServerAddress 读取参数 0x0017（备份服务器地址）。
func (p *Param) GetBackupServerAddress() (string, error) {
	if err := expectID(p, ParamBackupServerAddress, "BackupServerAddress"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetICClientDomainName 设置参数 0x001A（IC卡认证主服务器IP或域名）。
func (p *Param) SetICClientDomainName(s string) *Param {
	return p.SetString(ParamICClientDomainName, s)
}

// GetICClientDomainName 读取参数 0x001A（IC卡认证主服务器IP或域名）。
func (p *Param) GetICClientDomainName() (string, error) {
	if err := expectID(p, ParamICClientDomainName, "ICClientDomainName"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetICClientTCPPort 设置参数 0x001B（IC卡认证主服务器TCP端口）。
func (p *Param) SetICClientTCPPort(v uint32) *Param { return p.SetUint32(ParamICClientTCPPort, v) }

// GetICClientTCPPort 读取参数 0x001B（IC卡认证主服务器TCP端口）。
func (p *Param) GetICClientTCPPort() (uint32, error) {
	if err := expectID(p, ParamICClientTCPPort, "IC卡认证主服务器TCP端口"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetICClientUDPPort 设置参数 0x001C（IC卡认证主服务器UDP端口）。
func (p *Param) SetICClientUDPPort(v uint32) *Param { return p.SetUint32(ParamICClientUDPPort, v) }

// GetICClientUDPPort 读取参数 0x001C（IC卡认证主服务器UDP端口）。
func (p *Param) GetICClientUDPPort() (uint32, error) {
	if err := expectID(p, ParamICClientUDPPort, "IC卡认证主服务器UDP端口"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetICClientBackupDomainName 设置参数 0x001D（IC卡认证备份服务器IP或域名，端口同主服务器）。
func (p *Param) SetICClientBackupDomainName(s string) *Param {
	return p.SetString(ParamICClientBackupDomainName, s)
}

// GetICClientBackupDomainName 读取参数 0x001D（IC卡认证备份服务器IP或域名，端口同主服务器）。
func (p *Param) GetICClientBackupDomainName() (string, error) {
	if err := expectID(p, ParamICClientBackupDomainName, "ICClientBackupDomainName"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetLocationReportStrategy 设置参数 0x0020（位置汇报策略，0:定时汇报; 1:定距汇报; 2:定时和定距汇报）。
func (p *Param) SetLocationReportStrategy(v uint32) *Param {
	return p.SetUint32(ParamLocationReportStrategy, v)
}

// GetLocationReportStrategy 读取参数 0x0020（位置汇报策略，0:定时汇报; 1:定距汇报; 2:定时和定距汇报）。
func (p *Param) GetLocationReportStrategy() (uint32, error) {
	if err := expectID(p, ParamLocationReportStrategy, "位置汇报策略"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetLocationReportScheme 设置参数 0x0021（位置汇报方案，0:根据ACC状态; 1:根据登录状态和ACC状态，先判断登录状态，若登录再根据ACC状态）。
func (p *Param) SetLocationReportScheme(v uint32) *Param {
	return p.SetUint32(ParamLocationReportScheme, v)
}

// GetLocationReportScheme 读取参数 0x0021（位置汇报方案，0:根据ACC状态; 1:根据登录状态和ACC状态，先判断登录状态，若登录再根据ACC状态）。
func (p *Param) GetLocationReportScheme() (uint32, error) {
	if err := expectID(p, ParamLocationReportScheme, "位置汇报方案"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetDriverUnloginReportInterval 设置参数 0x0022（驾驶员未登录汇报时间间隔，秒，值大于0）。
func (p *Param) SetDriverUnloginReportInterval(v uint32) *Param {
	return p.SetUint32(ParamDriverUnloginReportInterval, v)
}

// GetDriverUnloginReportInterval 读取参数 0x0022（驾驶员未登录汇报时间间隔，秒，值大于0）。
func (p *Param) GetDriverUnloginReportInterval() (uint32, error) {
	if err := expectID(p, ParamDriverUnloginReportInterval, "驾驶员未登录汇报时间间隔"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetSlaveServerAPN 设置参数 0x0023（从服务器APN。该值为空时，终端应使用主服务器相同配置）。
func (p *Param) SetSlaveServerAPN(s string) *Param { return p.SetString(ParamSlaveServerAPN, s) }

// GetSlaveServerAPN 读取参数 0x0023（从服务器APN。该值为空时，终端应使用主服务器相同配置）。
func (p *Param) GetSlaveServerAPN() (string, error) {
	if err := expectID(p, ParamSlaveServerAPN, "SlaveServerAPN"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetSlaveServerUser 设置参数 0x0024（从服务器无线通信拨号用户名。该值为空时，终端应使用主服务器相同配置）。
func (p *Param) SetSlaveServerUser(s string) *Param { return p.SetString(ParamSlaveServerUser, s) }

// GetSlaveServerUser 读取参数 0x0024（从服务器无线通信拨号用户名。该值为空时，终端应使用主服务器相同配置）。
func (p *Param) GetSlaveServerUser() (string, error) {
	if err := expectID(p, ParamSlaveServerUser, "SlaveServerUser"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetSlaveServerPassword 设置参数 0x0025（从服务器无线通信拨号密码。该值为空时，终端应使用主服务器相同配置）。
func (p *Param) SetSlaveServerPassword(s string) *Param {
	return p.SetString(ParamSlaveServerPassword, s)
}

// GetSlaveServerPassword 读取参数 0x0025（从服务器无线通信拨号密码。该值为空时，终端应使用主服务器相同配置）。
func (p *Param) GetSlaveServerPassword() (string, error) {
	if err := expectID(p, ParamSlaveServerPassword, "SlaveServerPassword"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetSlaveServerAddress 设置参数 0x0026（从服务器地址，IP或域名，主机和端口可用冒号分割，多个服务器使用分号分割）。
func (p *Param) SetSlaveServerAddress(s string) *Param {
	return p.SetString(ParamSlaveServerAddress, s)
}

// GetSlaveServerAddress 读取参数 0x0026（从服务器地址，IP或域名，主机和端口可用冒号分割，多个服务器使用分号分割）。
func (p *Param) GetSlaveServerAddress() (string, error) {
	if err := expectID(p, ParamSlaveServerAddress, "SlaveServerAddress"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetSleepReportInterval 设置参数 0x0027（休眠时汇报时间间隔，单位秒，值大于0）。
func (p *Param) SetSleepReportInterval(v uint32) *Param {
	return p.SetUint32(ParamSleepReportInterval, v)
}

// GetSleepReportInterval 读取参数 0x0027（休眠时汇报时间间隔，单位秒，值大于0）。
func (p *Param) GetSleepReportInterval() (uint32, error) {
	if err := expectID(p, ParamSleepReportInterval, "休眠时汇报时间间隔"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetEmergencyReportInterval 设置参数 0x0028（紧急报警时汇报时间间隔，单位秒，值大于0）。
func (p *Param) SetEmergencyReportInterval(v uint32) *Param {
	return p.SetUint32(ParamEmergencyReportInterval, v)
}

// GetEmergencyReportInterval 读取参数 0x0028（紧急报警时汇报时间间隔，单位秒，值大于0）。
func (p *Param) GetEmergencyReportInterval() (uint32, error) {
	if err := expectID(p, ParamEmergencyReportInterval, "紧急报警时汇报时间间隔"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetDefaultReportInterval 设置参数 0x0029（缺省时间汇报间隔，单位秒，值大于0）。
func (p *Param) SetDefaultReportInterval(v uint32) *Param {
	return p.SetUint32(ParamDefaultReportInterval, v)
}

// GetDefaultReportInterval 读取参数 0x0029（缺省时间汇报间隔，单位秒，值大于0）。
func (p *Param) GetDefaultReportInterval() (uint32, error) {
	if err := expectID(p, ParamDefaultReportInterval, "缺省时间汇报间隔"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetDefaultDistanceReportInterval 设置参数 0x002C（缺省距离汇报间隔，单位米，值大于0）。
func (p *Param) SetDefaultDistanceReportInterval(v uint32) *Param {
	return p.SetUint32(ParamDefaultDistanceReportInterval, v)
}

// GetDefaultDistanceReportInterval 读取参数 0x002C（缺省距离汇报间隔，单位米，值大于0）。
func (p *Param) GetDefaultDistanceReportInterval() (uint32, error) {
	if err := expectID(p, ParamDefaultDistanceReportInterval, "缺省距离汇报间隔"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetDriverUnloginDistanceReportInterval 设置参数 0x002D（驾驶员未登录距离汇报间隔，单位米，值大于0）。
func (p *Param) SetDriverUnloginDistanceReportInterval(v uint32) *Param {
	return p.SetUint32(ParamDriverUnloginDistanceReportInterval, v)
}

// GetDriverUnloginDistanceReportInterval 读取参数 0x002D（驾驶员未登录距离汇报间隔，单位米，值大于0）。
func (p *Param) GetDriverUnloginDistanceReportInterval() (uint32, error) {
	if err := expectID(p, ParamDriverUnloginDistanceReportInterval, "驾驶员未登录距离汇报间隔"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetSleepDistanceReportInterval 设置参数 0x002E（休眠时距离汇报间隔，单位米，值大于0）。
func (p *Param) SetSleepDistanceReportInterval(v uint32) *Param {
	return p.SetUint32(ParamSleepDistanceReportInterval, v)
}

// GetSleepDistanceReportInterval 读取参数 0x002E（休眠时距离汇报间隔，单位米，值大于0）。
func (p *Param) GetSleepDistanceReportInterval() (uint32, error) {
	if err := expectID(p, ParamSleepDistanceReportInterval, "休眠时距离汇报间隔"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetEmergencyDistanceReportInterval 设置参数 0x002F（紧急报警时距离汇报间隔，单位米，值大于0）。
func (p *Param) SetEmergencyDistanceReportInterval(v uint32) *Param {
	return p.SetUint32(ParamEmergencyDistanceReportInterval, v)
}

// GetEmergencyDistanceReportInterval 读取参数 0x002F（紧急报警时距离汇报间隔，单位米，值大于0）。
func (p *Param) GetEmergencyDistanceReportInterval() (uint32, error) {
	if err := expectID(p, ParamEmergencyDistanceReportInterval, "紧急报警时距离汇报间隔"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetTurnAngleReport 设置参数 0x0030（拐点补传角度阈值，值小于180度）。
func (p *Param) SetTurnAngleReport(v uint32) *Param { return p.SetUint32(ParamTurnAngleReport, v) }

// GetTurnAngleReport 读取参数 0x0030（拐点补传角度阈值，值小于180度）。
func (p *Param) GetTurnAngleReport() (uint32, error) {
	if err := expectID(p, ParamTurnAngleReport, "拐点补传角度阈值"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

func (p *Param) SetElectronicFence(v uint16) *Param { return p.SetUint16(ParamElectronicFence, v) }

// GetElectronicFence 读取参数 0x0031（电子围栏半径，单位米）。
func (p *Param) GetElectronicFence() (uint16, error) {
	if err := expectID(p, ParamElectronicFence, "电子围栏半径"); err != nil {
		return 0, err
	}
	return p.GetUint16()
}

// TimeSection 表示参数 0x0032（违规行驶时段范围，精确到分钟）。
//
// 字节定义：
//
//	BYTE1..BYTE2：开始时间的小时和分钟；BYTE3..BYTE4: 结束时间的小时和分钟
//	即：StartHour, StartMinute, EndHour, EndMinute
type TimeSection struct {
	// 开始时间的小时
	StartHour byte
	// 开始时间的分钟
	StartMinute byte
	// 结束时间的小时
	EndHour byte
	// 结束时间的分钟
	EndMinute byte
}

// SetTimeSection 设置参数 0x0032（违规行驶时段范围，精确到分钟）。
func (p *Param) SetTimeSection(v *TimeSection) *Param {
	if v == nil {
		return p.SetBytes(ParamTimeSection, nil)
	}
	b := []byte{byte(v.StartHour), byte(v.StartMinute), byte(v.EndHour), byte(v.EndMinute)}
	return p.SetBytes(ParamTimeSection, b)
}

// GetTimeSection 读取参数 0x0032（违规行驶时段范围，精确到分钟）。
func (p *Param) GetTimeSection() (*TimeSection, error) {
	if err := expectID(p, ParamTimeSection, "违规行驶时段范围"); err != nil {
		return nil, err
	}
	if len(p.Data) < 4 {
		return nil, ErrInvalidBody
	}
	ts := &TimeSection{
		StartHour:   p.Data[0],
		StartMinute: p.Data[1],
		EndHour:     p.Data[2],
		EndMinute:   p.Data[3],
	}
	return ts, nil
}

// SetPhoneNumber 设置参数 0x0040（监控平台电话号码）。
func (p *Param) SetPhoneNumber(s string) *Param { return p.SetString(ParamPhoneNumber, s) }

// GetPhoneNumber 读取参数 0x0040（监控平台电话号码）。
func (p *Param) GetPhoneNumber() (string, error) {
	if err := expectID(p, ParamPhoneNumber, "PhoneNumber"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetResetPhoneNumber 设置参数 0x0041（复位电话号码）。
func (p *Param) SetResetPhoneNumber(s string) *Param { return p.SetString(ParamResetPhoneNumber, s) }

// GetResetPhoneNumber 读取参数 0x0041（复位电话号码）。
func (p *Param) GetResetPhoneNumber() (string, error) {
	if err := expectID(p, ParamResetPhoneNumber, "ResetPhoneNumber"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetRestoreFactoryPhoneNumber 设置参数 0x0042（恢复出厂设置电话号码）。
func (p *Param) SetRestoreFactoryPhoneNumber(s string) *Param {
	return p.SetString(ParamRestoreFactoryPhoneNumber, s)
}

// GetRestoreFactoryPhoneNumber 读取参数 0x0042（恢复出厂设置电话号码）。
func (p *Param) GetRestoreFactoryPhoneNumber() (string, error) {
	if err := expectID(p, ParamRestoreFactoryPhoneNumber, "RestoreFactoryPhoneNumber"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetSMSPhoneNumber 设置参数 0x0043（监控平台SMS电话号码）。
func (p *Param) SetSMSPhoneNumber(s string) *Param { return p.SetString(ParamSMSPhoneNumber, s) }

// GetSMSPhoneNumber 读取参数 0x0043（监控平台SMS电话号码）。
func (p *Param) GetSMSPhoneNumber() (string, error) {
	if err := expectID(p, ParamSMSPhoneNumber, "SMSPhoneNumber"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetSMSEventPhoneNumber 设置参数 0x0044（接收终端SMS文本报警号码）。
func (p *Param) SetSMSEventPhoneNumber(s string) *Param {
	return p.SetString(ParamSMSEventPhoneNumber, s)
}

// GetSMSEventPhoneNumber 读取参数 0x0044（接收终端SMS文本报警号码）。
func (p *Param) GetSMSEventPhoneNumber() (string, error) {
	if err := expectID(p, ParamSMSEventPhoneNumber, "SMSEventPhoneNumber"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetAnswerPhoneStrategy 设置参数 0x0045（终端电话接听策略，0:自动接听；1:ACC ON时自动接听，ACC OFF时不自动接听）。
func (p *Param) SetAnswerPhoneStrategy(v uint32) *Param {
	return p.SetUint32(ParamAnswerPhoneStrategy, v)
}

// GetAnswerPhoneStrategy 读取参数 0x0045（终端电话接听策略，0:自动接听；1:ACC ON时自动接听，ACC OFF时不自动接听）。
func (p *Param) GetAnswerPhoneStrategy() (uint32, error) {
	if err := expectID(p, ParamAnswerPhoneStrategy, "终端电话接听策略"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetMaxCallTime 设置参数 0x0046（每次最长通话时间，秒，0为不允许通话，0xFFFFFFFF为不限制）。
func (p *Param) SetMaxCallTime(v uint32) *Param { return p.SetUint32(ParamMaxCallTime, v) }

// GetMaxCallTime 读取参数 0x0046（每次最长通话时间，秒，0为不允许通话，0xFFFFFFFF为不限制）。
func (p *Param) GetMaxCallTime() (uint32, error) {
	if err := expectID(p, ParamMaxCallTime, "每次最长通话时间"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetMaxCallTimeInMonth 设置参数 0x0047（当月最长通话时间，秒，0为不允许通话，0xFFFFFFFF为不限制）。
func (p *Param) SetMaxCallTimeInMonth(v uint32) *Param {
	return p.SetUint32(ParamMaxCallTimeInMonth, v)
}

// GetMaxCallTimeInMonth 读取参数 0x0047（当月最长通话时间，秒，0为不允许通话，0xFFFFFFFF为不限制）。
func (p *Param) GetMaxCallTimeInMonth() (uint32, error) {
	if err := expectID(p, ParamMaxCallTimeInMonth, "当月最长通话时间"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetMonitorPhoneNumber 设置参数 0x0048（监听电话号码）。
func (p *Param) SetMonitorPhoneNumber(s string) *Param {
	return p.SetString(ParamMonitorPhoneNumber, s)
}

// GetMonitorPhoneNumber 读取参数 0x0048（监听电话号码）。
func (p *Param) GetMonitorPhoneNumber() (string, error) {
	if err := expectID(p, ParamMonitorPhoneNumber, "MonitorPhoneNumber"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetSupervisorPhoneNumber 设置参数 0x0049（监管平台特权短信号码）。
func (p *Param) SetSupervisorPhoneNumber(s string) *Param {
	return p.SetString(ParamSupervisorPhoneNumber, s)
}

// GetSupervisorPhoneNumber 读取参数 0x0049（监管平台特权短信号码）。
func (p *Param) GetSupervisorPhoneNumber() (string, error) {
	if err := expectID(p, ParamSupervisorPhoneNumber, "SupervisorPhoneNumber"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetAlarmMask 设置参数 0x0050（报警屏蔽字，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警被屏蔽）。
func (p *Param) SetAlarmMask(v uint32) *Param { return p.SetUint32(ParamAlarmMask, v) }

// GetAlarmMask 读取参数 0x0050（报警屏蔽字，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警被屏蔽）。
func (p *Param) GetAlarmMask() (uint32, error) {
	if err := expectID(p, ParamAlarmMask, "报警屏蔽字"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetSMSAlarmMask 设置参数 0x0051（报警发送文本SMS开关，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警时发送文本SMS）。
func (p *Param) SetSMSAlarmMask(v uint32) *Param { return p.SetUint32(ParamSMSAlarmMask, v) }

// GetSMSAlarmMask 读取参数 0x0051（报警发送文本SMS开关，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警时发送文本SMS）。
func (p *Param) GetSMSAlarmMask() (uint32, error) {
	if err := expectID(p, ParamSMSAlarmMask, "报警发送文本SMS开关"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetPhoneAlarmMask 设置参数 0x0052（报警拍摄开关，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警时拍摄）。
func (p *Param) SetPhoneAlarmMask(v uint32) *Param { return p.SetUint32(ParamPhoneAlarmMask, v) }

// GetPhoneAlarmMask 读取参数 0x0052（报警拍摄开关，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警时拍摄）。
func (p *Param) GetPhoneAlarmMask() (uint32, error) {
	if err := expectID(p, ParamPhoneAlarmMask, "报警拍摄开关"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetPhoneAlarmSaveMask 设置参数 0x0053（报警拍摄存储标志，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警时拍摄并存储，否则实时上传）。
func (p *Param) SetPhoneAlarmSaveMask(v uint32) *Param {
	return p.SetUint32(ParamPhoneAlarmSaveMask, v)
}

// GetPhoneAlarmSaveMask 读取参数 0x0053（报警拍摄存储标志，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警时拍摄并存储，否则实时上传）。
func (p *Param) GetPhoneAlarmSaveMask() (uint32, error) {
	if err := expectID(p, ParamPhoneAlarmSaveMask, "报警拍摄存储标志"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetAlarmShootMask 设置参数 0x0054（关键标志，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警为关键报警）。
func (p *Param) SetAlarmShootMask(v uint32) *Param { return p.SetUint32(ParamAlarmShootMask, v) }

// GetAlarmShootMask 读取参数 0x0054（关键标志，与位置信息汇报消息中的报警标志相对应，相应位为1则相应报警为关键报警）。
func (p *Param) GetAlarmShootMask() (uint32, error) {
	if err := expectID(p, ParamAlarmShootMask, "关键标志"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetMaxSpeed 设置参数 0x0055（最高速度，单位 km/h）。
func (p *Param) SetMaxSpeed(v uint32) *Param {
	return p.SetUint32(ParamMaxSpeed, v)
}

// GetMaxSpeed 读取参数 0x0055（最高速度，单位 km/h）。
func (p *Param) GetMaxSpeed() (uint32, error) {
	if err := expectID(p, ParamMaxSpeed, "最高速度"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetOverspeedDuration 设置参数 0x0056（超速持续时间，单位秒）。
func (p *Param) SetOverspeedDuration(v uint32) *Param { return p.SetUint32(ParamOverspeedDuration, v) }

// GetOverspeedDuration 读取参数 0x0056（超速持续时间，单位秒）。
func (p *Param) GetOverspeedDuration() (uint32, error) {
	if err := expectID(p, ParamOverspeedDuration, "超速持续时间"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetRunningTimeInterval 设置参数 0x0057（连续驾驶时间门限，单位秒）。
func (p *Param) SetRunningTimeInterval(v uint32) *Param {
	return p.SetUint32(ParamRunningTimeInterval, v)
}

// GetRunningTimeInterval 读取参数 0x0057（连续驾驶时间门限，单位秒）。
func (p *Param) GetRunningTimeInterval() (uint32, error) {
	if err := expectID(p, ParamRunningTimeInterval, "连续驾驶时间门限"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetBaseStationReportTimeinterval 设置参数 0x0058（当天累计驾驶时间门限，单位秒）。
func (p *Param) SetBaseStationReportTimeinterval(v uint32) *Param {
	return p.SetUint32(ParamBaseStationReportTimeinterval, v)
}

// GetBaseStationReportTimeinterval 读取参数 0x0058（当天累计驾驶时间门限，单位秒）。
func (p *Param) GetBaseStationReportTimeinterval() (uint32, error) {
	if err := expectID(p, ParamBaseStationReportTimeinterval, "当天累计驾驶时间门限"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetStopCarTimeThreshold 设置参数 0x0059（最小休息时间，单位秒）。
func (p *Param) SetStopCarTimeThreshold(v uint32) *Param {
	return p.SetUint32(ParamStopCarTimeThreshold, v)
}

// GetStopCarTimeThreshold 读取参数 0x0059（最小休息时间，单位秒）。
func (p *Param) GetStopCarTimeThreshold() (uint32, error) {
	if err := expectID(p, ParamStopCarTimeThreshold, "最小休息时间"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetMaxDrivingTimeOnce 设置参数 0x005A（最长停车时间，单位秒）。
func (p *Param) SetMaxDrivingTimeOnce(v uint32) *Param {
	return p.SetUint32(ParamMaxDrivingTimeOnce, v)
}

// GetMaxDrivingTimeOnce 读取参数 0x005A（最长停车时间，单位秒）。
func (p *Param) GetMaxDrivingTimeOnce() (uint32, error) {
	if err := expectID(p, ParamMaxDrivingTimeOnce, "最长停车时间"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetOverspeedThreshold 设置参数 0x005B（超速报警预警差值，单位为1/10Km/h）。
func (p *Param) SetOverspeedThreshold(v uint32) *Param {
	return p.SetUint32(ParamOverspeedThreshold, v)
}

// GetOverspeedThreshold 读取参数 0x005B（超速报警预警差值，单位为1/10Km/h）。
func (p *Param) GetOverspeedThreshold() (uint32, error) {
	if err := expectID(p, ParamOverspeedThreshold, "超速报警预警差值"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

func (p *Param) SetDriverDutyTime(v uint32) *Param { return p.SetUint32(ParamDriverDutyTime, v) }

// GetDriverDutyTime 读取参数 0x005C（驾驶员疲劳驾驶预警差值，单位秒，值大于0）。
func (p *Param) GetDriverDutyTime() (uint32, error) {
	if err := expectID(p, ParamDriverDutyTime, "驾驶员疲劳驾驶预警差值"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SpeedThreshold 表示参数 0x005D（碰撞报警参数设置）的结构化字段。
//
// 位域定义：
//   - b7-b0:  CollisionTimeMs（碰撞时间，单位毫秒(ms)，0~255）
//   - b15-b8: AccelerationDeciG（碰撞加速度，单位0.1g，范围0~79，默认10）
type SpeedThreshold struct {
	// 碰撞时间，单位毫秒(ms)，0~255
	CollisionTimeMs byte
	// 碰撞加速度，单位0.1g，范围0~79，默认10
	AccelerationDeciG byte
}

// SetSpeedThreshold 设置参数 0x005D（碰撞报警参数设置）。
func (p *Param) SetSpeedThreshold(v *SpeedThreshold) *Param {
	if v == nil {
		return p.SetUint16(ParamSpeedThreshold, 0)
	}
	value := (uint16(v.AccelerationDeciG&0xFF) << 8) | uint16(v.CollisionTimeMs)
	return p.SetUint16(ParamSpeedThreshold, value)
}

// GetSpeedThreshold 读取参数 0x005D（碰撞报警参数设置）。
func (p *Param) GetSpeedThreshold() (*SpeedThreshold, error) {
	if err := expectID(p, ParamSpeedThreshold, "碰撞报警参数设置"); err != nil {
		return nil, err
	}
	u, err := p.GetUint16()
	if err != nil {
		return nil, err
	}
	return &SpeedThreshold{
		CollisionTimeMs:   byte(u & 0x00FF),
		AccelerationDeciG: byte((u >> 8) & 0x00FF),
	}, nil
}

// SetSideFlipThreshold 设置参数 0x005E（侧翻报警参数设置：侧翻角度,单位度,默认为30度）。
func (p *Param) SetSideFlipThreshold(v uint16) *Param { return p.SetUint16(ParamSideFlipThreshold, v) }

// GetSideFlipThreshold 读取参数 0x005E（侧翻报警参数设置：侧翻角度,单位度,默认为30度）。
func (p *Param) GetSideFlipThreshold() (uint16, error) {
	if err := expectID(p, ParamSideFlipThreshold, "侧翻报警参数设置"); err != nil {
		return 0, err
	}
	return p.GetUint16()
}

// TimerShootingFlags 定时拍照控制位定义
//
// 位定义：
//   - bit0~bit4:  摄像通道 1~5 定时拍照开关（0:不允许；1:允许）
//   - bit5~bit7:  保留
//   - bit8~bit12: 摄像通道 1~5 定时拍照存储（0:存储；1:上传）
//   - bit13~bit15: 保留
//   - bit16:      定时时间单位（0:秒s，数值小于5s终端按5s处理；1:分min）
//   - bit17~31:   定时时间间隔（收到参数设置或重启后执行）
type TimerShootingFlags uint32

// SetChannelFlags 设置摄像通道（位）开关/存储（bit0~bit15）标志
//
// 位定义：
//   - bit0~bit4:  摄像通道 1~5 定时拍照开关（0:不允许；1:允许）
//   - bit5~bit7:  保留
//   - bit8~bit12: 摄像通道 1~5 定时拍照存储（0:存储；1:上传）
//   - bit13~bit15: 保留
//
// bit: 0~15 (超过15无效)
// enable: true (设置为1), false (设置为0)
func (f *TimerShootingFlags) SetChannelFlags(bit int, enable bool) {
	if bit < 0 || bit > 15 {
		return
	}
	if enable {
		*f |= 1 << uint(bit)
	} else {
		*f &= ^(1 << uint(bit))
	}
}

// GetChannelFlags 获取摄像通道（位）开关/存储（bit0~bit15）标志
//
// 位定义：
//   - bit0~bit4:  摄像通道 1~5 定时拍照开关（0:不允许；1:允许）
//   - bit5~bit7:  保留
//   - bit8~bit12: 摄像通道 1~5 定时拍照存储（0:存储；1:上传）
//   - bit13~bit15: 保留
//
// bit: 0~15 (超过15无效)
// return: true (bit为1), false (bit为0)
func (f TimerShootingFlags) GetChannelFlags(bit int) bool {
	if bit < 0 || bit > 15 {
		return false
	}
	return (f & (1 << uint(bit))) != 0
}

// SetTimeUnitMinutes 设置定时时间单位（bit16），
// false: 默认秒(s)，数值小于5s终端按5s处理
// true: 分(min)
func (f *TimerShootingFlags) SetTimeUnitMinutes(unit bool) {
	if unit {
		*f |= 1 << 16
	}
}

// GetTimeUnitMinutes 获取定时时间单位（bit16）
//
//	false: 秒(s)，数值小于5s终端按5s处理
//	true: 分(min)
func (f TimerShootingFlags) GetTimeUnitMinutes() bool {
	return (f & (1 << 16)) != 0
}

// SetTimeInterval 设置定时时间间隔（15位，0~32767）
func (f *TimerShootingFlags) SetTimeInterval(interval uint16) {
	if interval > 32767 {
		interval = 32767
	}
	// 仅清除 bit17~bit31，不影响 bit16（时间单位）
	const maskBits17To31 TimerShootingFlags = 0xFFFE0000
	*f = (*f &^ maskBits17To31) | TimerShootingFlags(uint32(interval&0x7FFF)<<17)
}

// GetTimeInterval 获取定时时间间隔（15位，0~32767），单位由 GetTimeUnitMinutes 指定
func (f TimerShootingFlags) GetTimeInterval() uint16 {
	return uint16((uint32(f) >> 17) & 0x7FFF)
}

// SetTimerShootingControl 设置参数 0x0064（定时拍照控制）。
func (p *Param) SetTimerShootingControl(v TimerShootingFlags) *Param {
	return p.SetUint32(ParamTimerShootingControl, uint32(v))
}

// DistanceShootingFlags 定距拍照控制位定义（参数 0x0065）
//
// 位定义：
//   - bit0~bit4:  摄像通道 1~5 定距拍照开关（0:不允许；1:允许）
//   - bit5~bit7:  保留
//   - bit8~bit12: 摄像通道 1~5 定距拍照存储（0:存储；1:上传）
//   - bit13~bit15: 保留
//   - bit16:      定距距离单位（0:米m，小于100m按100m处理；1:千米km）
//   - bit17~31:   定距离间隔（0~32767，收到参数设置或重启后执行）
type DistanceShootingFlags uint32

// SetChannelFlags 设置通道（位）开关/存储（bit0~bit15）标志
// bit: 0~15，enable: true(1)/false(0)
func (f *DistanceShootingFlags) SetChannelFlags(bit int, enable bool) {
	if bit < 0 || bit > 15 {
		return
	}
	if enable {
		*f |= 1 << uint(bit)
	} else {
		*f &= ^(1 << uint(bit))
	}
}

// GetChannelFlags 获取通道（位）开关/存储（bit0~bit15）标志
func (f DistanceShootingFlags) GetChannelFlags(bit int) bool {
	if bit < 0 || bit > 15 {
		return false
	}
	return (f & (1 << uint(bit))) != 0
}

// SetDistanceUnitKm 设置定距单位（bit16）
// false: 米(m)；true: 千米(km)
func (f *DistanceShootingFlags) SetDistanceUnitKm(unit bool) {
	if unit {
		*f |= 1 << 16
	}
}

// GetDistanceUnitKm 获取定距单位（bit16）
func (f DistanceShootingFlags) GetDistanceUnitKm() bool { return (f & (1 << 16)) != 0 }

// SetDistanceInterval 设置定距离间隔（bit17~31，0~32767），单位由 SetDistanceUnitKm 指定
func (f *DistanceShootingFlags) SetDistanceInterval(interval uint16) {
	if interval > 32767 {
		interval = 32767
	}
	const maskBits17To31 DistanceShootingFlags = 0xFFFE0000
	*f = (*f &^ maskBits17To31) | DistanceShootingFlags(uint32(interval&0x7FFF)<<17)
}

// GetDistanceInterval 获取定距离间隔（bit17~31，0~32767）
func (f DistanceShootingFlags) GetDistanceInterval() uint16 {
	return uint16((uint32(f) >> 17) & 0x7FFF)
}

// SetDistanceShootingControl 设置参数 0x0065（定距拍照控制）。
func (p *Param) SetDistanceShootingControl(v DistanceShootingFlags) *Param {
	return p.SetUint32(ParamDistanceShootingControl, uint32(v))
}

// GetDistanceShootingControl 读取参数 0x0065（定距拍照控制）。
func (p *Param) GetDistanceShootingControl() (DistanceShootingFlags, error) {
	if err := expectID(p, ParamDistanceShootingControl, "定距拍照控制"); err != nil {
		return 0, err
	}
	u, err := p.GetUint32()
	if err != nil {
		return 0, err
	}
	return DistanceShootingFlags(u), nil
}

// VideoRecordingStore 表示参数 0x0065（定距拍照控制）的结构化字段。
//
// 位布局：
//   - bit0~bit4:  摄像通道1~5 定距拍照开关（0:不允许；1:允许）
//   - bit5~bit7:  保留
//   - bit8~bit12: 摄像通道1~5 定距拍照存储（0:存储；1:上传）
//   - bit13~bit15: 保留
//   - bit16:      定距距离单位（0:米(m)，小于100m终端按100m处理；1:千米(km)）
//   - bit17~31:   定距离间隔（收到参数设置或重启后执行）
type VideoRecordingStore struct {
	// Enable[i] 对应通道 i+1（1..5）：true:允许，false:不允许
	Enable [5]bool
	// Upload[i] 对应通道 i+1（1..5）：true:上传，false:存储
	Upload [5]bool
	// true:千米(km)，false:米(m)，小于100m终端按100m处理
	DistanceUnitKm bool
	// 定距离间隔（15位，0~32767），单位由 DistanceUnitKm 指定，小于100m终端按100m处理
	DistanceInterval uint16
}

// SetVideoRecordingStore 设置参数 0x0065（定距拍照控制）。
func (p *Param) SetVideoRecordingStore(v *VideoRecordingStore) *Param {
	if v == nil {
		return p.SetUint32(ParamDistanceShootingControl, 0)
	}
	var u uint32
	// 开关 bit0~4
	for i := 0; i < 5; i++ {
		if v.Enable[i] {
			u |= 1 << uint(i)
		}
	}
	// 存储/上传 bit8~12（1=上传）
	for i := 0; i < 5; i++ {
		if v.Upload[i] {
			u |= 1 << uint(8+i)
		}
	}
	// 单位 bit16（1=km）
	if v.DistanceUnitKm {
		u |= 1 << 16
	}
	// 间隔 bit17~31（15位）
	u |= (uint32(v.DistanceInterval) & 0x7FFF) << 17

	return p.SetUint32(ParamDistanceShootingControl, u)
}

// GetVideoRecordingStore 读取参数 0x0065（定距拍照控制）。
func (p *Param) GetVideoRecordingStore() (*VideoRecordingStore, error) {
	if err := expectID(p, ParamDistanceShootingControl, "定距拍照控制"); err != nil {
		return nil, fmt.Errorf("fail to get VideoRecordingStore: %w", err)
	}
	u, err := p.GetUint32()
	if err != nil {
		return nil, fmt.Errorf("fail to get VideoRecordingStore: %w", err)
	}
	var v VideoRecordingStore
	for i := 0; i < 5; i++ {
		v.Enable[i] = ((u >> uint(i)) & 0x1) == 1
		v.Upload[i] = ((u >> uint(8+i)) & 0x1) == 1
	}
	v.DistanceUnitKm = ((u >> 16) & 0x1) == 1
	v.DistanceInterval = uint16((u >> 17) & 0x7FFF)
	return &v, nil
}

// SetImageQualitySetting 设置参数 0x0070（图像/视频质量，设置范围为1~10，1表示最优质量）。
func (p *Param) SetImageQualitySetting(v uint32) *Param {
	return p.SetUint32(ParamImageQualitySetting, v)
}

// GetImageQualitySetting 读取参数 0x0070（图像/视频质量，设置范围为1~10，1表示最优质量）。
func (p *Param) GetImageQualitySetting() (uint32, error) {
	if err := expectID(p, ParamImageQualitySetting, "图像/视频质量"); err != nil {
		return 0, fmt.Errorf("fail to get ImageQualitySetting: %w", err)
	}
	return p.GetUint32()
}

// SetBrightness 设置参数 0x0071（亮度，设置范围为 0～255）。
func (p *Param) SetBrightness(v uint32) *Param { return p.SetUint32(ParamBrightness, v) }

// GetBrightness 读取参数 0x0071（亮度，设置范围为 0～255）。
func (p *Param) GetBrightness() (uint32, error) {
	if err := expectID(p, ParamBrightness, "亮度"); err != nil {
		return 0, fmt.Errorf("fail to get Brightness: %w", err)
	}
	return p.GetUint32()
}

// SetContrast 设置参数 0x0072（对比度，设置范围为 0～127）。
func (p *Param) SetContrast(v uint32) *Param { return p.SetUint32(ParamContrast, v) }

// GetContrast 读取参数 0x0072（对比度，设置范围为 0～127）。
func (p *Param) GetContrast() (uint32, error) {
	if err := expectID(p, ParamContrast, "对比度"); err != nil {
		return 0, fmt.Errorf("fail to get Contrast: %w", err)
	}
	return p.GetUint32()
}

// SetSaturation 设置参数 0x0073（饱和度，设置范围为 0～127）。
func (p *Param) SetSaturation(v uint32) *Param { return p.SetUint32(ParamSaturation, v) }

// GetSaturation 读取参数 0x0073（饱和度，设置范围为 0～127）。
func (p *Param) GetSaturation() (uint32, error) {
	if err := expectID(p, ParamSaturation, "饱和度"); err != nil {
		return 0, fmt.Errorf("fail to get Saturation: %w", err)
	}
	return p.GetUint32()
}

// SetChroma 设置参数 0x0074（色度，设置范围为 0～255）。
func (p *Param) SetChroma(v uint32) *Param { return p.SetUint32(ParamChroma, v) }

// GetChroma 读取参数 0x0074（色度，设置范围为 0～255）。
func (p *Param) GetChroma() (uint32, error) {
	if err := expectID(p, ParamChroma, "色度"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetDeviceOdometer 设置参数 0x0080（车辆里程表读数，单位：1/10km）。
func (p *Param) SetDeviceOdometer(v uint32) *Param { return p.SetUint32(ParamDeviceOdometer, v) }

// GetDeviceOdometer 读取参数 0x0080（车辆里程表读数，单位：1/10km）。
func (p *Param) GetDeviceOdometer() (uint32, error) {
	if err := expectID(p, ParamDeviceOdometer, "车辆里程表读数"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetDeviceProvinceID 设置参数 0x0081（车辆所在的省域ID）。
func (p *Param) SetDeviceProvinceID(v byte) *Param { return p.SetByte(ParamDeviceProvinceID, v) }

// GetDeviceProvinceID 读取参数 0x0081（车辆所在的省域ID）。
func (p *Param) GetDeviceProvinceID() (byte, error) {
	if err := expectID(p, ParamDeviceProvinceID, "车辆所在的省域ID"); err != nil {
		return 0, err
	}
	return p.GetByte()
}

// SetDeviceCityID 设置参数 0x0082（车辆所在的市域ID）。
func (p *Param) SetDeviceCityID(v byte) *Param { return p.SetByte(ParamDeviceCityID, v) }

// GetDeviceCityID 读取参数 0x0082（车辆所在的市域ID）。
func (p *Param) GetDeviceCityID() (byte, error) {
	if err := expectID(p, ParamDeviceCityID, "车辆所在的市域ID"); err != nil {
		return 0, err
	}
	return p.GetByte()
}

// SetDevicePlateNumber 设置参数 0x0083（公安交通管理部门颁发的机动车号牌）。
func (p *Param) SetDevicePlateNumber(s string) *Param { return p.SetString(ParamDevicePlateNumber, s) }

// GetDevicePlateNumber 读取参数 0x0083（公安交通管理部门颁发的机动车号牌）。
func (p *Param) GetDevicePlateNumber() (string, error) {
	if err := expectID(p, ParamDevicePlateNumber, "公安交通管理部门颁发的机动车号牌"); err != nil {
		return "", err
	}
	return p.GetString()
}

// SetDevicePlateColor 设置参数 0x0084（车牌颜色，按照JT/T 697.7-2014中的规定，未上牌车辆填0）。
func (p *Param) SetDevicePlateColor(v byte) *Param { return p.SetByte(ParamDevicePlateColor, v) }

// GetDevicePlateColor 读取参数 0x0084（车牌颜色，按照JT/T 697.7-2014中的规定，未上牌车辆填0）。
func (p *Param) GetDevicePlateColor() (byte, error) {
	if err := expectID(p, ParamDevicePlateColor, "车牌颜色"); err != nil {
		return 0, err
	}
	return p.GetByte()
}

// GNSSAttrs GNSS 模块属性
//
// 位定义：
//
//	bit0：是否支持 GPS 定位（0:不支持；1:支持）
//	bit1：是否支持北斗定位（0:不支持；1:支持）
//	bit2：是否支持 GLONASS 定位（0:不支持；1:支持）
//	bit3：是否支持 Galileo 定位（0:不支持；1:支持）
type GNSSAttrs byte

// GetGPS 获取是否支持 GPS 定位
func (g GNSSAttrs) GetGPS() bool { return GetBitByte(byte(g), 0) }

// SetGPS 设置是否支持 GPS 定位
func (g *GNSSAttrs) SetGPS(v bool) { SetBitByte((*byte)(g), 0, v) }

// GetBeidou 获取是否支持北斗定位
func (g GNSSAttrs) GetBeidou() bool { return GetBitByte(byte(g), 1) }

// SetBeidou 设置是否支持北斗定位
func (g *GNSSAttrs) SetBeidou(v bool) { SetBitByte((*byte)(g), 1, v) }

// GetGLONASS 获取是否支持 GLONASS 定位
func (g GNSSAttrs) GetGLONASS() bool { return GetBitByte(byte(g), 2) }

// SetGLONASS 设置是否支持 GLONASS 定位
func (g *GNSSAttrs) SetGLONASS(v bool) { SetBitByte((*byte)(g), 2, v) }

// GetGalileo 获取是否支持 Galileo 定位
func (g GNSSAttrs) GetGalileo() bool { return GetBitByte(byte(g), 3) }

// SetGalileo 设置是否支持 Galileo 定位
func (g *GNSSAttrs) SetGalileo(v bool) { SetBitByte((*byte)(g), 3, v) }

// SetGNSS 设置参数 0x0090（GNSS 定位模式）。
func (p *Param) SetGNSS(v GNSSAttrs) *Param {
	return p.SetByte(ParamGNSS, byte(v))
}

// GetGNSS 读取参数 0x0090（GNSS 定位模式）。
func (p *Param) GetGNSS() (GNSSAttrs, error) {
	if err := expectID(p, ParamGNSS, "GNSS 定位模式"); err != nil {
		return 0, err
	}
	b, err := p.GetByte()
	if err != nil {
		return 0, err
	}
	return GNSSAttrs(b), nil
}

// SetGNSSBaudRate 设置参数 0x0091（GNSS 波特率）。
//
//	0x00：4800
//	0x01：9600
//	0x02：19200
//	0x03：38400
//	0x04：57600
//	0x05：115200
func (p *Param) SetGNSSBaudRate(v byte) *Param { return p.SetByte(ParamGNSSBaudRate, v) }

// GetGNSSBaudRate 读取参数 0x0091（GNSS 波特率）。
//
//	0x00：4800
//	0x01：9600
//	0x02：19200
//	0x03：38400
//	0x04：57600
//	0x05：115200
func (p *Param) GetGNSSBaudRate() (byte, error) {
	if err := expectID(p, ParamGNSSBaudRate, "GNSS 波特率"); err != nil {
		return 0, err
	}
	return p.GetByte()
}

// SetGNSSOutputFrequency 设置参数 0x0092（GNSS 模块详细定位数据输出频率）。
//
//	0x00：500ms
//	0x01：1000ms (默认值)
//	0x02：2000ms
//	0x03：3000ms
func (p *Param) SetGNSSOutputFrequency(v byte) *Param {
	return p.SetByte(ParamGNSSOutputFrequency, v)
}

// GetGNSSOutputFrequency 读取参数 0x0092（GNSS 模块详细定位数据输出频率）。
//
//	0x00：500ms
//	0x01：1000ms (默认值)
//	0x02：2000ms
//	0x03：3000ms
func (p *Param) GetGNSSOutputFrequency() (byte, error) {
	if err := expectID(p, ParamGNSSOutputFrequency, "GNSS 模块详细定位数据输出频率"); err != nil {
		return 0, err
	}
	return p.GetByte()
}

// SetGNSSCollectFrequency 设置参数 0x0093（GNSS 模块详细定位数据采集频率，单位为秒(s)，默认为1）。
func (p *Param) SetGNSSCollectFrequency(v uint32) *Param {
	return p.SetUint32(ParamGNSSCollectFrequency, v)
}

// GetGNSSCollectFrequency 读取参数 0x0093（GNSS 模块详细定位数据采集频率，单位为秒(s)，默认为1）。
func (p *Param) GetGNSSCollectFrequency() (uint32, error) {
	if err := expectID(p, ParamGNSSCollectFrequency, "GNSS 模块详细定位数据采集频率"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetGNSSUploadMode 设置参数 0x0094（GNSS 模块详细定位数据上传方式）。
//
//	0x00：按时间间隔上传，达到时间间隔后自动停止上传
//	0x01：按累计距离上传，达到距离后自动停止上传
//	0x02：按累计条数上传，达到上传条数后自动停止上传
//	0x0B：按时间间隔上传，达到时间间隔后自动停止上传
//	0x0C：按累计距离上传，达到距离后自动停止上传
//	0x0D：按累计条数上传，达到上传条数后自动停止上传
func (p *Param) SetGNSSUploadMode(v byte) *Param { return p.SetByte(ParamGNSSUploadMode, v) }

// GetGNSSUploadMode 读取参数 0x0094（GNSS 模块详细定位数据上传方式）。
//
//	0x00：按时间间隔上传，达到时间间隔后自动停止上传
//	0x01：按累计距离上传，达到距离后自动停止上传
//	0x02：按累计条数上传，达到上传条数后自动停止上传
//	0x0B：按时间间隔上传，达到时间间隔后自动停止上传
//	0x0C：按累计距离上传，达到距离后自动停止上传
//	0x0D：按累计条数上传，达到上传条数后自动停止上传
func (p *Param) GetGNSSUploadMode() (byte, error) {
	if err := expectID(p, ParamGNSSUploadMode, "GNSS 模块详细定位数据上传方式"); err != nil {
		return 0, err
	}
	return p.GetByte()
}

// SetGNSSUploadSetting 设置参数 0x0095（GNSS 模块详细定位数据上传设置）。
//
//	上传方式为0x01时，单位为秒(s)；
//	上传方式为0x02时，单位为米(m)；
//	上传方式为0x0B时，单位为秒(s)；
//	上传方式为0x0C时，单位为米(m)；
//	上传方式为0x0D时，单位为条
func (p *Param) SetGNSSUploadSetting(v uint32) *Param { return p.SetUint32(ParamGNSSUploadSetting, v) }

// GetGNSSUploadSetting 读取参数 0x0095（GNSS 模块详细定位数据上传设置）。
//
//	上传方式为0x01时，单位为秒(s)；
//	上传方式为0x02时，单位为米(m)；
//	上传方式为0x0B时，单位为秒(s)；
//	上传方式为0x0C时，单位为米(m)；
//	上传方式为0x0D时，单位为条
func (p *Param) GetGNSSUploadSetting() (uint32, error) {
	if err := expectID(p, ParamGNSSUploadSetting, "GNSS 模块详细定位数据上传设置"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetCANBusChannel1CollectInterval 设置参数 0x0100（CAN 总线通道1 采集时间间隔，单位为毫秒(ms)，0表示不采集）。
func (p *Param) SetCANBusChannel1CollectInterval(v uint32) *Param {
	return p.SetUint32(ParamCANBusChannel1CollectInterval, v)
}

// GetCANBusChannel1CollectInterval 读取参数 0x0100（CAN 总线通道1 采集时间间隔，单位为毫秒(ms)，0表示不采集）。
func (p *Param) GetCANBusChannel1CollectInterval() (uint32, error) {
	if err := expectID(p, ParamCANBusChannel1CollectInterval, "CAN 总线通道1 采集时间间隔"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetCANBusChannel1UploadInterval 设置参数 0x0101（CAN 总线通道1 上传时间间隔，单位为秒(s)，0表示不上传）。
func (p *Param) SetCANBusChannel1UploadInterval(v uint16) *Param {
	return p.SetUint16(ParamCANBusChannel1UploadInterval, v)
}

// GetCANBusChannel1UploadInterval 读取参数 0x0101（CAN 总线通道1 上传时间间隔，单位为秒(s)，0表示不上传）。
func (p *Param) GetCANBusChannel1UploadInterval() (uint16, error) {
	if err := expectID(p, ParamCANBusChannel1UploadInterval, "CAN 总线通道1 上传时间间隔"); err != nil {
		return 0, err
	}
	return p.GetUint16()
}

// SetCANBusChannel2CollectInterval 设置参数 0x0102（CAN 总线通道2 采集时间间隔，单位为毫秒(ms)，0表示不采集）。
func (p *Param) SetCANBusChannel2CollectInterval(v uint32) *Param {
	return p.SetUint32(ParamCANBusChannel2CollectInterval, v)
}

// GetCANBusChannel2CollectInterval 读取参数 0x0102（CAN 总线通道2 采集时间间隔，单位为毫秒(ms)，0表示不采集）。
func (p *Param) GetCANBusChannel2CollectInterval() (uint32, error) {
	if err := expectID(p, ParamCANBusChannel2CollectInterval, "CAN 总线通道2 采集时间间隔"); err != nil {
		return 0, err
	}
	return p.GetUint32()
}

// SetCANBusChannel2UploadInterval 设置参数 0x0103（CAN 总线通道2 上传时间间隔，单位为秒(s)，0表示不上传）。
func (p *Param) SetCANBusChannel2UploadInterval(v uint16) *Param {
	return p.SetUint16(ParamCANBusChannel2UploadInterval, v)
}

// GetCANBusChannel2UploadInterval 读取参数 0x0103（CAN 总线通道2 上传时间间隔，单位为秒(s)，0表示不上传）。
func (p *Param) GetCANBusChannel2UploadInterval() (uint16, error) {
	if err := expectID(p, ParamCANBusChannel2UploadInterval, "CAN 总线通道2 上传时间间隔"); err != nil {
		return 0, err
	}
	return p.GetUint16()
}

// CANID 表示 0x0110（CAN 总线ID 单独采集设置）的结构化字段。
//
//	位域定义：
//	- bit63~32: IntervalMs（采集时间间隔，毫秒，0 表示不采集）
//	- bit31: Channel（0: CAN1，1: CAN2）
//	- bit30: ExtendedFrame（0: 标准帧，1: 扩展帧）
//	- bit29: CalcValue（0: 原始数据，1: 采集区间的计算值）
//	- bit28~0: ID（CAN 总线ID）
type CANID struct {
	// 采集时间间隔，毫秒，0 表示不采集
	IntervalMs uint32
	// CAN 通道号，0：CAN1，1：CAN2
	Channel uint8
	// 帧类型，0：标准帧，1：扩展帧
	ExtendedFrame bool
	// 数据采集方式，0：原始数据，1：采集区间的计算值
	CalcValue bool
	// CAN 总线ID
	ID uint32
}

// SetCANID 设置参数 0x0110（CAN 总线ID 单独采集设置）。
//
//	IntervalMs 表示此ID 采集时间间隔（ms），0 表示不采集
//	Channel 表示CAN 通道号，0：CAN1，1：CAN2
//	ExtendedFrame 表示帧类型，0：标准帧，1：扩展帧
//	CalcValue 表示数据采集方式，0：原始数据，1：采集区间的计算值
//	ID 表示CAN 总线ID
func (p *Param) SetCANID(v *CANID) *Param {
	if v == nil {
		// 保持行为一致：若传入nil，则写入0
		return p.SetUint64(ParamCANID, 0)
	}
	var u uint64
	u |= uint64(v.IntervalMs) << 32
	u |= (uint64(v.Channel&0x01) << 31)
	if v.ExtendedFrame {
		u |= 1 << 30
	}
	if v.CalcValue {
		u |= 1 << 29
	}
	u |= uint64(v.ID & 0x1FFFFFFF)
	return p.SetUint64(ParamCANID, u)
}

// GetCANID 读取参数 0x0110（CAN 总线ID 单独采集设置）。
//
//	IntervalMs 表示此ID 采集时间间隔（ms），0 表示不采集
//	Channel 表示CAN 通道号，0：CAN1，1：CAN2
//	ExtendedFrame 表示帧类型，0：标准帧，1：扩展帧
//	CalcValue 表示数据采集方式，0：原始数据，1：采集区间的计算值
//	ID 表示CAN 总线ID
func (p *Param) GetCANID() (*CANID, error) {
	if err := expectID(p, ParamCANID, "CAN 总线ID 单独采集设置"); err != nil {
		return nil, err
	}
	raw, err := p.GetUint64()
	if err != nil {
		return nil, err
	}
	res := &CANID{}
	res.IntervalMs = uint32(raw >> 32)
	res.Channel = uint8((raw >> 31) & 0x1)
	res.ExtendedFrame = ((raw >> 30) & 0x1) == 1
	res.CalcValue = ((raw >> 29) & 0x1) == 1
	res.ID = uint32(raw & 0x1FFFFFFF)
	return res, nil
}
