package jtt

import (
	"fmt"
	"time"
)

// T808_0x0702 驾驶员身份信息采集上报
type T808_0x0702 struct {
	// 状态
	// 0x01: 从业资格证IC卡插入（驾驶员上班）
	// 0x02: 从业资格证IC卡拔出（驾驶员下班）
	Status byte
	// 插卡/拔卡时间，YY-MM-DD-hh-mm-ss
	// 以下字段在状态为0x01时才有效微调错误
	Time time.Time
	// IC卡读取结果
	// 0x00: IC卡读取成功
	// 0x01: 读卡失败，原因为卡片密钥认证未通过
	// 0x02: 读卡失败，原因为卡片已被锁定
	// 0x03: 读卡失败，原因为卡片被拔出
	// 0x04: 读卡失败，原因为数据校验错误
	// 以下字段在IC卡读取结果为0x00时才有效
	ICCardReadResult byte
	// 驾驶员姓名
	DriverName string
	// 从业资格证编码，长度20位，不足补0x00
	QualificationCode string
	// 发证机构名称，从业资格证发证机构名称
	IssuingAuthority string
	// 证件有效期，YYYYMMDD
	CertificateValidity time.Time
	// 驾驶员身份证号，长度20位，不足补0x00（2019版本新增）
	DriverIDCardNumber string

	protocolVersion VersionType
}

func (m *T808_0x0702) MsgID() MsgID { return MsgT808_0x0702 }

func (m *T808_0x0702) SetProtocolVersion(protocolVersion VersionType) {
	m.protocolVersion = protocolVersion
}

func (m *T808_0x0702) ProtocolVersion() VersionType { return m.protocolVersion }

func (m *T808_0x0702) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(m.Status)
	w.WriteBcdTime(m.Time)

	// 只有状态为0x01时才写入后续字段
	if m.Status == 0x01 {
		w.WriteByte(m.ICCardReadResult)

		// 只有IC卡读取成功时才写入后续字段
		if m.ICCardReadResult == 0x00 {
			// 驾驶员姓名长度
			driverNameLen, err := GB18030Length(m.DriverName)
			if err != nil {
				return nil, fmt.Errorf("get driver name length: %w", err)
			}
			w.WriteByte(byte(driverNameLen))

			// 驾驶员姓名
			if err := w.WriteString(m.DriverName); err != nil {
				return nil, fmt.Errorf("write driver name: %w", err)
			}

			// 从业资格证编码，固定20字节
			if err := w.WriteString(m.QualificationCode, 20); err != nil {
				return nil, fmt.Errorf("write qualification code: %w", err)
			}

			// 发证机构名称长度
			issuingAuthorityLen, err := GB18030Length(m.IssuingAuthority)
			if err != nil {
				return nil, fmt.Errorf("get issuing authority length: %w", err)
			}
			w.WriteByte(byte(issuingAuthorityLen))

			// 发证机构名称
			if err := w.WriteString(m.IssuingAuthority); err != nil {
				return nil, fmt.Errorf("write issuing authority: %w", err)
			}

			// 证件有效期，BCD[4] YYYYMMDD
			w.WriteBcd(m.CertificateValidity.Format("20060102"), 4)

			// 驾驶员身份证号，长度20位，不足补0x00（2019版本新增）
			if m.protocolVersion >= Version2019 {
				if err := w.WriteString(m.DriverIDCardNumber, 20); err != nil {
					return nil, fmt.Errorf("write driver id card number: %w", err)
				}
			}
		}
	}

	return w.Bytes(), nil
}

func (m *T808_0x0702) Decode(data []byte) (int, error) {
	if len(data) < 7 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	if m.Status, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read status: %w", err)
	}

	if m.Time, err = r.ReadBcdTime(); err != nil {
		return 0, fmt.Errorf("read time: %w", err)
	}

	// 只有状态为0x01时才读取后续字段
	if m.Status == 0x01 {
		if m.ICCardReadResult, err = r.ReadByte(); err != nil {
			return 0, fmt.Errorf("read ic card read result: %w", err)
		}

		// 只有IC卡读取成功时才读取后续字段
		if m.ICCardReadResult == 0x00 {
			// 驾驶员姓名长度
			var driverNameLength byte
			if driverNameLength, err = r.ReadByte(); err != nil {
				return 0, fmt.Errorf("read driver name length: %w", err)
			}

			if m.DriverName, err = r.ReadString(int(driverNameLength)); err != nil {
				return 0, fmt.Errorf("read driver name: %w", err)
			}

			// 从业资格证编码，固定20字节
			b, err := r.Read(20)
			if err != nil {
				return 0, fmt.Errorf("read qualification code: %w", err)
			}
			m.QualificationCode = string(trimRightZeros(b))

			// 发证机构名称长度
			var issuingAuthorityLength byte
			if issuingAuthorityLength, err = r.ReadByte(); err != nil {
				return 0, fmt.Errorf("read issuing authority length: %w", err)
			}

			if m.IssuingAuthority, err = r.ReadString(int(issuingAuthorityLength)); err != nil {
				return 0, fmt.Errorf("read issuing authority: %w", err)
			}

			// 证件有效期，BCD[4] YYYYMMDD
			validityStr, err := r.ReadBcd(4)
			if err != nil {
				return 0, fmt.Errorf("read certificate validity: %w", err)
			}

			if m.CertificateValidity, err = time.Parse("20060102", validityStr); err != nil {
				return 0, fmt.Errorf("parse certificate validity: %w", err)
			}

			// 驾驶员身份证号，长度20位，不足补0x00（2019版本新增）
			if m.protocolVersion >= Version2019 {
				b, err := r.Read(20)
				if err != nil {
					return 0, fmt.Errorf("read driver id card number: %w", err)
				}
				m.DriverIDCardNumber = string(trimRightZeros(b))
			}
		}
	}

	return len(data) - r.Len(), nil
}
